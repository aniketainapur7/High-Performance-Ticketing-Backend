import { redis } from "../config/redis.js";
import prisma from "../config/prisma.js";

export const confirmBooking = async (trainId, seatId, userId) => {
    if (!trainId || !seatId || !userId) {
        return res.status(400).json({ error: "TrainId or SeatId or userId is missing" });
    }

    const lockKey = `seat:${seatId}:lock`;
    const bitmapKey = `train:${trainId}:seats`;

    const lockedBy = await redis.get(lockKey);
    const isBooked = await redis.getbit(bitmapKey,seatId-1);
    
    if(isBooked){
        throw new Error("Seat_Already_Booked");
    }
    if (!lockedBy) {
        throw new Error("LOCK_EXPIRED");
    }
    if (lockedBy !== userId) {
        throw new Error("LOCK_NOT_OWNED");
    }

    await redis.setbit(bitmapKey, seatId - 1, 1);

    const booking = await prisma.booking.create({
        data:{
            trainId,
            seatId,
            userId,
            status:"CONFIRMED"
        }
    });

    await redis.del(lockKey);

    return booking;

}