import { useState, useEffect } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import { Calendar, Clock, User, CheckCircle, XCircle, CalendarOff } from 'lucide-react';
import api from '../api/client';
import toast from 'react-hot-toast';

export default function BookAppointment() {
    const { salonId } = useParams();
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const [services, setServices] = useState([]);
    const [staff, setStaff] = useState([]);
    const [selectedService, setSelectedService] = useState(searchParams.get('service') || '');
    const [selectedStaff, setSelectedStaff] = useState('');
    const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
    const [selectedTime, setSelectedTime] = useState('');
    const [slots, setSlots] = useState([]);
    const [promo, setPromo] = useState('');
    const [notes, setNotes] = useState('');
    const [discount, setDiscount] = useState(0);
    const [loading, setLoading] = useState(false);
    const [salon, setSalon] = useState(null);
    const [booked, setBooked] = useState(false);
    const [waitlistSlots, setWaitlistSlots] = useState([]);

    useEffect(() => {
        Promise.all([
            api.get(`/salons/${salonId}`),
            api.get(`/salons/${salonId}/services`),
        ]).then(([sRes, svRes]) => {
            setSalon(sRes.data);
            setServices(svRes.data);
        });
    }, [salonId]);

    useEffect(() => {
        if (selectedService) {
            api.get(`/salons/${salonId}/staff-for-service?service_id=${selectedService}`)
                .then(res => setStaff(res.data))
                .catch(() => setStaff([]));
        }
    }, [selectedService, salonId]);

    useEffect(() => {
        if (selectedService && selectedStaff && selectedDate) {
            api.get('/appointments/available-slots', {
                params: { salon_id: salonId, staff_id: selectedStaff, service_id: selectedService, date: selectedDate }
            }).then(res => {
                const availableSlots = res.data || [];
                setSlots(availableSlots);
            }).catch(() => setSlots([]));
        }
    }, [selectedService, selectedStaff, selectedDate, salonId]);

    const validatePromo = async () => {
        if (!promo) return;
        try {
            const res = await api.get(`/promos/validate?code=${promo}&salon_id=${salonId}`);
            setDiscount(res.data.discount_percent || 0);
            toast.success(`Promo applied! ${res.data.discount_percent}% off`);
        } catch {
            toast.error('Invalid promo code');
            setDiscount(0);
        }
    };

    const handleBook = async () => {
        setLoading(true);
        try {
            let apptRes = null;
            if (selectedTime) {
                apptRes = await api.post('/appointments', {
                    salon_id: salonId,
                    service_id: selectedService,
                    staff_id: selectedStaff,
                    date: selectedDate,
                    start_time: selectedTime,
                    promo_code: promo || undefined,
                    notes: notes || undefined,
                });
            }

            if (waitlistSlots.length > 0) {
                await api.post('/waitlist', {
                    salon_id: salonId,
                    service_id: selectedService,
                    staff_id: selectedStaff,
                    preferred_date: selectedDate,
                    preferred_times: waitlistSlots,
                });
            }

            if (selectedTime) {
                setBooked(true);
                toast.success('Appointment booked successfully!');
            } else if (waitlistSlots.length > 0) {
                toast.success('Joined waitlist! We will auto-assign you if a slot opens.');
                navigate('/appointments');
            }
        } catch (err) {
            toast.error(err.response?.data?.error || 'Operation failed');
        } finally {
            setLoading(false);
        }
    };

    const selectedServiceData = services.find(s => s.id === selectedService);
    const price = selectedServiceData?.price || 0;
    const finalPrice = price - (price * discount / 100);

    if (booked) {
        return (
            <div className="page">
                <div className="container" style={{ textAlign: 'center', padding: '80px 20px' }}>
                    <div style={{ width: 80, height: 80, borderRadius: '50%', background: 'rgba(16,185,129,0.15)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
                        <CheckCircle size={40} color="var(--success)" />
                    </div>
                    <h1 style={{ marginBottom: 8 }}>Booking Confirmed!</h1>
                    <p style={{ color: 'var(--text-secondary)', marginBottom: 32 }}>Your appointment has been scheduled. {waitlistSlots.length > 0 && "We've also noted your waitlist preferences."}</p>
                    <div className="glass" style={{ display: 'inline-block', padding: 32, borderRadius: 'var(--radius-lg)', textAlign: 'left' }}>
                        <p><strong>Salon:</strong> {salon?.name}</p>
                        <p><strong>Service:</strong> {selectedServiceData?.name}</p>
                        <p><strong>Date:</strong> {selectedDate}</p>
                        <p><strong>Time:</strong> {selectedTime}</p>
                        <p><strong>Total:</strong> ${finalPrice.toFixed(2)}</p>
                    </div>
                    <div style={{ marginTop: 32, display: 'flex', gap: 12, justifyContent: 'center' }}>
                        <button className="btn btn-primary" onClick={() => navigate('/appointments')}>My Bookings</button>
                        <button className="btn btn-secondary" onClick={() => navigate('/')}>Go Home</button>
                    </div>
                </div>
            </div>
        );
    }

    const isReady = selectedService && selectedStaff && selectedDate && (selectedTime || waitlistSlots.length > 0);

    return (
        <div className="page">
            <div className="container">
                <div className="page-header" style={{ marginBottom: 40 }}>
                    <h1>Book at <span className="gradient-text">{salon?.name || 'Salon'}</span></h1>
                    <p style={{ color: 'var(--text-secondary)' }}>Complete your selection to schedule your appointment</p>
                </div>

                <div className="booking-grid">
                    {/* Main Selection Area */}
                    <div className="booking-main">
                        {/* Section 1: Services */}
                        <section className="booking-section glass">
                            <h3 className="section-title"><CheckCircle size={18} color={selectedService ? 'var(--success)' : 'var(--text-muted)'} /> 1. Select Service</h3>
                            <div className="service-options">
                                {services.map(svc => (
                                    <div key={svc.id} className={`service-option ${selectedService === svc.id ? 'selected' : ''}`} onClick={() => setSelectedService(svc.id)}>
                                        <div style={{ flex: 1 }}>
                                            <h4 style={{ margin: 0 }}>{svc.name}</h4>
                                            <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{svc.duration_minutes} min</span>
                                        </div>
                                        <strong style={{ color: 'var(--primary-light)' }}>${svc.price.toFixed(2)}</strong>
                                    </div>
                                ))}
                            </div>
                        </section>

                        {/* Section 2: Stylist */}
                        <section className={`booking-section glass ${!selectedService ? 'disabled' : ''}`}>
                            <h3 className="section-title"><User size={18} color={selectedStaff ? 'var(--success)' : 'var(--text-muted)'} /> 2. Choose a Stylist</h3>
                            {!selectedService ? (
                                <p className="selection-helper">Select a service first</p>
                            ) : (
                                <div className="staff-grid">
                                    {staff.map(s => (
                                        <div key={s.id} className={`staff-card ${selectedStaff === s.id ? 'selected' : ''}`} onClick={() => setSelectedStaff(s.id)}>
                                            <div className="staff-avatar">
                                                <User size={20} />
                                            </div>
                                            <div className="staff-info">
                                                <span className="staff-name">{s.name}</span>
                                                <span className="staff-role">{s.specialization || 'Stylist'}</span>
                                            </div>
                                        </div>
                                    ))}
                                    {staff.length === 0 && <p className="selection-helper">No stylists available for this service</p>}
                                </div>
                            )}
                        </section>

                        {/* Section 3: Date & Time */}
                        <section className={`booking-section glass ${!selectedStaff ? 'disabled' : ''}`}>
                            <h3 className="section-title"><Calendar size={18} color={selectedDate ? 'var(--success)' : 'var(--text-muted)'} /> 3. Date & Time</h3>
                            {!selectedStaff ? (
                                <p className="selection-helper">Choose a stylist first</p>
                            ) : (
                                <>
                                    <div className="form-group" style={{ marginBottom: 20 }}>
                                        <input
                                            type="date"
                                            className="form-control"
                                            value={selectedDate}
                                            onChange={(e) => setSelectedDate(e.target.value)}
                                            min={new Date().toISOString().split('T')[0]}
                                        />
                                    </div>
                                    {selectedDate && (
                                        <div className="slots-container">
                                            {/* Salon Closure Section */}
                                            {slots.some(s => s.availability_reason === 'closed') && (
                                                <div className="slots-group closure-notice" style={{ background: 'rgba(239,68,68,0.1)', padding: 16, borderRadius: 'var(--radius-md)', border: '1px solid rgba(239, 68, 68, 0.2)', marginBottom: 24 }}>
                                                    <h4 style={{ color: 'var(--secondary)', display: 'flex', alignItems: 'center', gap: 8, margin: 0 }}>
                                                        <CalendarOff size={18} /> {slots.length === 1 ? "Salon is Closed Today" : "Salon Closed During These Hours"}
                                                    </h4>
                                                    {slots.length > 1 && (
                                                        <div className="time-slots" style={{ marginTop: 12 }}>
                                                            {slots.filter(s => s.availability_reason === 'closed').map(slot => (
                                                                <button key={slot.start_time} className="time-slot" disabled style={{ opacity: 0.5, cursor: 'not-allowed', background: 'rgba(239,68,68,0.1)', color: 'var(--secondary)' }}>
                                                                    {slot.start_time}
                                                                </button>
                                                            ))}
                                                        </div>
                                                    )}
                                                </div>
                                            )}

                                            {/* Available Slots Section */}
                                            {slots.some(s => s.availability_reason === 'available') ? (
                                                <div className="slots-group">
                                                    <h4 className="group-label"><CheckCircle size={14} /> Available Times</h4>
                                                    <div className="time-slots">
                                                        {slots.filter(slot => slot.availability_reason === 'available').map(slot => (
                                                            <button
                                                                key={slot.start_time}
                                                                className={`time-slot ${selectedTime === slot.start_time ? 'selected' : ''}`}
                                                                onClick={() => {
                                                                    setSelectedTime(prev => prev === slot.start_time ? '' : slot.start_time);
                                                                }}
                                                            >
                                                                {slot.start_time}
                                                            </button>
                                                        ))}
                                                    </div>
                                                </div>
                                            ) : (
                                                !slots.some(s => s.availability_reason === 'closed') && (
                                                    <div className="empty-state" style={{ padding: '20px 0' }}>
                                                        <p style={{ color: 'var(--text-secondary)' }}>No available slots for this selection.</p>
                                                    </div>
                                                )
                                            )}

                                            {/* Busy Slots Section (Waitlist) */}
                                            {slots.some(s => s.availability_reason === 'booked') && (
                                                <div className="slots-group" style={{ marginTop: 24 }}>
                                                    <h4 className="group-label"><Clock size={14} /> Busy Slots (Waitlist Only)</h4>
                                                    <p style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: 12 }}>
                                                        Select these to be auto-assigned if they open up.
                                                    </p>
                                                    <div className="time-slots">
                                                        {slots.filter(slot => slot.availability_reason === 'booked').map(slot => (
                                                            <button
                                                                key={slot.start_time}
                                                                className={`time-slot busy ${waitlistSlots.includes(slot.start_time) ? 'waitlist-selected' : ''}`}
                                                                onClick={() => {
                                                                    setWaitlistSlots(prev =>
                                                                        prev.includes(slot.start_time)
                                                                            ? prev.filter(t => t !== slot.start_time)
                                                                            : [...prev, slot.start_time]
                                                                    );
                                                                }}
                                                            >
                                                                {slot.start_time}
                                                            </button>
                                                        ))}
                                                    </div>
                                                </div>
                                            )}

                                            {/* Unserviceable Slots (Out of Shift) */}
                                            {slots.some(s => s.availability_reason === 'out_of_shift') && (
                                                <div className="slots-group" style={{ marginTop: 24 }}>
                                                    <h4 className="group-label"><XCircle size={14} /> Unserviceable</h4>
                                                    <p style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: 12 }}>
                                                        Stylist is not available during these hours.
                                                    </p>
                                                    <div className="time-slots">
                                                        {slots.filter(slot => slot.availability_reason === 'out_of_shift').map(slot => (
                                                            <button
                                                                key={slot.start_time}
                                                                className="time-slot"
                                                                disabled
                                                                style={{ opacity: 0.4, cursor: 'not-allowed', textDecoration: 'line-through' }}
                                                            >
                                                                {slot.start_time}
                                                            </button>
                                                        ))}
                                                    </div>
                                                </div>
                                            )}

                                            {slots.length === 0 && <p className="selection-helper">No slots configured for this date</p>}
                                        </div>
                                    )}
                                </>
                            )}
                        </section>
                    </div>

                    {/* Sidebar: Booking Summary */}
                    <div className="booking-sidebar">
                        <div className="glass summary-card">
                            <h3 style={{ marginBottom: 20 }}>Summary</h3>

                            <div className="summary-details">
                                <div className="summary-row">
                                    <span>Service</span>
                                    <strong>{selectedServiceData?.name || '-'}</strong>
                                </div>
                                <div className="summary-row">
                                    <span>Stylist</span>
                                    <strong>{staff.find(s => s.id === selectedStaff)?.name || '-'}</strong>
                                </div>
                                <div className="summary-row">
                                    <span>Date</span>
                                    <strong>{selectedDate || '-'}</strong>
                                </div>
                                <div className="summary-row">
                                    <span>Time (Booking)</span>
                                    <strong>{selectedTime || '-'}</strong>
                                </div>
                                {waitlistSlots.length > 0 && (
                                    <div className="summary-row" style={{ marginTop: 8 }}>
                                        <span>Waitlist Preferences</span>
                                        <div style={{ textAlign: 'right', display: 'flex', gap: 4, flexWrap: 'wrap', justifyContent: 'flex-end', maxWidth: 150 }}>
                                            {waitlistSlots.map(t => (
                                                <span key={t} className="badge badge-info" style={{ fontSize: '0.65rem', padding: '2px 6px' }}>{t}</span>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>

                            <hr style={{ margin: '20px 0', border: 'none', borderTop: '1px solid var(--border)' }} />

                            <div className="form-group">
                                <label style={{ fontSize: '0.85rem' }}>Special Requests</label>
                                <textarea
                                    className="form-control"
                                    placeholder="Add notes..."
                                    rows="2"
                                    value={notes}
                                    onChange={(e) => setNotes(e.target.value)}
                                    style={{ fontSize: '0.9rem' }}
                                ></textarea>
                            </div>

                            <div className="promo-inline">
                                <input type="text" className="form-control form-control-sm" placeholder="Promo code" value={promo} onChange={(e) => setPromo(e.target.value)} />
                                <button className="btn btn-secondary btn-sm" onClick={validatePromo}>Apply</button>
                            </div>

                            <div className="price-breakdown">
                                <div className="price-row">
                                    <span>Service Price</span>
                                    <span>${price.toFixed(2)}</span>
                                </div>
                                {discount > 0 && (
                                    <div className="price-row discount">
                                        <span>Discount ({discount}%)</span>
                                        <span>-${(price * discount / 100).toFixed(2)}</span>
                                    </div>
                                )}
                                <div className="price-row total">
                                    <span>Total</span>
                                    <span className="gradient-text">${finalPrice.toFixed(2)}</span>
                                </div>
                            </div>

                            <button
                                className={`btn btn-lg btn-block ${!selectedTime && waitlistSlots.length > 0 ? 'btn-secondary' : 'btn-primary'}`}
                                style={{ marginTop: 24, width: '100%' }}
                                onClick={handleBook}
                                disabled={(!selectedTime && waitlistSlots.length === 0) || loading}
                            >
                                {loading ? <div className="spinner" style={{ width: 20, height: 20, margin: 0 }} /> :
                                    (selectedTime && waitlistSlots.length > 0 ? '✨ Book + Waitlist' :
                                        !selectedTime && waitlistSlots.length > 0 ? '🤝 Join Waitlist' : '✨ Book Now')}
                            </button>

                            {!isReady && (
                                <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', textAlign: 'center', marginTop: 12 }}>
                                    Please complete all sections to book
                                </p>
                            )}
                        </div>
                    </div>
                </div>
            </div>

            <style>{`
        .booking-grid { display: grid; grid-template-columns: 1fr 340px; gap: 32px; align-items: start; }
        .booking-section { padding: 24px; border-radius: var(--radius-lg); margin-bottom: 24px; transition: var(--transition-normal); }
        .booking-section.disabled { opacity: 0.5; pointer-events: none; filter: grayscale(0.5); }
        .section-title { font-size: 1.1rem; margin-bottom: 20px; display: flex; align-items: center; gap: 10px; }
        .selection-helper { font-size: 0.9rem; color: var(--text-muted); text-align: center; padding: 20px; }
        
        .service-options { display: grid; gap: 10px; }
        .service-option { display: flex; align-items: center; padding: 14px 18px; border-radius: var(--radius-md); border: 1px solid var(--border); cursor: pointer; transition: all 0.2s ease; }
        .service-option:hover { border-color: var(--primary-light); background: rgba(255,255,255,0.02); }
        .service-option.selected { border-color: var(--primary); background: rgba(124,58,237,0.08); box-shadow: 0 0 0 1px var(--primary); }
        
        .staff-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 12px; }
        .staff-card { display: flex; align-items: center; gap: 12px; padding: 12px; border-radius: var(--radius-md); border: 1px solid var(--border); cursor: pointer; transition: all 0.2s ease; }
        .staff-card:hover { border-color: var(--primary-light); }
        .staff-card.selected { border-color: var(--primary); background: rgba(124,58,237,0.08); }
        .staff-avatar { width: 40px; height: 40px; border-radius: 50%; background: var(--bg-elevated); display: flex; align-items: center; justify-content: center; color: var(--text-secondary); }
        .staff-info { display: flex; flex-direction: column; }
        .staff-name { font-weight: 600; font-size: 0.95rem; }
        .staff-role { font-size: 0.75rem; color: var(--text-muted); }

        .time-slots { display: grid; grid-template-columns: repeat(auto-fill, minmax(80px, 1fr)); gap: 8px; }
        .slots-group { margin-bottom: 24px; }
        .group-label { font-size: 0.8rem; font-weight: 700; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
        .time-slot { padding: 8px; border-radius: var(--radius-sm); border: 1px solid var(--border); background: var(--bg-elevated); font-size: 0.85rem; font-weight: 600; cursor: pointer; transition: 0.2s; }
        .time-slot:hover:not(.busy) { border-color: var(--primary); }
        .time-slot.selected { background: var(--primary); color: white; border-color: var(--primary); }
        .time-slot.busy { border-style: dashed; opacity: 0.7; }
        .time-slot.waitlist-selected { background: var(--secondary); color: white; border-color: var(--secondary); opacity: 1; }

        .summary-card { padding: 24px; border-radius: var(--radius-xl); position: sticky; top: 100px; }
        .summary-row { display: flex; justify-content: space-between; font-size: 0.9rem; margin-bottom: 12px; }
        .summary-row span { color: var(--text-muted); }
        .promo-inline { display: flex; gap: 8px; margin-top: 20px; }
        .price-breakdown { margin-top: 24px; display: grid; gap: 10px; }
        .price-row { display: flex; justify-content: space-between; font-size: 0.95rem; }
        .price-row.discount { color: var(--success); }
        .price-row.total { font-size: 1.25rem; font-weight: 700; border-top: 1px solid var(--border); pt: 16px; margin-top: 6px; }

        @media (max-width: 992px) {
            .booking-grid { grid-template-columns: 1fr; }
            .booking-sidebar { order: -1; }
            .summary-card { position: static; }
        }
      `}</style>
        </div>
    );
}
