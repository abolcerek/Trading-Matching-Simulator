import { useNavigate, Link} from 'react-router-dom';
import './App.css'
import { useState,  } from 'react';


function SignUp (){
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const navigate = useNavigate()

    const handleSignUp = async () => {
        const response = await fetch('http://localhost:8080/api/users', {
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
            navigate('/login')
        } else {
            console.log(data)
        }
    }
    return (
        <div className='auth-page'>
            <div className='auth-card'>
                <h2>Sign up</h2>
                <label>Enter email: </label>
                    <input type="text" className='auth-input'
                    value={email}
                    onChange={e => setEmail(e.target.value)}/>
                <label>Enter password: </label>
                    <input type="password" className='auth-input'
                    value={password}
                    onChange={e => setPassword(e.target.value)}/>
                <button className='auth-button' onClick={handleSignUp}>
                    Submit
                </button>
                <Link to="/login">Already have an account? Log in</Link>
            </div>
        </div>
    )
}


export default SignUp