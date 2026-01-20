import { Router } from "express";
import { confirmBookingHandler } from "../controllers/booking.controller.js";

const bookingRouter = Router();

bookingRouter.post("/api/bookings/confirm",confirmBookingHandler);

export default bookingRouter