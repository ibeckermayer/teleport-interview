export function setLocalSession(session) {
  localStorage.setItem('session', JSON.stringify(session));
}

export function getLocalSession() {
  return JSON.parse(localStorage.getItem('session'));
}
