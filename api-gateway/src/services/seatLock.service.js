import { redis } from "../config/redis.js";

const LOCK_TTL_SECONDS = 30;

export const lockSeat = async (seatId, userId) => {
    const lockKey = `seat:${seatId}:lock`;

    const result = await redis.set(
        lockKey,
        userId,
        "NX",
        "EX",
        LOCK_TTL_SECONDS
    );

    return result ? true : false;
};

