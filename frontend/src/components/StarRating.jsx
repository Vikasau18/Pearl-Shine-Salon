import { Star } from 'lucide-react';

export default function StarRating({ rating, onRate, size = 20, interactive = false }) {
    return (
        <div className="stars" style={{ cursor: interactive ? 'pointer' : 'default' }}>
            {[1, 2, 3, 4, 5].map((star) => (
                <Star
                    key={star}
                    size={size}
                    fill={star <= rating ? '#facc15' : 'none'}
                    color={star <= rating ? '#facc15' : '#64748b'}
                    onClick={() => interactive && onRate?.(star)}
                    style={interactive ? { transition: 'transform 0.15s', cursor: 'pointer' } : {}}
                    onMouseEnter={(e) => interactive && (e.target.style.transform = 'scale(1.2)')}
                    onMouseLeave={(e) => interactive && (e.target.style.transform = 'scale(1)')}
                />
            ))}
        </div>
    );
}
