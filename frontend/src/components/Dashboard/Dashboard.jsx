import React from 'react';

function Dashboard() {
  return (
    <div>
      <header className="top-nav">
        <h1>
          <i className="material-icons">supervised_user_circle</i>
          User Management Dashboard
        </h1>
        <button className="button is-border" type="button">
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
