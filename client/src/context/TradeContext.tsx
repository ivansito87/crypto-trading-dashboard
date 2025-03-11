import { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

const TradeContext = createContext<any>(null);

export const TradeProvider = ({ children }: { children: React.ReactNode }) => {
  const [trades, setTrades] = useState([]);

  useEffect(() => {
    const fetchTrades = async () => {
      const response = await axios.get("http://localhost:8080/trades");
      setTrades(response.data);
    };
    fetchTrades();
  }, []);

  return <TradeContext.Provider value={trades}>{children}</TradeContext.Provider>;
};

export const useTrades = () => useContext(TradeContext);