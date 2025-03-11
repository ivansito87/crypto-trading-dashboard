import { memo } from "react";
import { useWebSocket } from "../hooks/useWebSocket";

const PriceTicker = memo(() => {
  const prices = useWebSocket();

  return (
    <div>
      <h2>📈 Live Crypto Prices</h2>
      <ul>
        {Object.entries(prices).map(([symbol, price]) => (
          <li key={symbol}>
            {symbol}: ${price.toFixed(2)}
          </li>
        ))}
      </ul>
    </div>
  );
});

export default PriceTicker;