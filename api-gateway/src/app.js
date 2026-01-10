import express from "express";
import trainRouter from "./routes/train.routes.js";


export const app = express();

app.use(express.json());
app.use(trainRouter);

app.get("/", (req, res) => {
    return res.json({ message: "Health Check OK" });
});

