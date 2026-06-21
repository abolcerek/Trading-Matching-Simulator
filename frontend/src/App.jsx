import { useState } from 'react';
import './App.css'
import OrderBook from './OrderBook';
import OrderEntry from './OrderEntry';
import Login from './Login';

function App() {
  const [token, setToken] = useState('')
  return (
    <>
      <Login setToken = {setToken}/>
      <OrderEntry token={token}/>
      <OrderBook />
    </>
  );
}

export default App;
