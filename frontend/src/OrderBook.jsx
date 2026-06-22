import './App.css'
import { useState, useEffect } from 'react';

function OrderBook() {

  const [book, setBook] = useState(
    {Bid: [{Price: 0, Quantity: 0}], 
    Ask: [{Price: 0, Quantity: 0}]
    });

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/api/ws')
    ws.onmessage = (event) => {
      setBook(JSON.parse(event.data))
    }

    return () => {
      ws.close()
    };
  }, [])
  return (
    <>
      <div >
        <h2>Book</h2>
        <div className='book'>
          <div className="asks">
            {[...book.Ask].reverse().map((level) => (
              <div className="ask-row" key={level.Price}>
                <span>{level.Price}</span><span>{level.Quantity}</span>
              </div>
            ))}
          </div>
          <div className="bids">
            {book.Bid.map((level) => (
              <div className="bid-row" key={level.Price}>
                <span>{level.Price}</span><span>{level.Quantity}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </>
  );
}

export default OrderBook;
