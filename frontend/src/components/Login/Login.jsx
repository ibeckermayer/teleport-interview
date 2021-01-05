import React from 'react';
import { useHistory } from 'react-router-dom';

const Login = () => {
  const history = useHistory();

  const tryLogin = async e => {
    // Form validation is handled by html5
    e.preventDefault();
    const res = await fetch('api/login', {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: e.target.email.value,
        password: e.target.password.value,
      }),
    });
    if (res.ok) {
      history.push('/dashboard');
    }
    // TODO: check for Unauthorized and alert user that username/pwd is incorrect
  };

  return (
    <form className="login-form" onSubmit={tryLogin}>
      <h1>Sign Into Your Account</h1>

      <div>
        <label htmlFor="email">Email Address</label>
        <input
          type="email"
          id="email"
          className="field"
          autoComplete="username"
          name="email"
          required
        />
      </div>

      <div>
        <label htmlFor="password">Password</label>
        <input
          type="password"
          id="password"
          className="field"
          autoComplete="current-password"
          name="password"
          required
        />
      </div>

      <input
        type="submit"
        value="Login to my Dashboard"
        className="button block"
      />
    </form>
  );
};

export default Login;
