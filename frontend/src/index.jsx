import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import './index.css';
import { SessionContextProvider } from './session';

const App = () => {
  return (
    <SessionContextProvider>
      <Router>
        <Switch>
          <Route exact path="/">
            <Login />
          </Route>
          <Route path="/dashboard">
            <Dashboard />
          </Route>
        </Switch>
      </Router>
    </SessionContextProvider>
  );
};

ReactDOM.render(<App />, document.getElementById('root'));
