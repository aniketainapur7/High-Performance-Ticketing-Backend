import { Router } from "express";
import { lockSeatHandler } from "../controllers/seatLock.controller.js";

const seatLockRouter = Router();

seatLockRouter.post("/api/seats/:seatId/lock",lockSeatHandler);

export default seatLockRouter