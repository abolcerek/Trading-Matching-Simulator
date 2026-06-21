import './App.css'
import OrderBook from './OrderBook';
import OrderEntry from './OrderEntry';

function TradingView ({token}) {
    return (
        <>
        <OrderBook/>
        <OrderEntry token={token}/>
        </>
    )
}

export default TradingView