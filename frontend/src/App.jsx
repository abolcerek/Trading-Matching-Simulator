import './App.css'
import { useState, useEffect } from 'react';

function App() {

  const [book, setBook] = useState(
    {Bid: [{Price: 100, Quantity: 7}], 
    Ask: [{Price: 150, Quantity: 3}]
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
      <section id="top">
        <div>
          <h1>Order Book</h1>
        </div>
      </section>
      <section id="center">
        <div>
          <h2>Book</h2>
          <ul>
            {Object.entries(book).map(([key, value]) => (
              <li key={key}>
                <p>{key}</p> 
                <ul>
                  {value.map((item, index) => (
                    <li key={index}>
                      Price: {item.Price}, Quantity: {item.Quantity}
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
        </div>
      </section>
    </>
  );
}

export default App;
