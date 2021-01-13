import React, { createContext, useState } from 'react';
import { getLocalSession, setLocalSession } from '../localStorage';

// Global React.Context value for components to pass into useContext in order to access the session
const SessionContext = createContext();

// AppContext should wrap the entire application, so that any component in the application
// can access the global sessions's state as a StoreContext.Consumer
// eslint-disable-next-line react/prop-types
const SessionContextProvider = ({ children }) => {
  // Initialize the app with session set to the session saved in localStorage (or null if no saved session exists)
  // setSessionI is "I" for internal use only (like a private method)
  const [session, setSessionI] = useState(getLocalSession());

  // Consumers can update the session by calling setSession, or clear the session by passing in null
  const setSession = updates => {
    const newSession = updates === null ? null : { ...session, updates };
    setSessionI(newSession); // update session
    setLocalSession(newSession); // keep localStorage in sync
  };

  return (
    <SessionContext.Provider value={[session, setSession]}>
      {children}
    </SessionContext.Provider>
  );
};

export { SessionContext, SessionContextProvider };
