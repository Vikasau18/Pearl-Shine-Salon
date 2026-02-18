import { useState, useEffect } from 'react';
import { Calendar, DollarSign, Users, BarChart3, Scissors, TrendingUp, Clock, Star, Plus, Settings } from 'lucide-react';
import api from '../api/client';
import toast from 'react-hot-toast';

export default function Dashboard() {
    const [salons, setSalons] = useState([]);
    const [selectedSalon, setSelectedSalon] = useState('');
    const [analytics, setAnalytics] = useState(null);
    const [appointments, setAppointments] = useState([]);
    const [tab, setTab] = useState('overview');
    const [loading, setLoading] = useState(true);
    const [showCreate, setShowCreate] = useState(false);
    const [createForm, setCreateForm] = useState({ name: '', address: '', city: '', phone: '', image_url: 'https://images.unsplash.com/photo-1560066984-138dadb4c035?w=800' });

    // Services state
    const [services, setServices] = useState([]);
    const [showAddService, setShowAddService] = useState(false);
    const [serviceForm, setServiceForm] = useState({ name: '', price: '', duration_minutes: 30, description: 'Standard service' });

    // Staff state
    const [staff, setStaff] = useState([]);
    const [showAddStaff, setShowAddStaff] = useState(false);
    const [staffForm, setStaffForm] = useState({ name: '', role: 'Stylist', specialization: '', bio: '', service_ids: [] });

    useEffect(() => {
        api.get('/dashboard/salons').then(res => {
            const data = res.data || [];
            setSalons(data);
            if (data.length > 0) setSelectedSalon(data[0].id);
        }).catch(() => { }).finally(() => setLoading(false));
    }, []);

    useEffect(() => {
        if (!selectedSalon) return;
        api.get(`/dashboard/salons/${selectedSalon}/analytics`).then(res => setAnalytics(res.data)).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/appointments`).then(res => setAppointments(res.data || [])).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/services`).then(res => setServices(res.data || [])).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/staff`).then(res => setStaff(res.data || [])).catch(() => { });
    }, [selectedSalon]);

    // ... existing markNoShow / completeAppt ...

    const handleCreateService = async (e) => {
        e.preventDefault();
        try {
            await api.post(`/dashboard/salons/${selectedSalon}/services`, {
                ...serviceForm,
                price: parseFloat(serviceForm.price),
                duration_minutes: parseInt(serviceForm.duration_minutes)
            });
            toast.success('Service added');
            setShowAddService(false);
            setServiceForm({ name: '', price: '', duration_minutes: 30, description: 'Standard service' });
            // Refresh
            const res = await api.get(`/dashboard/salons/${selectedSalon}/services`);
            setServices(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed to add service'); }
    };

    const handleCreateStaff = async (e) => {
        e.preventDefault();
        try {
            await api.post(`/dashboard/salons/${selectedSalon}/staff`, staffForm);
            toast.success('Staff added');
            setShowAddStaff(false);
            setStaffForm({ name: '', role: 'Stylist', specialization: '', bio: '', service_ids: [] });
            // Refresh
            const res = await api.get(`/dashboard/salons/${selectedSalon}/staff`);
            setStaff(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed to add staff'); }
    };

    const toggleServiceSelection = (id) => {
        const current = staffForm.service_ids;
        if (current.includes(id)) {
            setStaffForm({ ...staffForm, service_ids: current.filter(sid => sid !== id) });
        } else {
            setStaffForm({ ...staffForm, service_ids: [...current, id] });
        }
    };

    const markNoShow = async (id) => {
        try {
            await api.put(`/dashboard/salons/${selectedSalon}/appointments/${id}/no-show`);
            toast.success('Marked as no-show');
            const res = await api.get(`/dashboard/salons/${selectedSalon}/appointments`);
            setAppointments(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed'); }
    };

    const completeAppt = async (id) => {
        try {
            await api.put(`/dashboard/salons/${selectedSalon}/appointments/${id}/complete`);
            toast.success('Marked as completed');
            const res = await api.get(`/dashboard/salons/${selectedSalon}/appointments`);
            setAppointments(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed'); }
    };

    if (loading) return <div className="loading"><div className="spinner" /><p>Loading dashboard...</p></div>;

    const handleCreateSalon = async (e) => {
        e.preventDefault();
        try {
            await api.post('/dashboard/salons', createForm);
            toast.success('Salon created successfully!');
            window.location.reload();
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to create salon');
        }
    };

    if (salons.length === 0 || showCreate) {
        return (
            <div className="page">
                <div className="container" style={{ maxWidth: 600 }}>
                    <div className="glass" style={{ padding: 40, borderRadius: 'var(--radius-xl)' }}>
                        <div style={{ textAlign: 'center', marginBottom: 32 }}>
                            <Scissors size={48} color="var(--primary)" style={{ marginBottom: 16 }} />
                            <h1>Create Your Salon</h1>
                            <p style={{ color: 'var(--text-secondary)' }}>Enter your business details to get started</p>
                        </div>
                        <form onSubmit={handleCreateSalon}>
                            <div className="form-group">
                                <label>Salon Name</label>
                                <input type="text" className="form-control" value={createForm.name} onChange={e => setCreateForm({ ...createForm, name: e.target.value })} required placeholder="e.g. Luxe Beauty Bar" />
                            </div>
                            <div className="form-group">
                                <label>Address</label>
                                <input type="text" className="form-control" value={createForm.address} onChange={e => setCreateForm({ ...createForm, address: e.target.value })} required placeholder="123 Main St" />
                            </div>
                            <div className="form-group">
                                <label>City</label>
                                <input type="text" className="form-control" value={createForm.city} onChange={e => setCreateForm({ ...createForm, city: e.target.value })} required placeholder="New York" />
                            </div>
                            <div className="form-group">
                                <label>Phone</label>
                                <input type="tel" className="form-control" value={createForm.phone} onChange={e => setCreateForm({ ...createForm, phone: e.target.value })} placeholder="+1 555-0123" />
                            </div>
                            <div className="form-group">
                                <label>Image URL (Optional)</label>
                                <input type="url" className="form-control" value={createForm.image_url} onChange={e => setCreateForm({ ...createForm, image_url: e.target.value })} placeholder="https://..." />
                            </div>
                            <div style={{ display: 'flex', gap: 12, marginTop: 32 }}>
                                {salons.length > 0 && <button type="button" className="btn btn-secondary" style={{ flex: 1 }} onClick={() => setShowCreate(false)}>Cancel</button>}
                                <button type="submit" className="btn btn-primary btn-lg" style={{ flex: 1 }}>Create Salon</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        );
    }

    const statusColor = {
        confirmed: 'badge-success', pending: 'badge-warning',
        cancelled: 'badge-danger', completed: 'badge-primary', no_show: 'badge-danger',
    };

    return (
        <div className="page dashboard">
            <div className="container">
                {/* Header */}
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 32, flexWrap: 'wrap', gap: 16 }}>
                    <div>
                        <h1>Salon <span className="gradient-text">Dashboard</span></h1>
                        <p style={{ color: 'var(--text-secondary)' }}>Manage your salon operations</p>
                    </div>
                    <select className="form-control" style={{ width: 'auto', minWidth: 200 }} value={selectedSalon} onChange={(e) => setSelectedSalon(e.target.value)}>
                        {salons.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                    </select>
                </div>

                {/* Dashboard Tabs */}
                <div className="tabs" style={{ marginBottom: 32 }}>
                    {['overview', 'appointments', 'staff', 'services'].map(t => (
                        <button key={t} className={`tab ${tab === t ? 'active' : ''}`} onClick={() => setTab(t)}>
                            {t.charAt(0).toUpperCase() + t.slice(1)}
                        </button>
                    ))}
                </div>

                {/* Overview Tab */}
                {tab === 'overview' && analytics && (
                    <div>
                        <div className="grid-4">
                            <div className="stat-card">
                                <div className="stat-icon" style={{ background: 'rgba(124,58,237,0.15)', color: 'var(--primary)' }}>
                                    <Calendar size={22} />
                                </div>
                                <div className="stat-value">{analytics.total_appointments || 0}</div>
                                <div className="stat-label">Total Appointments</div>
                            </div>
                            <div className="stat-card">
                                <div className="stat-icon" style={{ background: 'rgba(16,185,129,0.15)', color: 'var(--success)' }}>
                                    <DollarSign size={22} />
                                </div>
                                <div className="stat-value">${(analytics.total_revenue || 0).toFixed(0)}</div>
                                <div className="stat-label">Total Revenue</div>
                            </div>
                            <div className="stat-card">
                                <div className="stat-icon" style={{ background: 'rgba(6,182,212,0.15)', color: 'var(--accent)' }}>
                                    <Users size={22} />
                                </div>
                                <div className="stat-value">{analytics.total_customers || 0}</div>
                                <div className="stat-label">Customers</div>
                            </div>
                            <div className="stat-card">
                                <div className="stat-icon" style={{ background: 'rgba(244,63,94,0.15)', color: 'var(--secondary)' }}>
                                    <Star size={22} />
                                </div>
                                <div className="stat-value">{(analytics.avg_rating || 0).toFixed(1)}</div>
                                <div className="stat-label">Avg Rating</div>
                            </div>
                        </div>

                        {/* Popular services & Peak hours */}
                        <div className="grid-2" style={{ marginTop: 32 }}>
                            <div className="glass" style={{ padding: 24, borderRadius: 'var(--radius-lg)' }}>
                                <h3 style={{ marginBottom: 20, display: 'flex', alignItems: 'center', gap: 8 }}>
                                    <TrendingUp size={18} /> Popular Services
                                </h3>
                                {(analytics.popular_services || []).map((svc, i) => (
                                    <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 0', borderBottom: '1px solid var(--border)' }}>
                                        <span>{svc.name}</span>
                                        <span style={{ color: 'var(--primary-light)', fontWeight: 600 }}>{svc.booking_count} bookings</span>
                                    </div>
                                ))}
                                {(!analytics.popular_services || analytics.popular_services.length === 0) && <p style={{ color: 'var(--text-muted)' }}>No data yet</p>}
                            </div>
                            <div className="glass" style={{ padding: 24, borderRadius: 'var(--radius-lg)' }}>
                                <h3 style={{ marginBottom: 20, display: 'flex', alignItems: 'center', gap: 8 }}>
                                    <Clock size={18} /> Peak Hours
                                </h3>
                                {(analytics.peak_hours || []).map((h, i) => (
                                    <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 0', borderBottom: '1px solid var(--border)' }}>
                                        <span>{h.hour}:00</span>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                                            <div style={{ width: Math.min(h.count * 20, 150), height: 8, borderRadius: 4, background: 'var(--primary)' }} />
                                            <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>{h.count}</span>
                                        </div>
                                    </div>
                                ))}
                                {(!analytics.peak_hours || analytics.peak_hours.length === 0) && <p style={{ color: 'var(--text-muted)' }}>No data yet</p>}
                            </div>
                        </div>
                    </div>
                )}

                {/* Appointments Tab */}
                {tab === 'appointments' && (
                    <div>
                        <h2 style={{ marginBottom: 20 }}>Appointments</h2>
                        {appointments.length === 0 ? (
                            <div className="empty-state"><h3>No appointments</h3></div>
                        ) : (
                            <div className="table-wrap glass" style={{ borderRadius: 'var(--radius-lg)', overflow: 'hidden' }}>
                                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                    <thead>
                                        <tr style={{ borderBottom: '1px solid var(--border)' }}>
                                            <th>Customer</th>
                                            <th>Service</th>
                                            <th>Staff</th>
                                            <th>Date</th>
                                            <th>Time</th>
                                            <th>Status</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {appointments.map(a => (
                                            <tr key={a.id} style={{ borderBottom: '1px solid var(--border)' }}>
                                                <td>{a.customer_name || 'Customer'}</td>
                                                <td>{a.service_name}</td>
                                                <td>{a.staff_name}</td>
                                                <td>{a.appointment_date}</td>
                                                <td>{a.start_time}</td>
                                                <td><span className={`badge ${statusColor[a.status]}`}>{a.status}</span></td>
                                                <td>
                                                    {a.status === 'confirmed' && (
                                                        <div style={{ display: 'flex', gap: 4 }}>
                                                            <button className="btn btn-sm btn-primary" onClick={() => completeAppt(a.id)}>Complete</button>
                                                            <button className="btn btn-sm btn-danger" onClick={() => markNoShow(a.id)}>No-show</button>
                                                        </div>
                                                    )}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                )}

                {/* Staff Tab */}
                {tab === 'staff' && (
                    <div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
                            <h2>Staff</h2>
                            <button className="btn btn-primary btn-sm" onClick={() => setShowAddStaff(!showAddStaff)}>
                                {showAddStaff ? 'Cancel' : <><Plus size={16} /> Add Staff</>}
                            </button>
                        </div>

                        {showAddStaff && (
                            <div className="glass" style={{ padding: 24, marginBottom: 24, borderRadius: 'var(--radius-lg)' }}>
                                <form onSubmit={handleCreateStaff} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
                                    <div className="form-group">
                                        <label>Name</label>
                                        <input type="text" className="form-control" value={staffForm.name} onChange={e => setStaffForm({ ...staffForm, name: e.target.value })} required placeholder="e.g. Sarah Smith" />
                                    </div>
                                    <div className="form-group">
                                        <label>Role</label>
                                        <select className="form-control" value={staffForm.role} onChange={e => setStaffForm({ ...staffForm, role: e.target.value })}>
                                            <option value="Stylist">Stylist</option>
                                            <option value="Manager">Manager</option>
                                            <option value="Assistant">Assistant</option>
                                        </select>
                                    </div>
                                    <div className="form-group">
                                        <label>Specialization</label>
                                        <input type="text" className="form-control" value={staffForm.specialization} onChange={e => setStaffForm({ ...staffForm, specialization: e.target.value })} placeholder="e.g. Colorist" />
                                    </div>
                                    <div className="form-group" style={{ gridColumn: '1 / -1' }}>
                                        <label>Services Performed</label>
                                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginTop: 8 }}>
                                            {services.map(s => (
                                                <button
                                                    key={s.id}
                                                    type="button"
                                                    className={`btn btn-sm ${staffForm.service_ids.includes(s.id) ? 'btn-primary' : 'btn-secondary'}`}
                                                    onClick={() => toggleServiceSelection(s.id)}
                                                >
                                                    {s.name}
                                                </button>
                                            ))}
                                            {services.length === 0 && <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>No services available. Create services first.</p>}
                                        </div>
                                    </div>
                                    <div style={{ gridColumn: '1 / -1' }}>
                                        <button type="submit" className="btn btn-primary" style={{ width: '100%' }}>Save Staff Member</button>
                                    </div>
                                </form>
                            </div>
                        )}

                        {staff.length === 0 ? (
                            <div className="empty-state">
                                <Users size={48} style={{ marginBottom: 16, opacity: 0.3 }} />
                                <h3>No staff members</h3>
                                <p>Add staff to perform services.</p>
                            </div>
                        ) : (
                            <div className="table-wrap glass" style={{ borderRadius: 'var(--radius-lg)', overflow: 'hidden' }}>
                                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                    <thead>
                                        <tr style={{ borderBottom: '1px solid var(--border)' }}>
                                            <th>Name</th>
                                            <th>Role</th>
                                            <th>Specialization</th>
                                            <th>Details</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {staff.map(s => (
                                            <tr key={s.id} style={{ borderBottom: '1px solid var(--border)' }}>
                                                <td style={{ fontWeight: 600 }}>{s.name}</td>
                                                <td><span className="badge badge-info">{s.role}</span></td>
                                                <td>{s.specialization}</td>
                                                <td style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                                                    {(s.services || []).length} services assigned
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                )}

                {/* Services Tab */}
                {tab === 'services' && (
                    <div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
                            <h2>Services</h2>
                            <button className="btn btn-primary btn-sm" onClick={() => setShowAddService(!showAddService)}>
                                {showAddService ? 'Cancel' : <><Plus size={16} /> Add Service</>}
                            </button>
                        </div>

                        {showAddService && (
                            <div className="glass" style={{ padding: 24, marginBottom: 24, borderRadius: 'var(--radius-lg)' }}>
                                <form onSubmit={handleCreateService} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
                                    <div className="form-group">
                                        <label>Service Name</label>
                                        <input type="text" className="form-control" value={serviceForm.name} onChange={e => setServiceForm({ ...serviceForm, name: e.target.value })} required placeholder="e.g. Haircut" />
                                    </div>
                                    <div className="form-group">
                                        <label>Price ($)</label>
                                        <input type="number" className="form-control" value={serviceForm.price} onChange={e => setServiceForm({ ...serviceForm, price: e.target.value })} required placeholder="25.00" />
                                    </div>
                                    <div className="form-group">
                                        <label>Duration (min)</label>
                                        <input type="number" className="form-control" value={serviceForm.duration_minutes} onChange={e => setServiceForm({ ...serviceForm, duration_minutes: e.target.value })} required step="15" />
                                    </div>
                                    <div className="form-group">
                                        <label>Description</label>
                                        <input type="text" className="form-control" value={serviceForm.description} onChange={e => setServiceForm({ ...serviceForm, description: e.target.value })} />
                                    </div>
                                    <div style={{ gridColumn: '1 / -1' }}>
                                        <button type="submit" className="btn btn-primary" style={{ width: '100%' }}>Save Service</button>
                                    </div>
                                </form>
                            </div>
                        )}

                        {services.length === 0 ? (
                            <div className="empty-state">
                                <Scissors size={48} style={{ marginBottom: 16, opacity: 0.3 }} />
                                <h3>No services yet</h3>
                                <p>Add services so customers can book them.</p>
                            </div>
                        ) : (
                            <div className="table-wrap glass" style={{ borderRadius: 'var(--radius-lg)', overflow: 'hidden' }}>
                                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                    <thead>
                                        <tr style={{ borderBottom: '1px solid var(--border)' }}>
                                            <th>Name</th>
                                            <th>Price</th>
                                            <th>Duration</th>
                                            <th>Description</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {services.map(s => (
                                            <tr key={s.id} style={{ borderBottom: '1px solid var(--border)' }}>
                                                <td style={{ fontWeight: 600 }}>{s.name}</td>
                                                <td>${s.price}</td>
                                                <td>{s.duration_minutes} min</td>
                                                <td style={{ color: 'var(--text-secondary)' }}>{s.description}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                )}
            </div>

            <style>{`
        .dashboard table { font-size: 0.9rem; }
        .dashboard th { text-align: left; padding: 14px 16px; color: var(--text-muted); font-weight: 600; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.5px; }
        .dashboard td { padding: 14px 16px; }
        .dashboard tr:hover td { background: rgba(124,58,237,0.03); }
        .tabs { background: var(--bg-secondary); padding: 4px; border-radius: var(--radius-md); display: inline-flex; }
        .tab { padding: 10px 24px; background: none; border: none; color: var(--text-secondary); cursor: pointer; font-family: inherit; font-weight: 500; border-radius: var(--radius-sm); transition: var(--transition-fast); white-space: nowrap; }
        .tab.active { background: var(--primary); color: white; }
      `}</style>
        </div>
    );
}
