import { env } from "./env";
import Redis from 'ioredis';

export const redis = new Redis({
    host: env.redis.host,
    port: env.redis.port,
})

