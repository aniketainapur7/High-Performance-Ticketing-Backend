package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"booking-service/bookingpb"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type server struct {
	bookingpb.UnimplementedBookingServiceServer
	redis *redis.Client
	db    *pgx.Conn
}

func (s *server) ConfirmBooking(
	ctx context.Context,
	req *bookingpb.ConfirmBookingRequest,
) (*bookingpb.ConfirmBookingResponse, error) {

	seatId := req.SeatId
	userId := req.UserId
	trainId := req.TrainId

	lockKey := fmt.Sprintf("seat:%d:lock", seatId)
	lockedBy, err := s.redis.Get(ctx, lockKey).Result()

	bitmapKey := fmt.Sprintf("train:%d:seats", trainId)
	seatIndex := int64(seatId - 1)

	if err == redis.Nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "LOCK_EXPIRED",
		}, nil
	}

	if err != nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "REDIS_ERROR",
		}, nil
	}

	if lockedBy != userId {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "Locked By Other User",
		}, nil
	}

	booked, err := s.redis.GetBit(ctx, bitmapKey, seatIndex).Result()

	if err != nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "REDIS_ERROR",
		}, nil
	}

	if booked == 1 {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "Seat Already Booked",
		}, nil
	}
	_, err = s.redis.SetBit(ctx, bitmapKey, seatIndex, 1).Result()

	if err != nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "REDIS_ERROR",
		}, nil
	}

	_, err = s.db.Exec(ctx, `INSERT INTO bookings(id,seat_id,train_id,user_id,status)
						  VALUES(gen_random_uuid(),$1,$2,$3,"CONFIRMED")`,
						seatId, trainId, userId);
	
	if err != nil {
		s.redis.SetBit(ctx,bitmapKey,seatIndex,0)
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error: "DB_ERROR",
		},nil
	}

	s.redis.Del(ctx,lockKey)


	return &bookingpb.ConfirmBookingResponse{
		Success: true,
	}, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	err := rdb.Ping(ctx).Err()

	if err != nil {
		log.Fatal("Redis Connection Failed", err)
	}

	log.Println("Redis Connection Successful !!")

	dbConn, err := pgx.Connect(
		context.Background(),
		"postgresql://ticket:ticket@localhost:5432/ticketing",
	)

	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	log.Println("Connected to Postgres")

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	bookingpb.RegisterBookingServiceServer(
		grpcServer,
		&server{
			redis: rdb,
			db:    dbConn,
		},
	)

	log.Println("GRPC SERVICE RUNNING ON PORT 50051")
	grpcServer.Serve(lis)
}
