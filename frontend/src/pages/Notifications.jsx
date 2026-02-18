import { useState, useEffect } from 'react';
import { Bell, Check, CheckCheck } from 'lucide-react';
import api from '../api/client';

export default function Notifications() {
    const [notifications, setNotifications] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetch = async () => {
        try {
            const res = await api.get('/notifications');
            setNotifications(res.data || []);
        } catch { setNotifications([]); }
        finally { setLoading(false); }
    };

    useEffect(() => { fetch(); }, []);

    const markRead = async (id) => {
        await api.put(`/notifications/${id}/read`);
        fetch();
    };

    const markAllRead = async () => {
        await api.put('/notifications/read-all');
        fetch();
    };

    return (
        <div className="page">
            <div className="container">
                <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <div>
                        <h1>Notifications</h1>
                        <p>Stay updated on your appointments</p>
                    </div>
                    <button className="btn btn-secondary btn-sm" onClick={markAllRead}>
                        <CheckCheck size={16} /> Mark all read
                    </button>
                </div>

                {loading ? (
                    <div className="loading"><div className="spinner" /><p>Loading...</p></div>
                ) : notifications.length === 0 ? (
                    <div className="empty-state">
                        <Bell size={48} style={{ marginBottom: 16, opacity: 0.3 }} />
                        <h3>No notifications</h3>
                    </div>
                ) : (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                        {notifications.map(n => (
                            <div key={n.id} className={`notif-item glass ${n.is_read ? '' : 'unread'}`} onClick={() => !n.is_read && markRead(n.id)}>
                                <div className="notif-dot" />
                                <div className="notif-body">
                                    <strong>{n.title}</strong>
                                    <p>{n.message}</p>
                                    <span className="notif-time">{new Date(n.created_at).toLocaleString()}</span>
                                </div>
                                {!n.is_read && <Check size={16} style={{ color: 'var(--text-muted)' }} />}
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <style>{`
        .notif-item { display: flex; align-items: center; gap: 16px; padding: 20px; border-radius: var(--radius-md); cursor: pointer; transition: var(--transition-fast); }
        .notif-item:hover { border-color: var(--primary-light); }
        .notif-item.unread { border-left: 3px solid var(--primary); }
        .notif-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--text-muted); flex-shrink: 0; }
        .notif-item.unread .notif-dot { background: var(--primary); }
        .notif-body { flex: 1; }
        .notif-body strong { display: block; margin-bottom: 4px; }
        .notif-body p { color: var(--text-secondary); font-size: 0.9rem; }
        .notif-time { font-size: 0.75rem; color: var(--text-muted); }
      `}</style>
        </div>
    );
}
