import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Search, MapPin, SlidersHorizontal } from 'lucide-react';
import api from '../api/client';
import SalonCard from '../components/SalonCard';

export default function SalonList() {
    const [searchParams, setSearchParams] = useSearchParams();
    const [salons, setSalons] = useState([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState(searchParams.get('search') || '');
    const [city, setCity] = useState(searchParams.get('city') || '');
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    const fetchSalons = async () => {
        setLoading(true);
        try {
            const params = { page, limit: 12 };
            if (search) params.search = search;
            if (city) params.city = city;
            const res = await api.get('/salons', { params });
            setSalons(res.data.data || []);
            setTotalPages(res.data.total_pages || 1);
        } catch {
            setSalons([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchSalons(); }, [page]);

    const handleSearch = (e) => {
        e.preventDefault();
        setPage(1);
        fetchSalons();
    };

    return (
        <div className="page">
            <div className="container">
                <div className="page-header">
                    <h1>Explore <span className="gradient-text">Salons</span></h1>
                    <p>Discover the best beauty services near you</p>
                </div>

                <form onSubmit={handleSearch} className="search-bar glass">
                    <div className="search-field">
                        <Search size={18} />
                        <input type="text" placeholder="Search salons or services..." value={search} onChange={(e) => setSearch(e.target.value)} />
                    </div>
                    <div className="search-field">
                        <MapPin size={18} />
                        <input type="text" placeholder="City" value={city} onChange={(e) => setCity(e.target.value)} />
                    </div>
                    <button type="submit" className="btn btn-primary">
                        <SlidersHorizontal size={16} /> Filter
                    </button>
                </form>

                {loading ? (
                    <div className="loading"><div className="spinner" /><p>Finding salons...</p></div>
                ) : salons.length === 0 ? (
                    <div className="empty-state">
                        <h3>No salons found</h3>
                        <p>Try adjusting your search or filters</p>
                    </div>
                ) : (
                    <>
                        <div className="grid-3" style={{ marginTop: 32 }}>
                            {salons.map(salon => (
                                <SalonCard key={salon.id} salon={salon} />
                            ))}
                        </div>
                        {totalPages > 1 && (
                            <div className="pagination">
                                <button className="btn btn-secondary btn-sm" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>Previous</button>
                                <span>Page {page} of {totalPages}</span>
                                <button className="btn btn-secondary btn-sm" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>Next</button>
                            </div>
                        )}
                    </>
                )}
            </div>

            <style>{`
        .search-bar { display: flex; gap: 12px; padding: 12px; border-radius: 60px; flex-wrap: wrap; }
        .search-field { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 200px; padding: 0 12px; }
        .search-field input { background: none; border: none; outline: none; color: var(--text); font-size: 0.95rem; font-family: inherit; width: 100%; }
        .search-field input::placeholder { color: var(--text-muted); }
        .search-bar .btn { border-radius: 50px; white-space: nowrap; }
        .pagination { display: flex; align-items: center; justify-content: center; gap: 16px; margin-top: 40px; color: var(--text-secondary); }
      `}</style>
        </div>
    );
}
