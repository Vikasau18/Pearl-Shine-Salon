import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { MapPin, Clock, Phone, Mail, Star, Calendar, Heart, Users } from 'lucide-react';
import api from '../api/client';
import StarRating from '../components/StarRating';

export default function SalonDetail() {
    const { id } = useParams();
    const [salon, setSalon] = useState(null);
    const [services, setServices] = useState([]);
    const [staff, setStaff] = useState([]);
    const [reviews, setReviews] = useState([]);
    const [gallery, setGallery] = useState([]);
    const [tab, setTab] = useState('services');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        Promise.all([
            api.get(`/salons/${id}`),
            api.get(`/salons/${id}/services`),
            api.get(`/salons/${id}/staff`),
            api.get(`/salons/${id}/reviews`),
            api.get(`/salons/${id}/gallery`),
        ]).then(([sRes, svRes, stRes, rRes, gRes]) => {
            setSalon(sRes.data);
            setServices(svRes.data);
            setStaff(stRes.data);
            setReviews(rRes.data);
            setGallery(gRes.data);
        }).catch(() => { }).finally(() => setLoading(false));
    }, [id]);

    if (loading) return <div className="loading"><div className="spinner" /><p>Loading salon...</p></div>;
    if (!salon) return <div className="empty-state"><h3>Salon not found</h3></div>;

    const categories = [...new Set(services.map(s => s.category).filter(Boolean))];

    return (
        <div className="page salon-detail">
            <div className="container">
                {/* Hero Banner */}
                <div className="salon-hero">
                    <img src={salon.image_url || 'https://images.unsplash.com/photo-1560066984-138dadb4c035?w=1200'} alt={salon.name} />
                    <div className="salon-hero-overlay">
                        <div className="salon-hero-info">
                            <h1>{salon.name}</h1>
                            <div className="salon-meta">
                                <span><MapPin size={16} /> {salon.address}, {salon.city}</span>
                                <span><Clock size={16} /> {salon.opening_time} – {salon.closing_time}</span>
                                <span><Star size={16} fill="#facc15" color="#facc15" /> {salon.rating?.toFixed(1)} ({salon.total_reviews} reviews)</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Action Row */}
                <div className="salon-actions">
                    <Link to={`/book/${id}`} className="btn btn-primary btn-lg">
                        <Calendar size={18} /> Book Appointment
                    </Link>
                    <button className="btn btn-secondary">
                        <Phone size={16} /> {salon.phone || 'Call'}
                    </button>
                    <button className="btn btn-secondary">
                        <Heart size={16} /> Favorite
                    </button>
                </div>

                {/* Description */}
                {salon.description && <p className="salon-desc">{salon.description}</p>}

                {/* Tabs */}
                <div className="tabs">
                    {['services', 'staff', 'reviews', 'gallery'].map(t => (
                        <button key={t} className={`tab ${tab === t ? 'active' : ''}`} onClick={() => setTab(t)}>
                            {t.charAt(0).toUpperCase() + t.slice(1)}
                            {t === 'reviews' && ` (${reviews.length})`}
                        </button>
                    ))}
                </div>

                {/* Services Tab */}
                {tab === 'services' && (
                    <div className="tab-content">
                        {categories.length > 0 ? categories.map(cat => (
                            <div key={cat} className="service-category">
                                <h3>{cat}</h3>
                                <div className="service-list">
                                    {services.filter(s => s.category === cat).map(svc => (
                                        <div key={svc.id} className="service-item glass">
                                            <div>
                                                <h4>{svc.name}</h4>
                                                {svc.description && <p>{svc.description}</p>}
                                                <span className="service-duration">{svc.duration_minutes} min</span>
                                            </div>
                                            <div className="service-right">
                                                <span className="service-price">${svc.price.toFixed(2)}</span>
                                                <Link to={`/book/${id}?service=${svc.id}`} className="btn btn-primary btn-sm">Book</Link>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )) : (
                            <div className="service-list">
                                {services.map(svc => (
                                    <div key={svc.id} className="service-item glass">
                                        <div>
                                            <h4>{svc.name}</h4>
                                            {svc.description && <p>{svc.description}</p>}
                                            <span className="service-duration">{svc.duration_minutes} min</span>
                                        </div>
                                        <div className="service-right">
                                            <span className="service-price">${svc.price.toFixed(2)}</span>
                                            <Link to={`/book/${id}?service=${svc.id}`} className="btn btn-primary btn-sm">Book</Link>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

                {/* Staff Tab */}
                {tab === 'staff' && (
                    <div className="tab-content grid-3">
                        {staff.map(s => (
                            <div key={s.id} className="staff-card card">
                                <div className="staff-avatar">
                                    {s.avatar_url ? <img src={s.avatar_url} alt={s.name} /> : <Users size={32} />}
                                </div>
                                <div className="card-body" style={{ textAlign: 'center' }}>
                                    <h4>{s.name}</h4>
                                    <span className="badge badge-primary">{s.role}</span>
                                    {s.specialization && <p style={{ marginTop: 8, fontSize: '0.85rem', color: 'var(--text-secondary)' }}>{s.specialization}</p>}
                                    {s.services?.length > 0 && (
                                        <div style={{ marginTop: 8, display: 'flex', flexWrap: 'wrap', gap: 4, justifyContent: 'center' }}>
                                            {s.services.slice(0, 3).map(sv => (
                                                <span key={sv.id} className="badge badge-info">{sv.name}</span>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}

                {/* Reviews Tab */}
                {tab === 'reviews' && (
                    <div className="tab-content">
                        {reviews.length === 0 ? (
                            <div className="empty-state"><h3>No reviews yet</h3><p>Be the first to leave a review!</p></div>
                        ) : reviews.map(r => (
                            <div key={r.id} className="review-item glass">
                                <div className="review-header">
                                    <div className="review-user">
                                        {r.customer_avatar ? <img src={r.customer_avatar} alt="" className="review-avatar" /> : <div className="review-avatar-placeholder">{r.customer_name?.[0]}</div>}
                                        <div>
                                            <strong>{r.customer_name}</strong>
                                            <StarRating rating={r.rating} size={14} />
                                        </div>
                                    </div>
                                    <span className="review-date">{new Date(r.created_at).toLocaleDateString()}</span>
                                </div>
                                {r.comment && <p className="review-comment">{r.comment}</p>}
                            </div>
                        ))}
                    </div>
                )}

                {/* Gallery Tab */}
                {tab === 'gallery' && (
                    <div className="tab-content grid-3">
                        {gallery.length === 0 ? (
                            <div className="empty-state" style={{ gridColumn: '1/-1' }}><h3>No gallery images</h3></div>
                        ) : gallery.map(g => (
                            <div key={g.id} className="gallery-item">
                                <img src={g.image_url} alt={g.caption} />
                                {g.caption && <span>{g.caption}</span>}
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <style>{`
        .salon-hero { position: relative; height: 360px; border-radius: var(--radius-xl); overflow: hidden; margin-bottom: 24px; }
        .salon-hero img { width: 100%; height: 100%; object-fit: cover; }
        .salon-hero-overlay { position: absolute; inset: 0; background: linear-gradient(transparent 40%, rgba(0,0,0,0.8)); display: flex; align-items: flex-end; padding: 32px; }
        .salon-hero-info h1 { font-size: 2.2rem; font-weight: 800; margin-bottom: 8px; }
        .salon-meta { display: flex; flex-wrap: wrap; gap: 20px; color: var(--text-secondary); font-size: 0.9rem; }
        .salon-meta span { display: flex; align-items: center; gap: 6px; }
        .salon-actions { display: flex; gap: 12px; margin-bottom: 24px; flex-wrap: wrap; }
        .salon-desc { color: var(--text-secondary); line-height: 1.7; margin-bottom: 32px; max-width: 700px; }
        .tabs { display: flex; gap: 4px; margin-bottom: 32px; background: var(--bg-secondary); padding: 4px; border-radius: var(--radius-md); }
        .tab { padding: 12px 24px; background: none; border: none; color: var(--text-secondary); cursor: pointer; font-family: inherit; font-weight: 500; border-radius: var(--radius-sm); transition: var(--transition-fast); }
        .tab.active { background: var(--primary); color: white; }
        .tab:hover:not(.active) { color: var(--text); }
        .tab-content { min-height: 200px; }
        .service-category { margin-bottom: 32px; }
        .service-category h3 { font-size: 1.2rem; margin-bottom: 16px; padding-bottom: 8px; border-bottom: 1px solid var(--border); }
        .service-list { display: flex; flex-direction: column; gap: 12px; }
        .service-item { display: flex; justify-content: space-between; align-items: center; padding: 20px; border-radius: var(--radius-md); }
        .service-item h4 { font-size: 1rem; margin-bottom: 4px; }
        .service-item p { color: var(--text-secondary); font-size: 0.85rem; margin-bottom: 4px; }
        .service-duration { font-size: 0.8rem; color: var(--text-muted); }
        .service-right { display: flex; align-items: center; gap: 16px; }
        .service-price { font-size: 1.2rem; font-weight: 700; color: var(--primary-light); }
        .staff-card { text-align: center; }
        .staff-avatar { width: 80px; height: 80px; border-radius: 50%; margin: 20px auto 0; overflow: hidden; background: var(--bg-elevated); display: flex; align-items: center; justify-content: center; }
        .staff-avatar img { width: 100%; height: 100%; object-fit: cover; }
        .review-item { padding: 20px; border-radius: var(--radius-md); margin-bottom: 12px; }
        .review-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
        .review-user { display: flex; align-items: center; gap: 12px; }
        .review-avatar { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; }
        .review-avatar-placeholder { width: 40px; height: 40px; border-radius: 50%; background: var(--primary); display: flex; align-items: center; justify-content: center; font-weight: 700; color: white; }
        .review-date { color: var(--text-muted); font-size: 0.8rem; }
        .review-comment { color: var(--text-secondary); line-height: 1.6; }
        .gallery-item { border-radius: var(--radius-md); overflow: hidden; position: relative; }
        .gallery-item img { width: 100%; height: 250px; object-fit: cover; }
        .gallery-item span { position: absolute; bottom: 0; left: 0; right: 0; padding: 8px 12px; background: rgba(0,0,0,0.6); font-size: 0.85rem; }
        @media (max-width: 768px) {
          .salon-hero { height: 240px; }
          .salon-hero-info h1 { font-size: 1.5rem; }
          .tabs { overflow-x: auto; }
        }
      `}</style>
        </div>
    );
}
