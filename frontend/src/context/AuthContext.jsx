import { createContext, useContext, useState, useEffect } from 'react';
import api from '../api/client';
import { registerPushNotifications, unregisterPushNotifications } from '../utils/pushNotifications';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const saved = localStorage.getItem('user');
        if (token && saved) {
            setUser(JSON.parse(saved));
            // Verify token is still valid
            api.get('/auth/me').then(res => {
                setUser(res.data);
                localStorage.setItem('user', JSON.stringify(res.data));
                // Register push notifications for returning users
                registerPushNotifications();
            }).catch(() => {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
                setUser(null);
            }).finally(() => setLoading(false));
        } else {
            setLoading(false);
        }
    }, []);

    const login = async (email, password) => {
        const res = await api.post('/auth/login', { email, password });
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('user', JSON.stringify(res.data.user));
        setUser(res.data.user);
        // Register push notifications on login
        registerPushNotifications();
        return res.data;
    };

    const register = async (data) => {
        const res = await api.post('/auth/register', data);
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('user', JSON.stringify(res.data.user));
        setUser(res.data.user);
        // Register push notifications on register
        registerPushNotifications();
        return res.data;
    };

    const logout = () => {
        // Unregister push notifications on logout
        unregisterPushNotifications();
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, register, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error('useAuth must be used within AuthProvider');
    return ctx;
}
