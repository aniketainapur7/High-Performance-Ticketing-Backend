import express from "express";
import trainRouter from "./routes/train.routes.js";
import seatRouter from "./routes/seat.routes.js";


export const app = express();

app.use(express.json());
app.use(trainRouter);
app.use(seatRouter);

app.get("/", (req, res) => {
    return res.json({ message: "Health Check OK" });
});

