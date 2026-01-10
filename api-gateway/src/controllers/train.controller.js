import prisma from "../config/prisma.js";

export const searchTrains = async (req, res) => {
    try {
        const {source, destination} = req.query;
        
        if(!source || !destination){
            return res.status(400).json({error: "Source and Destination are required"});
        }
    
        const trains = await prisma.train.findMany({
            where:{
                source : source,
                destination: destination
            },
            select:{
                id:true,
                name:true,
                source:true,
                destination:true
            }
        });
        
        res.json(trains);
    } catch (error) {
        res.status(500).json({error:" Internal Server Error"});
    }
}