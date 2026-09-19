import { createContext, useContext, useState, useEffect } from 'react';
import { api } from '../api/client';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
      const [user, setUser] = useState(null);
      const [loading, setLoading] = useState(true);

      useEffect(() => {
            const token = localStorage.getItem('token');
            const storedUser = localStorage.getItem('user');

            if (token && storedUser) {
                  setUser(JSON.parse(storedUser));
            }

            setLoading(false);
      }, []);

      async function login(email, password) {
            const data = await api.post('/login', { email, password });
            localStorage.setItem('token', data.token);

            // Decode the JWT payload to learn the user's role, without a second API call.
            const payload = JSON.parse(atob(data.token.split('.')[1]));
            const loggedInUser = { id: payload.sub, role: payload.role };

            localStorage.setItem('user', JSON.stringify(loggedInUser));
            setUser(loggedInUser);
            return loggedInUser;
      }

      async function register(name, email, password) {
            await api.post('/register', { name, email, password, role: 'student' });
      }

      function logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            setUser(null);
      }

      return (
            <AuthContext.Provider value={{ user, loading, login, register, logout }}>
                  {children}
            </AuthContext.Provider>
      );
}

export function useAuth() {
      return useContext(AuthContext);
}