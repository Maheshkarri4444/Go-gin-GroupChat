import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';
import Allapi from '../common';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Check cookies on mount and updates
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await axios.get(Allapi.checkauth.url, {
          withCredentials: true
        });
        
        if (response.data && response.data.userid && response.data.username) {
          setUser({
            userid: response.data.userid,
            username: response.data.username
          });
        }
      } catch (error) {
        console.error('Auth verification error:', error);
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = (userData) => {
    setUser(userData);
  };

  const logout = async () => {
    try {
      const response = await axios.post(Allapi.logout.url, {}, {
        withCredentials: true
      });

      if (response.status === 200) {
        setUser(null);
        return true;
      }
      return false;
    } catch (error) {
      console.error('Logout error:', error);
      return false;
    }
  };

  if (loading) {
    return null; // or a loading spinner
  }

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};