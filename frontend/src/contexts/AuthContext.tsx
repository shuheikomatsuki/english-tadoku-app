import React, { /*createContext,*/ useState,/* useContext,*/ useEffect } from 'react';
import type { ReactNode } from 'react';
import { AuthContext } from './authContext';

// interface AuthContextType {
//     isAuthenticated: boolean;
//     login: (token: string) => void;
//     logout: () => void;
// }

// export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC< { children: ReactNode }> = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
        }
    }, []);

    const login = (token: string) => {
        localStorage.setItem('token', token);
        setIsAuthenticated(true);
    };

    const logout = () => {
        localStorage.removeItem('token');
        setIsAuthenticated(false);
    };

    const value = { isAuthenticated, login, logout };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthProvider;

// export const useAuth = () => {
//     const context = useContext(AuthContext);
//     if (context === undefined) {
//         throw new Error('useAuth must be used within an AuthProvider');
//     }
//     return context;
// };