import { useState, useEffect } from "react";

const WS_URL = "ws://localhost:8080/ws";

export function useWebSocket() {
  const [prices, setPrices] = useState<{ [symbol: string]: number }>({});

  useEffect(() => {
    let socket: WebSocket;
    let reconnectAttempts = 0;

    const connectWebSocket = () => {
      socket = new WebSocket(WS_URL);

      socket.onmessage = (event) => {
        setPrices(JSON.parse(event.data));
        reconnectAttempts = 0; // Reset counter on successful message
      };

      socket.onerror = () => {
        console.error("WebSocket error. Attempting reconnect...");
        reconnect();
      };

      socket.onclose = () => {
        console.warn("WebSocket closed. Reconnecting...");
        reconnect();
      };
    };

    const reconnect = () => {
      reconnectAttempts++;
      setTimeout(connectWebSocket, Math.min(5000, reconnectAttempts * 1000)); // Exponential backoff
    };

    connectWebSocket();
    return () => socket.close();
  }, []);

  return prices;
}