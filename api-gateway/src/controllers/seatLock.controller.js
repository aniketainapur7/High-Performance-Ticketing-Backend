import { redis } from "../config/redis.js";
import { lockSeat } from "../services/seatLock.service.js";

export const lockSeatHandler = async (req,res) => {
    const seatId = Number(req.params.seatId);
    const {userId} = req.body;

    if(!seatId || !userId){
        return res.status(400).json({error : "seatId and userId required"});
    }

    const locked = await lockSeat(seatId,userId);

    if(!locked){
        return res.status(409).json({
            error:"seat already locked"
        });
    }

    return res.json({
        message:"Seat Locked Succcesfully",
        expiresInSeconds : 30
    });
};