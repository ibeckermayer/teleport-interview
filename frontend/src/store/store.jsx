import React, { createContext, useReducer, useEffect } from 'react';
import { getLocalStore, setLocalStore } from '../localStorage';

// Default state for the store
const initialStore = {
  sessionID: '', // Empty string represents user not authenticated, otherwise they are authenticated
  plan: '', // "FREE" or "ENTERPRISE"
};

// On app load, try to fetch a saved store from localStorage
const localStore = getLocalStore();

// Updates the store with newStoreValues
// Pass `newStoreValues = null` to reset the store
const storeReducer = (store, newStoreValues) => {
  if (newStoreValues === null) {
    localStorage.removeItem('store');
    return initialStore;
  }
  return { ...store, ...newStoreValues };
};

// Global React.Context value for components to pass into useContext in order to access the store
const StoreContext = createContext();

// createStore should wrap the entire application, so that any component in the application
// can access the global store's state as a StoreContext.Consumer
const createStore = WrappedComponent => props => {
  // Initialize with the store saved in localStorage if it exists, else initialize with initialStore
  const [store, dispatch] = useReducer(
    storeReducer,
    localStore || initialStore
  );

  // Save the store to localStorage every time it's changed
  useEffect(() => {
    setLocalStore(store);
  }, [store]);

  return (
    <StoreContext.Provider value={[store, dispatch]}>
      <WrappedComponent {...props} />
    </StoreContext.Provider>
  );
};

export { StoreContext, createStore };
