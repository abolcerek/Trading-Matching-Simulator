import './App.css'
import { useState,  } from 'react';

function OrderEntry({token}) {
    const [side, setSide] = useState('buy');
    const [price, setPrice] = useState('');
    const [type, setType] = useState('limit');
    const [quantity, setQuantity] = useState('');
    const handleOrder = async () => {
        const response = await fetch('http://localhost:8080/api/orders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify({
                side: side,
                type: type,
                price: Number(price),
                quantity: Number(quantity)
            })
        });
        console.log(response)
    }

    return (
        <div>
            <div className='auth-card'>
                <label> Enter order side: </label>
                    <select value={side} className='auth-input' onChange={e => setSide(e.target.value)}>
                        <option value="buy">Buy</option>
                        <option value="sell">Sell</option>
                    </select>
                <label> Enter order price: </label>
                    <input type="number"
                    value={price} className='auth-input'
                    onChange={e => setPrice(e.target.value)} />
                <label> Enter order type: </label>
                    <select value={type} className='auth-input' onChange={e => setType(e.target.value)}>
                        <option value="limit">Limit</option>
                        <option value="market">Market</option>
                    </select>
                <label> Enter order quantity: </label>
                    <input type="number" className='auth-input'
                    value={quantity}
                    onChange={e => setQuantity(e.target.value)} />
                <button className='auth-button' onClick={handleOrder}>
                    Submit Order
                </button>
            </div>
        </div>
    )

}

export default OrderEntry;