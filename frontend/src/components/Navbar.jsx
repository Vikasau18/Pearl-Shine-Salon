import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Scissors, Bell, Heart, Menu, X, LogOut, LayoutDashboard, User, Calendar } from 'lucide-react';
import { useState, useEffect } from 'react';
import api from '../api/client';
import './Navbar.css';

export default function Navbar() {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [menuOpen, setMenuOpen] = useState(false);
    const [unread, setUnread] = useState(0);

    useEffect(() => {
        if (user) {
            api.get('/notifications/unread-count')
                .then(res => setUnread(res.data.unread_count))
                .catch(() => { });
        }
    }, [user]);

    const handleLogout = () => {
        logout();
        navigate('/');
        setMenuOpen(false);
    };

    return (
        <nav className="navbar glass">
            <div className="container navbar-inner">
                <Link to="/" className="navbar-brand">
                    <Scissors size={24} />
                    <span className="gradient-text">Saloon</span>
                </Link>

                <div className={`navbar-links ${menuOpen ? 'active' : ''}`}>
                    <Link to="/salons" onClick={() => setMenuOpen(false)}>Explore</Link>

                    {user ? (
                        <>
                            <Link to="/appointments" onClick={() => setMenuOpen(false)}>
                                <Calendar size={16} /> My Bookings
                            </Link>
                            <Link to="/favorites" onClick={() => setMenuOpen(false)}>
                                <Heart size={16} /> Favorites
                            </Link>

                            {(user.role === 'salon_owner' || user.role === 'admin') && (
                                <Link to="/dashboard" onClick={() => setMenuOpen(false)}>
                                    <LayoutDashboard size={16} /> Dashboard
                                </Link>
                            )}

                            <Link to="/notifications" className="nav-notif" onClick={() => setMenuOpen(false)}>
                                <Bell size={16} />
                                {unread > 0 && <span className="notif-badge">{unread}</span>}
                            </Link>

                            <div className="nav-user">
                                <div className="nav-avatar">
                                    {user.avatar_url ? (
                                        <img src={user.avatar_url} alt={user.name} />
                                    ) : (
                                        <User size={18} />
                                    )}
                                </div>
                                <div className="nav-dropdown">
                                    <div className="nav-dropdown-header">
                                        <strong>{user.name}</strong>
                                        <span>{user.email}</span>
                                    </div>
                                    <Link to="/profile" onClick={() => setMenuOpen(false)}>Profile</Link>
                                    <button onClick={handleLogout}><LogOut size={14} /> Logout</button>
                                </div>
                            </div>
                        </>
                    ) : (
                        <div className="nav-auth">
                            <Link to="/login" className="btn btn-secondary btn-sm" onClick={() => setMenuOpen(false)}>Login</Link>
                            <Link to="/register" className="btn btn-primary btn-sm" onClick={() => setMenuOpen(false)}>Sign Up</Link>
                        </div>
                    )}
                </div>

                <button className="navbar-toggle" onClick={() => setMenuOpen(!menuOpen)}>
                    {menuOpen ? <X size={24} /> : <Menu size={24} />}
                </button>
            </div>
        </nav>
    );
}
