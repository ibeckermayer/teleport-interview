import React, { useContext } from 'react';
import api from '../../api';
import { StoreContext } from '../../store';

function Dashboard() {
  const { setStore } = useContext(StoreContext);
  const logout = async e => {
    e.preventDefault();
    try {
      await api.delete('/logout');
      setStore(null);
    } catch (error) {
      // TODO: Probably should alert the user something along the lines of
      // "application error, please try again and contact customer support if the error persists"
      // eslint-disable-next-line no-console
      console.error(error);
    }
  };

  return (
    <div>
      <header className="top-nav">
        <h1>
          <i className="material-icons">supervised_user_circle</i>
          User Management Dashboard
        </h1>
        <button className="button is-border" type="button" onClick={logout}>
          Logout
        </button>
      </header>

      <div className="alert is-error">
        You have exceeded the maximum number of users for your account, please
        upgrade your plan to increaese the limit.
      </div>
      <div className="alert is-success">
        Your account has been upgraded successfully!
      </div>

      <div className="plan">
        <header>Startup Plan - $100/Month</header>

        <div className="plan-content">
          <div className="progress-bar">
            <div style={{ width: '35%' }} className="progress-bar-usage" />
          </div>

          <h3>Users: 35/100</h3>
        </div>

        <footer>
          <button className="button is-success" type="button">
            Upgrade to Enterprise Plan
          </button>
        </footer>
      </div>
    </div>
  );
}

export default Dashboard;
