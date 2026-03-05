import { useState, useEffect } from 'react';
import { Calendar, DollarSign, Users, BarChart3, Scissors, TrendingUp, Clock, Star, Plus, Settings, CalendarOff, Trash2, AlertTriangle } from 'lucide-react';
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
    const [createForm, setCreateForm] = useState({ name: '', address: '', city: '', phone: '', opening_time: '09:00', closing_time: '21:00', image: null });

    const [showEdit, setShowEdit] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', address: '', city: '', phone: '', opening_time: '09:00', closing_time: '21:00', description: '', image: null });

    // Services state
    const [services, setServices] = useState([]);
    const [showAddService, setShowAddService] = useState(false);
    const [serviceForm, setServiceForm] = useState({ name: '', price: '', duration_minutes: 30, description: 'Standard service' });

    // Staff state
    const [staff, setStaff] = useState([]);
    const [showAddStaff, setShowAddStaff] = useState(false);
    const [staffForm, setStaffForm] = useState({ name: '', role: 'Stylist', specialization: '', bio: '', service_ids: [] });

    // Waitlist state
    const [waitlist, setWaitlist] = useState([]);

    // Adjustment state
    const [adjustingAppt, setAdjustingAppt] = useState(null);
    const [newEndTime, setNewEndTime] = useState('');

    // Staff Hours state
    const [editingStaffHours, setEditingStaffHours] = useState(null);
    const [staffHoursData, setStaffHoursData] = useState([]);

    // Closures state
    const [closures, setClosures] = useState([]);
    const [showAddClosure, setShowAddClosure] = useState(false);
    const [closureForm, setClosureForm] = useState({ start_date: '', end_date: '', start_time: '', end_time: '', reason: '' });
    const [closureConflicts, setClosureConflicts] = useState(null);

    useEffect(() => {
        api.get('/dashboard/salons').then(res => {
            const data = res.data || [];
            setSalons(data);
            if (data.length > 0) setSelectedSalon(data[0].id);
        }).catch(() => { }).finally(() => setLoading(false));
    }, []);

    useEffect(() => {
        api.get(`/dashboard/salons/${selectedSalon}/analytics`).then(res => setAnalytics(res.data)).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/appointments`).then(res => setAppointments(res.data || [])).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/services`).then(res => setServices(res.data || [])).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/staff`).then(res => setStaff(res.data || [])).catch(() => { });
        api.get(`/dashboard/salons/${selectedSalon}/waitlist`).then(res => setWaitlist(res.data || [])).catch(() => { });
        if (tab === 'closures') {
            api.get(`/dashboard/salons/${selectedSalon}/closures`).then(res => setClosures(res.data || [])).catch(() => { });
        }
    }, [selectedSalon, tab]);

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

    const approveAppt = async (id) => {
        try {
            await api.put(`/dashboard/salons/${selectedSalon}/appointments/${id}/approve`);
            toast.success('Appointment approved');
            const res = await api.get(`/dashboard/salons/${selectedSalon}/appointments`);
            setAppointments(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed to approve'); }
    };

    const cancelAppt = async (id) => {
        if (!window.confirm('Are you sure you want to cancel this appointment?')) return;
        try {
            await api.put(`/appointments/${id}/cancel`);
            toast.success('Appointment cancelled');
            const res = await api.get(`/dashboard/salons/${selectedSalon}/appointments`);
            setAppointments(res.data || []);
        } catch (err) { toast.error(err.response?.data?.error || 'Failed to cancel'); }
    };

    const handleAdjustTime = async (e) => {
        e.preventDefault();
        try {
            await api.put(`/dashboard/salons/${selectedSalon}/appointments/${adjustingAppt.id}/adjust-time`, {
                new_end_time: newEndTime
            });
            toast.success('Schedule adjusted! Subsequent appointments shifted.');
            setAdjustingAppt(null);
            // Refresh
            const res = await api.get(`/dashboard/salons/${selectedSalon}/appointments`);
            setAppointments(res.data || []);
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to adjust schedule');
        }
    };

    const fetchStaffHours = async (staffId) => {
        try {
            const res = await api.get(`/dashboard/salons/${selectedSalon}/staff/${staffId}/hours`);
            // Ensure 7 days are represented
            const days = [0, 1, 2, 3, 4, 5, 6];
            const hours = days.map(d => {
                const existing = (res.data || []).find(h => h.day_of_week === d);
                return existing || { day_of_week: d, start_time: '09:00', end_time: '17:00', is_off: true };
            });
            setStaffHoursData(hours);
            setEditingStaffHours(staff.find(s => s.id === staffId));
        } catch (err) {
            toast.error('Failed to fetch working hours');
        }
    };

    const handleUpdateStaffHours = async (e) => {
        e.preventDefault();
        try {
            await api.put(`/dashboard/salons/${selectedSalon}/staff/${editingStaffHours.id}/hours`, {
                hours: staffHoursData
            });
            toast.success('Working hours updated');
            setEditingStaffHours(null);
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to update hours');
        }
    };

    const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

    const toggleDayOff = (index) => {
        const newData = [...staffHoursData];
        newData[index] = { ...newData[index], is_off: !newData[index].is_off };
        setStaffHoursData(newData);
    };

    const updateDayTime = (index, field, value) => {
        const newData = [...staffHoursData];
        newData[index] = { ...newData[index], [field]: value };
        setStaffHoursData(newData);
    };

    if (loading) return <div className="loading"><div className="spinner" /><p>Loading dashboard...</p></div>;

    const handleCreateSalon = async (e) => {
        e.preventDefault();
        try {
            const formData = new FormData();
            formData.append('name', createForm.name);
            formData.append('address', createForm.address);
            formData.append('city', createForm.city);
            if (createForm.phone) formData.append('phone', createForm.phone);
            formData.append('opening_time', createForm.opening_time);
            formData.append('closing_time', createForm.closing_time);
            if (createForm.image) {
                formData.append('image', createForm.image);
            }

            await api.post('/dashboard/salons', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                }
            });
            toast.success('Salon created successfully!');
            window.location.reload();
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to create salon');
        }
    };

    const handleUpdateSalon = async (e) => {
        e.preventDefault();
        try {
            const formData = new FormData();
            formData.append('name', editForm.name);
            formData.append('address', editForm.address);
            formData.append('city', editForm.city);
            if (editForm.phone) formData.append('phone', editForm.phone);
            if (editForm.description) formData.append('description', editForm.description);
            formData.append('opening_time', editForm.opening_time);
            formData.append('closing_time', editForm.closing_time);
            if (editForm.image) {
                formData.append('image', editForm.image);
            }

            await api.put(`/dashboard/salons/${selectedSalon}`, formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });
            toast.success('Salon updated successfully!');
            setShowEdit(false);
            const res = await api.get('/dashboard/salons');
            setSalons(res.data || []);
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to update salon');
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
                            <div className="grid-2">
                                <div className="form-group">
                                    <label>City</label>
                                    <input type="text" className="form-control" value={createForm.city} onChange={e => setCreateForm({ ...createForm, city: e.target.value })} required placeholder="New York" />
                                </div>
                                <div className="form-group">
                                    <label>Phone</label>
                                    <input type="tel" className="form-control" value={createForm.phone} onChange={e => setCreateForm({ ...createForm, phone: e.target.value })} placeholder="+1 555-0123" />
                                </div>
                                <div className="form-group">
                                    <label>Opening Time</label>
                                    <input type="time" className="form-control" value={createForm.opening_time} onChange={e => setCreateForm({ ...createForm, opening_time: e.target.value })} required />
                                </div>
                                <div className="form-group">
                                    <label>Closing Time</label>
                                    <input type="time" className="form-control" value={createForm.closing_time} onChange={e => setCreateForm({ ...createForm, closing_time: e.target.value })} required />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Image</label>
                                <input type="file" accept="image/*" className="form-control" onChange={e => setCreateForm({ ...createForm, image: e.target.files[0] })} />
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
                    <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
                        <select className="form-control" style={{ width: 'auto', minWidth: 200 }} value={selectedSalon} onChange={(e) => setSelectedSalon(e.target.value)}>
                            {salons.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                        </select>
                        <button className="btn btn-secondary btn-sm" onClick={() => {
                            const currentSalon = salons.find(s => s.id === selectedSalon);
                            setEditForm({
                                name: currentSalon?.name || '',
                                address: currentSalon?.address || '',
                                city: currentSalon?.city || '',
                                phone: currentSalon?.phone || '',
                                opening_time: currentSalon?.opening_time || '09:00',
                                closing_time: currentSalon?.closing_time || '21:00',
                                description: currentSalon?.description || '',
                                image: null
                            });
                            setShowEdit(true);
                        }}>
                            <Settings size={16} style={{ marginRight: 6 }} /> Edit Details
                        </button>
                    </div>
                </div>

                {/* Edit Salon Modal */}
                {showEdit && (
                    <div className="glass" style={{ padding: 24, marginBottom: 32, borderRadius: 'var(--radius-lg)' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
                            <h2>Edit Salon Details</h2>
                            <button type="button" className="btn btn-secondary btn-sm" onClick={() => setShowEdit(false)}>Close</button>
                        </div>
                        <form onSubmit={handleUpdateSalon}>
                            <div className="form-group">
                                <label>Salon Name</label>
                                <input type="text" className="form-control" value={editForm.name} onChange={e => setEditForm({ ...editForm, name: e.target.value })} required />
                            </div>
                            <div className="form-group">
                                <label>Address</label>
                                <input type="text" className="form-control" value={editForm.address} onChange={e => setEditForm({ ...editForm, address: e.target.value })} required />
                            </div>
                            <div className="grid-2">
                                <div className="form-group">
                                    <label>City</label>
                                    <input type="text" className="form-control" value={editForm.city} onChange={e => setEditForm({ ...editForm, city: e.target.value })} required />
                                </div>
                                <div className="form-group">
                                    <label>Phone</label>
                                    <input type="tel" className="form-control" value={editForm.phone} onChange={e => setEditForm({ ...editForm, phone: e.target.value })} />
                                </div>
                                <div className="form-group">
                                    <label>Opening Time</label>
                                    <input type="time" className="form-control" value={editForm.opening_time} onChange={e => setEditForm({ ...editForm, opening_time: e.target.value })} required />
                                </div>
                                <div className="form-group">
                                    <label>Closing Time</label>
                                    <input type="time" className="form-control" value={editForm.closing_time} onChange={e => setEditForm({ ...editForm, closing_time: e.target.value })} required />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Description</label>
                                <textarea className="form-control" value={editForm.description} onChange={e => setEditForm({ ...editForm, description: e.target.value })} rows="3"></textarea>
                            </div>
                            <div className="form-group">
                                <label>Update Image (Optional)</label>
                                <input type="file" accept="image/*" className="form-control" onChange={e => setEditForm({ ...editForm, image: e.target.files[0] })} />
                            </div>
                            <button type="submit" className="btn btn-primary" style={{ width: '100%', marginTop: 16 }}>Save Changes</button>
                        </form>
                    </div>
                )}

                {/* Adjust Time Modal */}
                {adjustingAppt && (
                    <div className="glass modal-overlay" style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)' }}>
                        <div className="glass" style={{ width: 400, padding: 32, borderRadius: 'var(--radius-lg)', background: 'var(--bg-primary)' }}>
                            <h2 style={{ marginBottom: 8 }}>Adjust Schedule</h2>
                            <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', marginBottom: 24 }}>
                                Update end time for <strong>{adjustingAppt.customer_name}</strong>. Subsequent appointments will shift automatically with a 5-min gap.
                            </p>
                            <form onSubmit={handleAdjustTime}>
                                <div className="form-group">
                                    <label>Current Start Time</label>
                                    <input type="text" className="form-control" value={adjustingAppt.start_time.substring(0, 5)} disabled />
                                </div>
                                <div className="form-group">
                                    <label>New End Time</label>
                                    <input
                                        type="time"
                                        className="form-control"
                                        value={newEndTime}
                                        onChange={e => setNewEndTime(e.target.value)}
                                        required
                                    />
                                </div>
                                <div style={{ display: 'flex', gap: 12, marginTop: 32 }}>
                                    <button type="button" className="btn btn-secondary" style={{ flex: 1 }} onClick={() => setAdjustingAppt(null)}>Cancel</button>
                                    <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Apply Shift</button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Staff Hours Modal */}
                {editingStaffHours && (
                    <div className="glass modal-overlay" style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)' }}>
                        <div className="glass" style={{ width: 500, padding: 32, borderRadius: 'var(--radius-lg)', background: 'var(--bg-primary)', maxHeight: '90vh', overflowY: 'auto' }}>
                            <h2 style={{ marginBottom: 8 }}>Adjust Working Hours</h2>
                            <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', marginBottom: 24 }}>
                                Set the weekly schedule for <strong>{editingStaffHours.name}</strong>.
                            </p>
                            <form onSubmit={handleUpdateStaffHours}>
                                <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
                                    {staffHoursData.map((day, idx) => (
                                        <div key={idx} style={{ display: 'flex', alignItems: 'center', gap: 12, padding: '12px 16px', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)' }}>
                                            <div style={{ width: 100, fontWeight: 600 }}>{dayNames[idx]}</div>
                                            <div style={{ display: 'flex', alignItems: 'center', gap: 8, flex: 1 }}>
                                                <input
                                                    type="time"
                                                    className="form-control"
                                                    style={{ padding: '4px 8px' }}
                                                    value={day.start_time.substring(0, 5)}
                                                    onChange={e => updateDayTime(idx, 'start_time', e.target.value)}
                                                    disabled={day.is_off}
                                                />
                                                <span style={{ opacity: 0.5 }}>-</span>
                                                <input
                                                    type="time"
                                                    className="form-control"
                                                    style={{ padding: '4px 8px' }}
                                                    value={day.end_time.substring(0, 5)}
                                                    onChange={e => updateDayTime(idx, 'end_time', e.target.value)}
                                                    disabled={day.is_off}
                                                />
                                            </div>
                                            <button
                                                type="button"
                                                className={`btn btn-sm ${day.is_off ? 'btn-danger' : 'btn-success'}`}
                                                style={{ width: 80 }}
                                                onClick={() => toggleDayOff(idx)}
                                            >
                                                {day.is_off ? 'Off' : 'Active'}
                                            </button>
                                        </div>
                                    ))}
                                </div>
                                <div style={{ display: 'flex', gap: 12, marginTop: 32 }}>
                                    <button type="button" className="btn btn-secondary" style={{ flex: 1 }} onClick={() => setEditingStaffHours(null)}>Cancel</button>
                                    <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Save Schedule</button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Dashboard Tabs */}
                <div className="tabs" style={{ marginBottom: 32 }}>
                    {['overview', 'appointments', 'staff', 'services', 'waitlist', 'closures'].map(t => (
                        <button key={t} className={`tab ${tab === t ? 'active' : ''}`} onClick={() => {
                            setTab(t);
                            if (t === 'closures' && selectedSalon) {
                                api.get(`/dashboard/salons/${selectedSalon}/closures`).then(res => setClosures(res.data || [])).catch(() => { });
                            }
                        }}>
                            {t.charAt(0).toUpperCase() + t.slice(1)}
                            {t === 'waitlist' && waitlist.length > 0 && <span className="badge-dot" />}
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
                                            <th>Phone</th>
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
                                                <td>{a.customer_phone || '-'}</td>
                                                <td>{a.service_name}</td>
                                                <td>{a.staff_name}</td>
                                                <td>{a.appointment_date}</td>
                                                <td>{a.start_time}</td>
                                                <td><span className={`badge ${statusColor[a.status]}`}>{a.status}</span></td>
                                                <td>
                                                    <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                                                        {a.status === 'pending' && (
                                                            <>
                                                                <button className="btn btn-sm btn-success" onClick={() => approveAppt(a.id)}>Approve</button>
                                                                <button className="btn btn-sm btn-danger" onClick={() => cancelAppt(a.id)}>Cancel</button>
                                                            </>
                                                        )}
                                                        {a.status === 'confirmed' && (
                                                            <>
                                                                <button className="btn btn-sm btn-primary" onClick={() => completeAppt(a.id)}>Complete</button>
                                                                <button className="btn btn-sm btn-danger" onClick={() => cancelAppt(a.id)}>Cancel</button>
                                                                <button className="btn btn-sm btn-secondary" onClick={() => markNoShow(a.id)}>No-show</button>
                                                                <button className="btn btn-sm btn-info" style={{ background: 'rgba(6,182,212,0.1)', color: 'var(--accent)', border: '1px solid var(--accent)' }} onClick={() => {
                                                                    setAdjustingAppt(a);
                                                                    setNewEndTime(a.end_time.substring(0, 5));
                                                                }}>Adjust Time</button>
                                                            </>
                                                        )}
                                                        {(a.status === 'cancelled' || a.status === 'no_show') && (
                                                            <div className="fill-opportunity">
                                                                <span style={{ fontSize: '0.75rem', color: 'var(--success)', fontWeight: 600 }}>Fill Opportunity! ✨</span>
                                                            </div>
                                                        )}
                                                    </div>
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
                                                <td>
                                                    <button className="btn btn-sm btn-info" style={{ background: 'rgba(6,182,212,0.1)', color: 'var(--accent)', border: '1px solid var(--accent)' }} onClick={() => fetchStaffHours(s.id)}>
                                                        Adjust Hours
                                                    </button>
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

                {/* Waitlist Tab */}
                {tab === 'waitlist' && (
                    <div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
                            <h2>Waitlist</h2>
                            <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>Customers waiting for an opening</p>
                        </div>
                        {waitlist.length === 0 ? (
                            <div className="empty-state"><h3>Waitlist is empty</h3></div>
                        ) : (
                            <div className="table-wrap glass" style={{ borderRadius: 'var(--radius-lg)', overflow: 'hidden' }}>
                                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                    <thead>
                                        <tr style={{ borderBottom: '1px solid var(--border)' }}>
                                            <th>Customer</th>
                                            <th>Service</th>
                                            <th>Preferred Date</th>
                                            <th>Requested</th>
                                            <th>Status</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {waitlist.map(w => (
                                            <tr key={w.id} style={{ borderBottom: '1px solid var(--border)' }}>
                                                <td>{w.customer_name}</td>
                                                <td>{w.service_name}</td>
                                                <td>{w.preferred_date}</td>
                                                <td>{new Date(w.created_at).toLocaleDateString()}</td>
                                                <td><span className="badge badge-warning">Waiting</span></td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                )}

                {/* Closures & Holidays Tab */}
                {tab === 'closures' && (
                    <div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
                            <h2>Closures & Holidays</h2>
                            <button className="btn btn-primary btn-sm" onClick={() => setShowAddClosure(!showAddClosure)}>
                                <Plus size={16} /> Block Dates
                            </button>
                        </div>

                        {showAddClosure && (
                            <div className="glass" style={{ padding: 24, marginBottom: 24, borderRadius: 'var(--radius-lg)' }}>
                                <h3 style={{ marginBottom: 16 }}>Add Closure / Holiday</h3>
                                <form onSubmit={async (e) => {
                                    e.preventDefault();
                                    try {
                                        const body = {
                                            start_date: closureForm.start_date,
                                            end_date: closureForm.end_date || closureForm.start_date,
                                            reason: closureForm.reason,
                                        };
                                        if (closureForm.start_time && closureForm.end_time) {
                                            body.start_time = closureForm.start_time;
                                            body.end_time = closureForm.end_time;
                                        }
                                        const res = await api.post(`/dashboard/salons/${selectedSalon}/closures`, body);
                                        toast.success('Closure created!');
                                        setClosures(prev => [...prev, res.data.closure]);
                                        if (res.data.conflicting_count > 0) {
                                            setClosureConflicts(res.data);
                                        } else {
                                            setShowAddClosure(false);
                                            setClosureForm({ start_date: '', end_date: '', start_time: '', end_time: '', reason: '' });
                                        }
                                    } catch (err) {
                                        toast.error(err.response?.data?.error || 'Failed to create closure');
                                    }
                                }}>
                                    <div className="grid-2">
                                        <div className="form-group">
                                            <label>Start Date</label>
                                            <input type="date" className="form-control" value={closureForm.start_date} onChange={e => setClosureForm({ ...closureForm, start_date: e.target.value })} required />
                                        </div>
                                        <div className="form-group">
                                            <label>End Date</label>
                                            <input type="date" className="form-control" value={closureForm.end_date} onChange={e => setClosureForm({ ...closureForm, end_date: e.target.value })} />
                                        </div>
                                    </div>
                                    <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginBottom: 12 }}>Leave times empty for a full-day closure.</p>
                                    <div className="grid-2">
                                        <div className="form-group">
                                            <label>Start Time (optional)</label>
                                            <input type="time" className="form-control" value={closureForm.start_time} onChange={e => setClosureForm({ ...closureForm, start_time: e.target.value })} />
                                        </div>
                                        <div className="form-group">
                                            <label>End Time (optional)</label>
                                            <input type="time" className="form-control" value={closureForm.end_time} onChange={e => setClosureForm({ ...closureForm, end_time: e.target.value })} />
                                        </div>
                                    </div>
                                    <div className="form-group">
                                        <label>Reason</label>
                                        <input type="text" className="form-control" placeholder="e.g. Vacation, Christmas, Emergency" value={closureForm.reason} onChange={e => setClosureForm({ ...closureForm, reason: e.target.value })} required />
                                    </div>
                                    <div style={{ display: 'flex', gap: 12, marginTop: 16 }}>
                                        <button type="button" className="btn btn-secondary" onClick={() => { setShowAddClosure(false); setClosureConflicts(null); }}>Cancel</button>
                                        <button type="submit" className="btn btn-primary">Save Closure</button>
                                    </div>
                                </form>

                                {closureConflicts && closureConflicts.conflicting_count > 0 && (
                                    <div style={{ marginTop: 24, padding: 20, background: 'rgba(244,63,94,0.08)', borderRadius: 'var(--radius-md)', border: '1px solid rgba(244,63,94,0.2)' }}>
                                        <h4 style={{ color: 'var(--secondary)', display: 'flex', alignItems: 'center', gap: 8 }}>
                                            <AlertTriangle size={18} /> {closureConflicts.conflicting_count} Existing Appointment(s) Affected
                                        </h4>
                                        <div style={{ maxHeight: 200, overflowY: 'auto', marginTop: 12, marginBottom: 16 }}>
                                            {closureConflicts.conflicting_bookings.map(b => (
                                                <div key={b.id} style={{ padding: '8px 0', borderBottom: '1px solid var(--border)', fontSize: '0.85rem' }}>
                                                    <strong>{b.customer_name}</strong> — {b.service_name} on {b.appointment_date} at {b.start_time}
                                                </div>
                                            ))}
                                        </div>
                                        <button className="btn btn-secondary btn-sm" style={{ marginRight: 8 }} onClick={async () => {
                                            try {
                                                await api.post(`/dashboard/salons/${selectedSalon}/closures/${closureConflicts.closure.id}/cancel-appointments`);
                                                toast.success('All conflicting appointments cancelled!');
                                                setClosureConflicts(null);
                                                setShowAddClosure(false);
                                                setClosureForm({ start_date: '', end_date: '', start_time: '', end_time: '', reason: '' });
                                            } catch (err) {
                                                toast.error('Failed to cancel appointments');
                                            }
                                        }}>
                                            Cancel All & Notify
                                        </button>
                                        <button className="btn btn-outline btn-sm" onClick={() => { setClosureConflicts(null); setShowAddClosure(false); setClosureForm({ start_date: '', end_date: '', start_time: '', end_time: '', reason: '' }); }}>
                                            Ignore for Now
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}

                        {closures.length === 0 ? (
                            <div className="empty-state" style={{ textAlign: 'center', padding: '60px 20px' }}>
                                <CalendarOff size={48} style={{ opacity: 0.3, marginBottom: 12 }} />
                                <h3>No Closures Scheduled</h3>
                                <p style={{ color: 'var(--text-secondary)' }}>Your salon is currently available all days.</p>
                            </div>
                        ) : (
                            <div className="table-wrap glass" style={{ borderRadius: 'var(--radius-lg)', overflow: 'hidden' }}>
                                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                    <thead>
                                        <tr style={{ borderBottom: '1px solid var(--border)' }}>
                                            <th>Dates</th>
                                            <th>Time</th>
                                            <th>Reason</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {closures.map(cl => (
                                            <tr key={cl.id} style={{ borderBottom: '1px solid var(--border)' }}>
                                                <td>{cl.start_date === cl.end_date ? cl.start_date : `${cl.start_date} → ${cl.end_date}`}</td>
                                                <td>{cl.start_time && cl.end_time ? `${cl.start_time} – ${cl.end_time}` : 'Full Day'}</td>
                                                <td>{cl.reason}</td>
                                                <td>
                                                    <button className="btn btn-sm" style={{ background: 'rgba(244,63,94,0.1)', color: 'var(--secondary)', border: 'none' }} onClick={async () => {
                                                        try {
                                                            await api.delete(`/dashboard/salons/${selectedSalon}/closures/${cl.id}`);
                                                            setClosures(prev => prev.filter(c => c.id !== cl.id));
                                                            toast.success('Closure removed');
                                                        } catch (err) {
                                                            toast.error('Failed to remove');
                                                        }
                                                    }}>
                                                        <Trash2 size={14} /> Remove
                                                    </button>
                                                </td>
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
        .badge-dot { width: 8px; height: 8px; background: var(--secondary); border-radius: 50%; display: inline-block; margin-left: 6px; }
        .fill-opportunity { animation: pulse 2s infinite; padding: 4px 8px; background: rgba(16,185,129,0.1); border-radius: 4px; }
        @keyframes pulse { 0% { opacity: 0.6; } 50% { opacity: 1; } 100% { opacity: 0.6; } }
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
