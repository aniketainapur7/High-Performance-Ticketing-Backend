import { Router } from "express";
import { getSeats } from "../controllers/seat.controller.js";


const seatRouter = Router();

seatRouter.get("/api/trains/:trainId/seats", getSeats);

export default seatRouter