import { getSeatAvailability } from "../services/seatAvailability.service.js";

export const getSeats = async (req, res) => {
    const trainId = Number(req.params.trainId);
    if (!trainId) {
        return res.status(400).json({ error: "invalid trainId" });
    }
    const seats = await getSeatAvailability(trainId);

    return res.json(seats);
    
}