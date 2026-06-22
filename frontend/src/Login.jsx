import './App.css'
import { useState,  } from 'react';
import { useNavigate, Link} from 'react-router-dom';


function Login ({setToken}){
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

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
        if (response.ok) {
            setToken(data.token)
            navigate('/trade')
        } else {
            console.log(data) // Implement later
        }
    }
    return (
        <div className='auth-page'>
            <div className='auth-card'>
                <h2>Log In</h2>
                <label>Enter email: </label>
                    <input type="text" className='auth-input'
                    value={email}
                    onChange={e => setEmail(e.target.value)}/>
                <label>Enter password: </label>
                    <input type="password" className='auth-input'
                    value={password}
                    onChange={e => setPassword(e.target.value)}/>
                <button className='auth-button' onClick={handleLogin}>
                    Submit
                </button>
                <Link to="/signup">Need an account? Sign up</Link>
            </div>
        </div>
    )
}


export default Login