import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Calendar, Clock, MapPin, XCircle, RefreshCw, Star, ChevronDown } from 'lucide-react';
import api from '../api/client';
import toast from 'react-hot-toast';

export default function MyAppointments() {
    const [appointments, setAppointments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState('upcoming');

    const fetchAppointments = async () => {
        setLoading(true);
        try {
            const res = await api.get('/appointments', { params: { status: filter !== 'all' ? filter : undefined } });
            setAppointments(res.data || []);
        } catch { setAppointments([]); }
        finally { setLoading(false); }
    };

    useEffect(() => { fetchAppointments(); }, [filter]);

    const cancel = async (id) => {
        if (!confirm('Cancel this appointment?')) return;
        try {
            await api.put(`/appointments/${id}/cancel`);
            toast.success('Appointment cancelled');
            fetchAppointments();
        } catch (err) { toast.error(err.response?.data?.error || 'Cancel failed'); }
    };

    const statusColor = {
        confirmed: 'badge-success',
        pending: 'badge-warning',
        cancelled: 'badge-danger',
        completed: 'badge-primary',
        no_show: 'badge-danger',
    };

    return (
        <div className="page">
            <div className="container">
                <div className="page-header">
                    <h1>My <span className="gradient-text">Appointments</span></h1>
                    <p>Manage your upcoming and past bookings</p>
                </div>

                <div className="filter-bar">
                    {['upcoming', 'completed', 'cancelled', 'all'].map(f => (
                        <button key={f} className={`btn btn-sm ${filter === f ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setFilter(f)}>
                            {f.charAt(0).toUpperCase() + f.slice(1)}
                        </button>
                    ))}
                </div>

                {loading ? (
                    <div className="loading"><div className="spinner" /><p>Loading appointments...</p></div>
                ) : appointments.length === 0 ? (
                    <div className="empty-state">
                        <Calendar size={48} style={{ marginBottom: 16, opacity: 0.3 }} />
                        <h3>No appointments found</h3>
                        <p>Book your first appointment to get started!</p>
                        <Link to="/salons" className="btn btn-primary" style={{ marginTop: 16 }}>Explore Salons</Link>
                    </div>
                ) : (
                    <div className="appointments-list">
                        {appointments.map(appt => (
                            <div key={appt.id} className="appointment-card glass">
                                <div className="appt-left">
                                    <div className="appt-date-box">
                                        <span className="appt-month">{new Date(appt.appointment_date).toLocaleDateString('en', { month: 'short' })}</span>
                                        <span className="appt-day">{new Date(appt.appointment_date).getDate()}</span>
                                    </div>
                                    <div className="appt-details">
                                        <h3>{appt.service_name || 'Service'}</h3>
                                        <p><MapPin size={14} /> {appt.salon_name || 'Salon'}</p>
                                        <p><Clock size={14} /> {appt.start_time} - {appt.end_time}</p>
                                        <p><Star size={14} /> Staff: {appt.staff_name || 'N/A'}</p>
                                    </div>
                                </div>
                                <div className="appt-right">
                                    <span className={`badge ${statusColor[appt.status] || 'badge-info'}`}>{appt.status}</span>
                                    <strong className="appt-price">${(appt.total_price || 0).toFixed(2)}</strong>
                                    {(appt.status === 'confirmed' || appt.status === 'pending') && (
                                        <div className="appt-actions">
                                            <button className="btn btn-sm btn-danger" onClick={() => cancel(appt.id)}><XCircle size={14} /> Cancel</button>
                                        </div>
                                    )}
                                    {appt.status === 'completed' && !appt.reviewed && (
                                        <Link to={`/salons/${appt.salon_id}`} className="btn btn-sm btn-outline">Leave Review</Link>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <style>{`
        .filter-bar { display: flex; gap: 8px; margin-bottom: 32px; flex-wrap: wrap; }
        .appointments-list { display: flex; flex-direction: column; gap: 16px; }
        .appointment-card { display: flex; justify-content: space-between; align-items: center; padding: 24px; border-radius: var(--radius-lg); }
        .appt-left { display: flex; gap: 20px; align-items: center; }
        .appt-date-box { display: flex; flex-direction: column; align-items: center; background: var(--bg-elevated); padding: 12px 16px; border-radius: var(--radius-md); min-width: 60px; }
        .appt-month { font-size: 0.7rem; text-transform: uppercase; color: var(--primary-light); font-weight: 700; }
        .appt-day { font-size: 1.6rem; font-weight: 800; }
        .appt-details h3 { font-size: 1.05rem; margin-bottom: 4px; }
        .appt-details p { display: flex; align-items: center; gap: 6px; color: var(--text-secondary); font-size: 0.85rem; }
        .appt-right { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }
        .appt-price { font-size: 1.2rem; font-weight: 700; color: var(--primary-light); }
        .appt-actions { display: flex; gap: 8px; }
        @media (max-width: 768px) {
          .appointment-card { flex-direction: column; align-items: flex-start; gap: 16px; }
          .appt-right { align-items: flex-start; flex-direction: row; flex-wrap: wrap; }
        }
      `}</style>
        </div>
    );
}
