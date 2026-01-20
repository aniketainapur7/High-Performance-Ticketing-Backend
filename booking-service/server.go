package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"booking-service/bookingpb"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type server struct {
	bookingpb.UnimplementedBookingServiceServer
	redis *redis.Client
}

func (s *server) ConfirmBooking(
	ctx context.Context,
	req *bookingpb.ConfirmBookingRequest,
) (*bookingpb.ConfirmBookingResponse, error) {

	seatId := req.SeatId
	userId := req.UserId

	lockKey := fmt.Sprintf("seat:%d:lock", seatId)
	lockedBy, err := s.redis.Get(ctx, lockKey).Result()

	bitmapKey := fmt.Sprintf("train:%d:seats",req.TrainId)
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
			Error: "Locked By Other User",
		},nil
	}

	booked,err := s.redis.GetBit(ctx,bitmapKey,seatIndex).Result()

	if err != nil {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error:   "REDIS_ERROR",
		}, nil
	}

	if booked == 1 {
		return &bookingpb.ConfirmBookingResponse{
			Success: false,
			Error: "Seat Already Booked",
		},nil
	}
	_,err = s.redis.SetBit(ctx,bitmapKey,seatIndex,1).Result()

	if err != nil {
    	return &bookingpb.ConfirmBookingResponse{
    	Success: false,
    	Error:   "REDIS_ERROR",
    }, nil
  	}


	return &bookingpb.ConfirmBookingResponse{
		Success:   true,
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

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	bookingpb.RegisterBookingServiceServer(
		grpcServer,
		&server{redis: rdb},
	)

	log.Println("GRPC SERVICE RUNNING ON PORT 50051")
	grpcServer.Serve(lis)
}
