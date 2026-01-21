import prisma from "../config/prisma.js";

export const addSeats = async (req,res) =>{
    const {seatNo, trainId, berthType } = req.body;
    if(!seatNo || !trainId || !berthType){
        return res.status(400).json({error:"All Fields Necessary"});
    }

    await prisma.seat.create({
        data:{
            seatNo,
            berthType,
            trainId
        }
    });

    return res.json({message:"success"})
}

export const addTrain = async (req,res) =>{
    const {name, source, destination } = req.body;
    if(!name || !source || !destination){
        return res.status(400).json({error:"All Fields Necessary"});
    }

    await prisma.train.create({
        data:{
            name,
            source,
            destination
        }
    });

    return res.json({message:"success"})
}