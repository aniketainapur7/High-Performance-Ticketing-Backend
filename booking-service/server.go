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

	bitmapKey := fmt.Sprintf("train:%d:seats", trainId)
	seatIndex := int64(seatId - 1)

	result, err := s.redis.Eval(ctx, bookSeatLua, []string{lockKey, bitmapKey}, seatIndex, userId).Result()

	if err != nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if result != "OK" {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "UNKNOWN_ERROR",
		}, nil
	}
	_, err = s.db.Exec(ctx,
		`INSERT INTO "Booking" ("seatId","trainId","userId","status")
   		VALUES ($1,$2,$3,'CONFIRMED')`,
		seatId, trainId, userId,
	)

	if err != nil {
		s.redis.SetBit(ctx, bitmapKey, seatIndex, 0)
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "DB_ERROR",
		}, nil
	}

	return &bookingpb.ConfirmBookingResponse{
		Success: true,
	}, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx := context.Background()

	err := rdb.Ping(ctx).Err()

	if err != nil {
		log.Fatal("Redis Connection Failed", err)
	}

	log.Println("Redis Connection Successful !!")

	dbConn, err := pgx.Connect(
		context.Background(),
		"postgresql://ticket:ticket@postgres:5432/ticketing",
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

var bookSeatLua = `
-- KEYS[1] = lockKey
-- KEYS[2] = bitmapKey
-- ARGV[1] = seatIndex
-- ARGV[2] = userId

local lockedBy = redis.call("GET",KEYS[1])

if not lockedBy then 
	return {err = "LOCK_EXPIRED"}
end

if lockedBy ~= ARGV[2] then
	return {err = "LOCK_NOT_OWNED"}
end

local booked = redis.call("GETBIT",KEYS[2],ARGV[1])

if booked == 1 then
	return {err = "SEAT_ALREADY_BOOKED"}
end

redis.call("SETBIT",KEYS[2],ARGV[1],1)
redis.call("DEL",KEYS[1])

return "OK"
`
