import './App.css'
import { useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import TradingView from './TradingView';
import Login from './Login';
import SignUp from './SignUp';

function App() {
  const [token, setToken] = useState('')
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/signup' element={< SignUp/>} />
        <Route path='/login' element={<Login setToken = {setToken}/>} />
        <Route path='/trade' element={token ? <TradingView token={token} /> : <Navigate to="/login" />} />
        <Route path='/' element={<Login setToken={setToken} />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;


