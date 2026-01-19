import { Router } from "express";


const seatRouter = Router();

seatRouter.get("/api/trains/:trainId/seats", getSeats);

export default seatRouter