import { useState, useCallback } from "react";
import axios from "axios";

const TradeForm = () => {
  const [symbol, setSymbol] = useState("BTC");
  const [type, setType] = useState("buy");
  const [amount, setAmount] = useState(1);
  const [message, setMessage] = useState("");

  const placeTrade = useCallback(async () => {
    try {
      const response = await axios.post("http://localhost:8080/order", {
        symbol,
        type,
        amount: Number(amount),
      });
      setMessage(`Trade successful: ${response.data.symbol} ${response.data.type} ${response.data.amount} @ $${response.data.price}`);
    } catch (error) {
      setMessage("Error placing trade");
    }
  }, [symbol, type, amount]);

  return (
    <div>
      <h2>🛒 Place Trade</h2>
      <select value={symbol} onChange={(e) => setSymbol(e.target.value)}>
        <option value="BTC">BTC</option>
        <option value="ETH">ETH</option>
        <option value="ADA">ADA</option>
      </select>
      <select value={type} onChange={(e) => setType(e.target.value)}>
        <option value="buy">Buy</option>
        <option value="sell">Sell</option>
      </select>
      <input type="number" value={amount} onChange={(e) => setAmount(Number(e.target.value))} min="1" />
      <button onClick={placeTrade}>Execute Trade</button>
      {message && <p>{message}</p>}
    </div>
  );
};

export default TradeForm;