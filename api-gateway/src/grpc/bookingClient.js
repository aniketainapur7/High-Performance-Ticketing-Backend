import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const protoPath = path.join(__dirname, "../../proto/booking.proto")

const packageDef = protoLoader.loadSync(protoPath, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const proto = grpc.loadPackageDefinition(packageDef);

export const bookingClient = new proto.booking.BookingService(
    "booking-service:50051",
    grpc.credentials.createInsecure()
);