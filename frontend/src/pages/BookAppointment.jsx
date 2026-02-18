import { useState, useEffect } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import { Calendar, Clock, User, CheckCircle } from 'lucide-react';
import api from '../api/client';
import toast from 'react-hot-toast';

export default function BookAppointment() {
    const { salonId } = useParams();
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const [step, setStep] = useState(1);
    const [services, setServices] = useState([]);
    const [staff, setStaff] = useState([]);
    const [selectedService, setSelectedService] = useState(searchParams.get('service') || '');
    const [selectedStaff, setSelectedStaff] = useState('');
    const [selectedDate, setSelectedDate] = useState('');
    const [selectedTime, setSelectedTime] = useState('');
    const [slots, setSlots] = useState([]);
    const [promo, setPromo] = useState('');
    const [discount, setDiscount] = useState(0);
    const [loading, setLoading] = useState(false);
    const [salon, setSalon] = useState(null);
    const [booked, setBooked] = useState(false);

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
            }).then(res => setSlots(res.data || []))
                .catch(() => setSlots([]));
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
            await api.post('/appointments', {
                salon_id: salonId,
                service_id: selectedService,
                staff_id: selectedStaff,
                date: selectedDate,
                start_time: selectedTime,
                promo_code: promo || undefined,
            });
            setBooked(true);
            toast.success('Appointment booked successfully!');
        } catch (err) {
            toast.error(err.response?.data?.error || 'Booking failed');
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
                    <p style={{ color: 'var(--text-secondary)', marginBottom: 32 }}>Your appointment has been scheduled. You'll receive a confirmation email shortly.</p>
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

    return (
        <div className="page">
            <div className="container">
                <div className="page-header">
                    <h1>Book at <span className="gradient-text">{salon?.name || 'Salon'}</span></h1>
                </div>

                {/* Steps indicator */}
                <div className="steps-bar">
                    {['Service', 'Staff', 'Date & Time', 'Confirm'].map((s, i) => (
                        <div key={i} className={`step-item ${step === i + 1 ? 'active' : step > i + 1 ? 'done' : ''}`}>
                            <div className="step-num">{step > i + 1 ? '✓' : i + 1}</div>
                            <span>{s}</span>
                        </div>
                    ))}
                </div>

                <div className="booking-content glass" style={{ padding: 32, borderRadius: 'var(--radius-xl)', marginTop: 32 }}>
                    {/* Step 1: Service */}
                    {step === 1 && (
                        <div>
                            <h2 style={{ marginBottom: 20 }}>Select a Service</h2>
                            <div className="service-options">
                                {services.map(svc => (
                                    <div key={svc.id} className={`service-option ${selectedService === svc.id ? 'selected' : ''}`} onClick={() => setSelectedService(svc.id)}>
                                        <div>
                                            <h4>{svc.name}</h4>
                                            <span>{svc.duration_minutes} min</span>
                                        </div>
                                        <strong>${svc.price.toFixed(2)}</strong>
                                    </div>
                                ))}
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 24 }}>
                                <button className="btn btn-primary" disabled={!selectedService} onClick={() => setStep(2)}>Next</button>
                            </div>
                        </div>
                    )}

                    {/* Step 2: Staff */}
                    {step === 2 && (
                        <div>
                            <h2 style={{ marginBottom: 20 }}>Choose a Stylist</h2>
                            <div className="grid-3">
                                {staff.map(s => (
                                    <div key={s.id} className={`staff-option card ${selectedStaff === s.id ? 'selected' : ''}`} onClick={() => setSelectedStaff(s.id)}>
                                        <div className="card-body" style={{ textAlign: 'center', padding: 24 }}>
                                            <div style={{ width: 64, height: 64, borderRadius: '50%', background: 'var(--bg-elevated)', margin: '0 auto 12px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                                <User size={28} />
                                            </div>
                                            <h4>{s.name}</h4>
                                            {s.specialization && <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>{s.specialization}</p>}
                                        </div>
                                    </div>
                                ))}
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 24 }}>
                                <button className="btn btn-secondary" onClick={() => setStep(1)}>Back</button>
                                <button className="btn btn-primary" disabled={!selectedStaff} onClick={() => setStep(3)}>Next</button>
                            </div>
                        </div>
                    )}

                    {/* Step 3: Date & Time */}
                    {step === 3 && (
                        <div>
                            <h2 style={{ marginBottom: 20 }}>Pick Date & Time</h2>
                            <div className="form-group">
                                <label><Calendar size={16} style={{ marginRight: 6 }} /> Date</label>
                                <input type="date" className="form-control" value={selectedDate} onChange={(e) => setSelectedDate(e.target.value)} min={new Date().toISOString().split('T')[0]} />
                            </div>
                            {slots.length > 0 && (
                                <div className="form-group">
                                    <label><Clock size={16} style={{ marginRight: 6 }} /> Available Slots</label>
                                    <div className="time-slots">
                                        {slots.map(slot => (
                                            <button
                                                key={slot.start_time}
                                                className={`time-slot ${selectedTime === slot.start_time ? 'selected' : ''} ${!slot.available ? 'disabled' : ''}`}
                                                onClick={() => slot.available && setSelectedTime(slot.start_time)}
                                                disabled={!slot.available}
                                            >
                                                {slot.start_time}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            )}
                            {selectedDate && slots.length === 0 && (
                                <p style={{ color: 'var(--text-muted)', textAlign: 'center', padding: 20 }}>No available slots for this date. Try another date.</p>
                            )}
                            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 24 }}>
                                <button className="btn btn-secondary" onClick={() => setStep(2)}>Back</button>
                                <button className="btn btn-primary" disabled={!selectedTime} onClick={() => setStep(4)}>Next</button>
                            </div>
                        </div>
                    )}

                    {/* Step 4: Confirm */}
                    {step === 4 && (
                        <div>
                            <h2 style={{ marginBottom: 20 }}>Confirm Booking</h2>
                            <div style={{ display: 'grid', gap: 16 }}>
                                <div className="confirm-row"><span>Service</span> <strong>{selectedServiceData?.name}</strong></div>
                                <div className="confirm-row"><span>Stylist</span> <strong>{staff.find(s => s.id === selectedStaff)?.name}</strong></div>
                                <div className="confirm-row"><span>Date</span> <strong>{selectedDate}</strong></div>
                                <div className="confirm-row"><span>Time</span> <strong>{selectedTime}</strong></div>
                                <div className="confirm-row"><span>Duration</span> <strong>{selectedServiceData?.duration_minutes} min</strong></div>
                                <hr style={{ border: 'none', borderTop: '1px solid var(--border)' }} />

                                <div style={{ display: 'flex', gap: 8 }}>
                                    <input type="text" className="form-control" placeholder="Promo code" value={promo} onChange={(e) => setPromo(e.target.value)} style={{ flex: 1 }} />
                                    <button className="btn btn-secondary" onClick={validatePromo}>Apply</button>
                                </div>

                                <div className="confirm-row"><span>Price</span> <strong>${price.toFixed(2)}</strong></div>
                                {discount > 0 && <div className="confirm-row"><span>Discount ({discount}%)</span> <strong style={{ color: 'var(--success)' }}>-${(price * discount / 100).toFixed(2)}</strong></div>}
                                <div className="confirm-row total"><span>Total</span> <strong className="gradient-text">${finalPrice.toFixed(2)}</strong></div>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 32 }}>
                                <button className="btn btn-secondary" onClick={() => setStep(3)}>Back</button>
                                <button className="btn btn-primary btn-lg" onClick={handleBook} disabled={loading}>
                                    {loading ? <div className="spinner" style={{ width: 20, height: 20, marginBottom: 0 }} /> : '✨ Confirm Booking'}
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            <style>{`
        .steps-bar { display: flex; justify-content: center; gap: 40px; }
        .step-item { display: flex; align-items: center; gap: 8px; color: var(--text-muted); }
        .step-item.active { color: var(--primary-light); }
        .step-item.done { color: var(--success); }
        .step-num { width: 32px; height: 32px; border-radius: 50%; background: var(--bg-elevated); display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 0.85rem; border: 2px solid var(--border); }
        .step-item.active .step-num { border-color: var(--primary); background: var(--primary); color: white; }
        .step-item.done .step-num { border-color: var(--success); background: var(--success); color: white; }
        .service-options { display: flex; flex-direction: column; gap: 8px; }
        .service-option { display: flex; justify-content: space-between; align-items: center; padding: 16px 20px; border-radius: var(--radius-md); border: 1px solid var(--border); cursor: pointer; transition: var(--transition-fast); }
        .service-option:hover { border-color: var(--primary-light); }
        .service-option.selected { border-color: var(--primary); background: rgba(124,58,237,0.1); }
        .service-option h4 { margin-bottom: 2px; }
        .service-option span { font-size: 0.8rem; color: var(--text-muted); }
        .service-option strong { font-size: 1.1rem; color: var(--primary-light); }
        .staff-option { cursor: pointer; }
        .staff-option.selected { border-color: var(--primary); box-shadow: var(--shadow-glow); }
        .time-slots { display: flex; flex-wrap: wrap; gap: 8px; }
        .time-slot { padding: 10px 20px; border: 1px solid var(--border); border-radius: var(--radius-md); background: var(--bg-elevated); color: var(--text); cursor: pointer; font-family: inherit; font-weight: 500; transition: var(--transition-fast); }
        .time-slot:hover { border-color: var(--primary-light); }
        .time-slot.selected { background: var(--primary); color: white; border-color: var(--primary); }
        .confirm-row { display: flex; justify-content: space-between; padding: 8px 0; }
        .confirm-row span { color: var(--text-secondary); }
        .confirm-row.total { font-size: 1.2rem; padding-top: 16px; }
        @media (max-width: 768px) { .steps-bar { gap: 16px; font-size: 0.85rem; } .steps-bar span { display: none; } }
      `}</style>
        </div>
    );
}
