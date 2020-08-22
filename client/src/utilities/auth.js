import { createContext, useContext } from 'react';

export const AuthContext = createContext();

export function useAuth() {
  return useContext(AuthContext);
}

export function logout() {
  localStorage.removeItem("tokens")
}