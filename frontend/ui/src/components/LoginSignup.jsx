import React, { useState } from 'react';
import { signup, login } from '../api';

export default function LoginSignup({ onAuth }){
  const [isLogin, setIsLogin] = useState(true);
  const [form, setForm] = useState({username: '', email: '', password: ''});
  const [error, setError] = useState('');

  const handleInput = (e) => {
    setForm({...form, [e.target.name]: e.target.value});
  };

  async function submit(e) {
    e.preventDefault();
    setError('');
    
    try {
      if (isLogin) {
        const res = await login({ email: form.email, password: form.password });
        onAuth(res); 
      } else {
        await signup({ 
          username: form.username, 
          email: form.email, 
          password: form.password 
        });
        const loginRes = await login({ email: form.email, password: form.password });
        onAuth(loginRes);
      }
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="login-signup-card">
      <h5>{isLogin ? 'Login' : 'Sign Up'}</h5>
      {error && <div className="alert alert-danger p-2">{error}</div>}
      
      <form onSubmit={submit}>
        {!isLogin && (
          <div className="mb-2">
            <input 
              className="form-control" 
              placeholder="Username" 
              name="username"
              value={form.username} 
              onChange={handleInput}
              required
            />
          </div>
        )}
        <div className="mb-2">
          <input 
            type="email" 
            className="form-control" 
            placeholder="Email" 
            name="email"
            value={form.email} 
            onChange={handleInput}
            required
          />
        </div>
        <div className="mb-2">
          <input 
            type="password" 
            className="form-control" 
            placeholder="Password" 
            name="password"
            value={form.password} 
            onChange={handleInput}
            required
          />
        </div>
        <div className="d-grid gap-2">
          <button className="btn btn-primary" type="submit">
            {isLogin ? 'Login' : 'Sign Up'}
          </button>
          <button 
            type="button" 
            className="btn btn-outline-secondary" 
            onClick={() => {
              setIsLogin(!isLogin);
              setError('');
            }}
          >
            {isLogin ? 'Need an account? Sign Up' : 'Have an account? Login'}
          </button>
        </div>
      </form>
    </div>
  );
}