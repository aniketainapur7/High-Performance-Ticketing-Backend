import express from "express";
import trainRouter from "./routes/train.routes.js";
import seatRouter from "./routes/seat.routes.js";
import seatLockRouter from "./routes/seatLock.routes.js";
import bookingRouter from "./routes/booking.routes.js";
import seedRouter from "./routes/seed.routes.js";


export const app = express();

app.use(express.json());
app.use(trainRouter);
app.use(seatRouter);
app.use(seatLockRouter)
app.use(bookingRouter);
app.use(seedRouter);

app.get("/", (req, res) => {
    return res.json({ message: "Health Check OK" });
});

