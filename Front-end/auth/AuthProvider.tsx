import React, { createContext, useContext, useState, useEffect } from 'react';

interface TokenProfile {
  acr: string;
  at_hash: string;
  aud: string;
  azp: string;
  email: string;
  email_verified: boolean;
  exp: number;
  family_name: string;
  given_name: string;
  iat: number;
  iss: string;
  jti: string;
  name: string;
  preferred_username: string;
  sid: string;
  sub: string;
  typ: string;
}

interface TokenType {
  access_token: string;
  profile?: TokenProfile;
}

interface AuthContextType {
  token: TokenType | null;
  setToken: (token: TokenType) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setTokenState] = useState<TokenType | null>(() => {
    const stored = sessionStorage.getItem('Token');
    return stored ? JSON.parse(stored) : null;
  });

  const setToken = (t: TokenType) => {
    setTokenState(t);
    sessionStorage.setItem('Token', JSON.stringify(t));
  };

  const logout = () => {
    //setTokenState(null);
    //sessionStorage.clear();
    //localStorage.clear();
    //window.location.href = '/login';
  };

  useEffect(() => {
    console.log('Token actualizado:', token?.access_token);
  }, [token]);

  const isAuthenticated = Boolean(token?.access_token);

  return (
    <AuthContext.Provider value={{ token, setToken, logout, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
