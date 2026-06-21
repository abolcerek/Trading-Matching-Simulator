import './App.css'
import { useState,  } from 'react';


function Login ({setToken}){
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const handleLogin = async () => {
        const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        });
        const data = await response.json()
        setToken(data.token)
    }
    return (
        <div>
            <label>Enter email: </label>
                <input type="text" 
                value={email}
                onChange={e => setEmail(e.target.value)}/>
            <label>Enter password: </label>
                <input type="password" 
                value={password}
                onChange={e => setPassword(e.target.value)}/>
            <button onClick={handleLogin}>
                Submit
            </button>
        </div>
    )
}


export default Login