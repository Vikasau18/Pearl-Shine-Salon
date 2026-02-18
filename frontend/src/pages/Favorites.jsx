import { useState, useEffect } from 'react';
import { Heart } from 'lucide-react';
import api from '../api/client';
import SalonCard from '../components/SalonCard';

export default function Favorites() {
    const [favorites, setFavorites] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetch = async () => {
        try {
            const res = await api.get('/favorites');
            setFavorites(res.data || []);
        } catch { setFavorites([]); }
        finally { setLoading(false); }
    };

    useEffect(() => { fetch(); }, []);

    return (
        <div className="page">
            <div className="container">
                <div className="page-header">
                    <h1>My <span className="gradient-text">Favorites</span></h1>
                    <p>Salons you've saved for quick access</p>
                </div>

                {loading ? (
                    <div className="loading"><div className="spinner" /><p>Loading favorites...</p></div>
                ) : favorites.length === 0 ? (
                    <div className="empty-state">
                        <Heart size={48} style={{ marginBottom: 16, opacity: 0.3 }} />
                        <h3>No favorites yet</h3>
                        <p>Heart a salon to save it here</p>
                    </div>
                ) : (
                    <div className="grid-3">
                        {favorites.map(salon => (
                            <SalonCard key={salon.id} salon={{ ...salon, is_favorited: true }} onFavToggle={fetch} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
