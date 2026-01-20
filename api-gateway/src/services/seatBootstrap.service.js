import prisma from "../config/prisma.js";
import { redis } from "../config/redis.js";

export const bootstrapSeatsForTrain = async (trainId) => {
    const seats = await prisma.seat.findMany({
        where:{trainId},
        orderBy:{seatNo:"asc"}
    });
    const redisKey = `train:${trainId}:seats`;
    console.log(seats)
    for(let i = 0; i < seats.length ; i++){
        await redis.setbit(redisKey,i,0);
        console.log(i);
    }

    return seats.length;
}

await bootstrapSeatsForTrain(1);