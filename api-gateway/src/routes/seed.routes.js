import {Router} from "express"
import { addSeats, addTrain } from "../controllers/seed.controller.js";
const seedRouter = Router();

seedRouter.post("/api/add-train",addTrain);
seedRouter.post("/api/add-seats",addSeats);

export default seedRouter;