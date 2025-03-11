import { useTrades } from "../context/TradeContext";


export interface Trade {
  id: number;
  symbol: string;
  type: string;
  amount: number;
  price: number;
}

const TradeHistory = () => {
  const trades = useTrades();

  return (
    <div>
      <h2>📜 Trade History</h2>
      <ul>
        {trades.map((trade: Trade) => (
          <li key={trade.id}>
            {trade.symbol} {trade.type} {trade.amount} @ ${trade.price}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TradeHistory;