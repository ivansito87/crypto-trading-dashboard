import PriceTicker from "./components/PriceTicker";
import TradeForm from "./components/TradeForm";
import TradeHistory from "./components/TradeHistory";
import { TradeProvider } from "./context/TradeContext";

function App() {
  return (
    <TradeProvider>
      <div>
        <h1>🚀 Crypto Trading Dashboard</h1>
        <PriceTicker />
        <TradeForm />
        <TradeHistory />
      </div>
    </TradeProvider>
  );
}

export default App;