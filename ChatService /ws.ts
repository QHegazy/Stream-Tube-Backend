import WebSocket from "ws";
import * as dotenv from "dotenv";
dotenv.config();
// Create WebSocket server
const wss = new WebSocket.Server({ port: 5050 });
wss.on("connection", (ws) => {
  console.log("Client connected");
});
