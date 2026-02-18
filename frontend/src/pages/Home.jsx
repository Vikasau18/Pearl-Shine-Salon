import { Link } from 'react-router-dom';
import { Scissors, Search, Star, Calendar, Shield, Sparkles, ArrowRight } from 'lucide-react';
import { useState, useEffect } from 'react';
import api from '../api/client';
import SalonCard from '../components/SalonCard';

export default function Home() {
    const [featured, setFeatured] = useState([]);
    const [search, setSearch] = useState('');

    useEffect(() => {
        api.get('/salons?limit=6').then(res => setFeatured(res.data.data || [])).catch(() => { });
    }, []);

    return (
        <div className="home-page">
            {/* Hero */}
            <section className="hero">
                <div className="container">
                    <div className="hero-content">
                        <div className="hero-badge">
                            <Sparkles size={16} /> Premium Salon Booking
                        </div>
                        <h1>
                            Discover & Book <span className="gradient-text">Beauty Services</span> Near You
                        </h1>
                        <p>Find the best salons, book appointments instantly, and look your best. All in one place.</p>

                        <div className="hero-search glass">
                            <Search size={20} />
                            <input
                                type="text"
                                placeholder="Search salons, services, or city..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                            <Link to={`/salons${search ? `?search=${search}` : ''}`} className="btn btn-primary">
                                Search
                            </Link>
                        </div>

                        <div className="hero-stats">
                            <div className="hero-stat">
                                <strong>500+</strong>
                                <span>Salons</span>
                            </div>
                            <div className="hero-stat">
                                <strong>10K+</strong>
                                <span>Bookings</span>
                            </div>
                            <div className="hero-stat">
                                <strong>4.8</strong>
                                <span>Avg Rating</span>
                            </div>
                        </div>
                    </div>
                    <div className="hero-visual">
                        <div className="hero-image-grid">
                            <img src="https://images.unsplash.com/photo-1560066984-138dadb4c035?w=400" alt="Salon" className="hero-img-1" />
                            <img src="https://images.unsplash.com/photo-1522337360788-8b13dee7a37e?w=300" alt="Styling" className="hero-img-2" />
                            <img src="https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=300" alt="Spa" className="hero-img-3" />
                        </div>
                    </div>
                </div>
            </section>

            {/* Features */}
            <section className="features">
                <div className="container">
                    <h2>Why Choose <span className="gradient-text">Saloon</span></h2>
                    <div className="grid-4 features-grid">
                        <div className="feature-card glass">
                            <div className="feature-icon" style={{ background: 'rgba(124,58,237,0.15)' }}>
                                <Search size={24} color="var(--primary)" />
                            </div>
                            <h3>Discover Salons</h3>
                            <p>Browse salons near you with ratings, reviews, and galleries</p>
                        </div>
                        <div className="feature-card glass">
                            <div className="feature-icon" style={{ background: 'rgba(6,182,212,0.15)' }}>
                                <Calendar size={24} color="var(--accent)" />
                            </div>
                            <h3>Easy Booking</h3>
                            <p>Book appointments in seconds with real-time availability</p>
                        </div>
                        <div className="feature-card glass">
                            <div className="feature-icon" style={{ background: 'rgba(16,185,129,0.15)' }}>
                                <Star size={24} color="var(--success)" />
                            </div>
                            <h3>Verified Reviews</h3>
                            <p>Read genuine reviews from verified customers</p>
                        </div>
                        <div className="feature-card glass">
                            <div className="feature-icon" style={{ background: 'rgba(244,63,94,0.15)' }}>
                                <Shield size={24} color="var(--secondary)" />
                            </div>
                            <h3>Secure Payments</h3>
                            <p>Safe and secure payment processing for all services</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Featured Salons */}
            {featured.length > 0 && (
                <section className="featured-section">
                    <div className="container">
                        <div className="section-header">
                            <h2>Featured <span className="gradient-text">Salons</span></h2>
                            <Link to="/salons" className="btn btn-outline btn-sm">
                                View All <ArrowRight size={16} />
                            </Link>
                        </div>
                        <div className="grid-3">
                            {featured.map(salon => (
                                <SalonCard key={salon.id} salon={salon} />
                            ))}
                        </div>
                    </div>
                </section>
            )}

            {/* CTA */}
            <section className="cta-section">
                <div className="container">
                    <div className="cta-card glass">
                        <div className="cta-content">
                            <Scissors size={32} />
                            <h2>Own a Salon?</h2>
                            <p>Join our platform and reach thousands of customers. Manage appointments, staff, and analytics — all from one dashboard.</p>
                            <Link to="/register" className="btn btn-primary btn-lg">
                                Get Started Free <ArrowRight size={20} />
                            </Link>
                        </div>
                    </div>
                </div>
            </section>

            <style>{`
        .hero {
          padding: 80px 0 60px;
          overflow: hidden;
        }
        .hero .container {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 60px;
          align-items: center;
        }
        .hero-badge {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          background: rgba(124,58,237,0.15);
          color: var(--primary-light);
          padding: 8px 16px;
          border-radius: 30px;
          font-size: 0.85rem;
          font-weight: 600;
          margin-bottom: 24px;
          border: 1px solid rgba(124,58,237,0.2);
        }
        .hero h1 {
          font-size: 3.2rem;
          font-weight: 800;
          line-height: 1.15;
          margin-bottom: 20px;
        }
        .hero p {
          font-size: 1.15rem;
          color: var(--text-secondary);
          margin-bottom: 32px;
          max-width: 500px;
        }
        .hero-search {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 8px 8px 8px 20px;
          border-radius: 60px;
          margin-bottom: 40px;
          max-width: 520px;
        }
        .hero-search input {
          flex: 1;
          background: none;
          border: none;
          outline: none;
          color: var(--text);
          font-size: 1rem;
          font-family: inherit;
        }
        .hero-search input::placeholder { color: var(--text-muted); }
        .hero-search .btn { border-radius: 50px; padding: 10px 24px; }
        .hero-stats {
          display: flex;
          gap: 40px;
        }
        .hero-stat strong {
          display: block;
          font-size: 1.8rem;
          font-weight: 800;
          background: linear-gradient(135deg, var(--primary-light), var(--accent));
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
        }
        .hero-stat span {
          font-size: 0.85rem;
          color: var(--text-muted);
        }
        .hero-visual {
          display: flex;
          justify-content: center;
        }
        .hero-image-grid {
          display: grid;
          grid-template-columns: 1fr 1fr;
          grid-template-rows: 1fr 1fr;
          gap: 12px;
          max-width: 450px;
        }
        .hero-image-grid img {
          border-radius: var(--radius-lg);
          object-fit: cover;
          height: 200px;
          width: 100%;
        }
        .hero-img-1 {
          grid-row: 1/3;
          height: 100% !important;
        }
        .features {
          padding: 80px 0;
        }
        .features h2 {
          font-size: 2.2rem;
          text-align: center;
          margin-bottom: 48px;
          font-weight: 800;
        }
        .feature-card {
          padding: 32px;
          border-radius: var(--radius-lg);
          text-align: center;
          transition: var(--transition);
        }
        .feature-card:hover {
          transform: translateY(-4px);
          box-shadow: var(--shadow-glow);
        }
        .feature-icon {
          width: 56px;
          height: 56px;
          border-radius: var(--radius-md);
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 20px;
        }
        .feature-card h3 {
          font-size: 1.1rem;
          margin-bottom: 8px;
        }
        .feature-card p {
          color: var(--text-secondary);
          font-size: 0.9rem;
        }
        .featured-section {
          padding: 60px 0;
        }
        .section-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 32px;
        }
        .section-header h2 {
          font-size: 1.8rem;
          font-weight: 800;
        }
        .cta-section {
          padding: 60px 0 80px;
        }
        .cta-card {
          border-radius: var(--radius-xl);
          padding: 60px;
          text-align: center;
          background: linear-gradient(135deg, rgba(124,58,237,0.1), rgba(244,63,94,0.05));
        }
        .cta-content h2 {
          font-size: 2rem;
          margin: 16px 0 12px;
        }
        .cta-content p {
          color: var(--text-secondary);
          max-width: 500px;
          margin: 0 auto 24px;
        }
        @media (max-width: 768px) {
          .hero .container { grid-template-columns: 1fr; gap: 40px; }
          .hero h1 { font-size: 2rem; }
          .hero-visual { display: none; }
          .hero-stats { gap: 24px; }
          .cta-card { padding: 40px 24px; }
        }
      `}</style>
        </div>
    );
}
