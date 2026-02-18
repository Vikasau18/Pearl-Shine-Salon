import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { UserPlus, Scissors } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Register() {
    const { register } = useAuth();
    const navigate = useNavigate();
    const [form, setForm] = useState({ name: '', email: '', password: '', phone: '', role: 'customer' });
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            await register(form);
            toast.success('Account created! Welcome!');
            navigate(form.role === 'salon_owner' ? '/dashboard' : '/');
        } catch (err) {
            toast.error(err.response?.data?.error || 'Registration failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-card glass">
                <div className="auth-header">
                    <Scissors size={32} />
                    <h1>Create Account</h1>
                    <p>Join Saloon and start booking today</p>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Full Name</label>
                        <input type="text" className="form-control" placeholder="John Doe" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
                    </div>
                    <div className="form-group">
                        <label>Email</label>
                        <input type="email" className="form-control" placeholder="you@example.com" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} required />
                    </div>
                    <div className="form-group">
                        <label>Password</label>
                        <input type="password" className="form-control" placeholder="Min 6 characters" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} required minLength={6} />
                    </div>
                    <div className="form-group">
                        <label>Phone (optional)</label>
                        <input type="tel" className="form-control" placeholder="+1234567890" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
                    </div>
                    <div className="form-group">
                        <label>I am a</label>
                        <div style={{ display: 'flex', gap: 12 }}>
                            <button type="button" className={`btn ${form.role === 'customer' ? 'btn-primary' : 'btn-secondary'}`} style={{ flex: 1 }} onClick={() => setForm({ ...form, role: 'customer' })}>
                                Customer
                            </button>
                            <button type="button" className={`btn ${form.role === 'salon_owner' ? 'btn-primary' : 'btn-secondary'}`} style={{ flex: 1 }} onClick={() => setForm({ ...form, role: 'salon_owner' })}>
                                Salon Owner
                            </button>
                        </div>
                    </div>
                    <button type="submit" className="btn btn-primary btn-lg" style={{ width: '100%' }} disabled={loading}>
                        {loading ? <div className="spinner" style={{ width: 20, height: 20, marginBottom: 0 }} /> : <><UserPlus size={18} /> Create Account</>}
                    </button>
                </form>

                <p style={{ textAlign: 'center', marginTop: 24, color: 'var(--text-secondary)' }}>
                    Already have an account? <Link to="/login" style={{ fontWeight: 600 }}>Sign In</Link>
                </p>
            </div>

            <style>{`
        .auth-page { min-height: calc(100vh - 80px); display: flex; align-items: center; justify-content: center; padding: 40px 20px; }
        .auth-card { max-width: 440px; width: 100%; padding: 40px; border-radius: var(--radius-xl); }
        .auth-header { text-align: center; margin-bottom: 32px; }
        .auth-header h1 { font-size: 1.8rem; margin: 12px 0 4px; }
        .auth-header p { color: var(--text-secondary); }
      `}</style>
        </div>
    );
}
