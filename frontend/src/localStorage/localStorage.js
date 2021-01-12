export function setLocalStore(session) {
  localStorage.setItem('store', JSON.stringify(session));
}

export function getLocalStore() {
  return JSON.parse(localStorage.getItem('store'));
}
