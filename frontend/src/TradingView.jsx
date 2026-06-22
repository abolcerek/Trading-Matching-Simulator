import './App.css'
import OrderBook from './OrderBook';
import OrderEntry from './OrderEntry';

function TradingView({token}) {
  return (
    <div className="trading-view">
      <h1 className="page-title">Order Book</h1>
      <div className="trading-row">
        <OrderBook />
        <OrderEntry token={token} />
      </div>
    </div>
  )
}

export default TradingView