import { Link } from 'react-router-dom';
import { Star, MapPin, Heart } from 'lucide-react';
import api from '../api/client';
import { useState } from 'react';

export default function SalonCard({ salon, onFavToggle }) {
    const [fav, setFav] = useState(salon.is_favorited || false);

    const toggleFav = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        try {
            const res = await api.post(`/favorites/${salon.id}`);
            setFav(res.data.favorited);
            onFavToggle?.();
        } catch { }
    };

    return (
        <Link to={`/salons/${salon.id}`} className="salon-card card">
            <div className="salon-card-img-wrap">
                <img
                    src={salon.image_url || 'https://images.unsplash.com/photo-1560066984-138dadb4c035?w=400'}
                    alt={salon.name}
                    className="card-img"
                />
                <button className={`fav-btn ${fav ? 'active' : ''}`} onClick={toggleFav}>
                    <Heart size={18} fill={fav ? '#f43f5e' : 'none'} />
                </button>
                {salon.distance != null && (
                    <span className="distance-badge">{salon.distance.toFixed(1)} km</span>
                )}
            </div>
            <div className="card-body">
                <h3 className="salon-card-name">{salon.name}</h3>
                <div className="salon-card-location">
                    <MapPin size={14} />
                    <span>{salon.city}{salon.state ? `, ${salon.state}` : ''}</span>
                </div>
                <div className="salon-card-footer">
                    <div className="salon-card-rating">
                        <Star size={14} fill="#facc15" color="#facc15" />
                        <span>{salon.rating?.toFixed(1) || '0.0'}</span>
                        <span className="review-count">({salon.total_reviews || 0})</span>
                    </div>
                    <span className="badge badge-primary">Open</span>
                </div>
            </div>

            <style>{`
        .salon-card { display: block; color: inherit; text-decoration: none; }
        .salon-card-img-wrap { position: relative; overflow: hidden; }
        .salon-card-img-wrap .card-img { transition: transform 0.4s ease; }
        .salon-card:hover .card-img { transform: scale(1.05); }
        .fav-btn {
          position: absolute; top: 12px; right: 12px;
          width: 36px; height: 36px; border-radius: 50%;
          background: rgba(0,0,0,0.5); backdrop-filter: blur(8px);
          border: none; cursor: pointer; color: white;
          display: flex; align-items: center; justify-content: center;
          transition: var(--transition-fast);
        }
        .fav-btn.active { color: #f43f5e; }
        .fav-btn:hover { transform: scale(1.1); }
        .distance-badge {
          position: absolute; bottom: 12px; left: 12px;
          background: rgba(0,0,0,0.6); backdrop-filter: blur(8px);
          color: white; padding: 4px 10px; border-radius: 20px;
          font-size: 0.8rem; font-weight: 600;
        }
        .salon-card-name { font-size: 1.1rem; font-weight: 700; margin-bottom: 6px; }
        .salon-card-location {
          display: flex; align-items: center; gap: 4px;
          color: var(--text-secondary); font-size: 0.85rem; margin-bottom: 12px;
        }
        .salon-card-footer {
          display: flex; align-items: center; justify-content: space-between;
        }
        .salon-card-rating {
          display: flex; align-items: center; gap: 4px;
          font-weight: 600; font-size: 0.9rem;
        }
        .review-count { color: var(--text-muted); font-weight: 400; font-size: 0.8rem; }
      `}</style>
        </Link>
    );
}
