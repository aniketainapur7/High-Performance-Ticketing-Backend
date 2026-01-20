import prisma from "../config/prisma.js"
import { redis } from "../config/redis.js"


export const getSeatAvailability = async (trainId) => {
    const seats = await prisma.seat.findMany({
        where:{trainId},
        orderBy:{seatNo:"asc"}
    });

    const bitmapKey = `train:${trainId}:seats`;

    const result = [];

    for(let i = 0; i < seats.length; i++){
        const booked = await redis.getbit(bitmapKey,i);
        const lockKey = `seat:${seats[i].id}:lock`
        const lockedBy = await redis.get(lockKey);

        result.push({
            seatId : seats[i].id,
            seatNo : seats[i].seatNo,
            berthType : seats[i].berthType,
            status: booked == 1 ? "BOOKED" : lockedBy ? "LOCKED" : "AVAILABLE"
        });
    }

    return result;
}