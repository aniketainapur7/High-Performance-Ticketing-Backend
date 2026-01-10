import { Router } from "express";
import { searchTrains } from "../controllers/train.controller.js";

const trainRouter = Router();

trainRouter.get("/api/trains", searchTrains);

export default trainRouter;