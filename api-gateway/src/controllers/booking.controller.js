import { confirmBooking } from "../services/booking.service.js";

export const confirmBookingHandler = async (req, res) => {
    const { seatId, trainId, userId } = req.body;

    if (!seatId || !trainId || !userId) {
        return res.status(400).json({
            error: "seatId, trainId and userId required",
        });
    }

    try {
        const booking = await confirmBooking(trainId, seatId, userId);
        return res.json({
            message: "Booking Confirmed",
            bookingId: booking.id
        });
    } catch (err) {
        if (err.message === "LOCK_EXPIRED") {
            return res.status(410).json({ error: "Seat lock expired" });
        }
        if (err.message === "LOCK_NOT_OWNED") {
            return res.status(403).json({ error: "Seat locked by another user" });
        }
        if (err.message === "Seat_Already_Booked") {
            return res.status(403).json({ error: "Seat booked by another user" });
        }

        console.error(err);
        res.status(500).json({ error: "internal server error" });
    }
}