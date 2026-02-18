import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { LogIn, Eye, EyeOff, Scissors } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Login() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPw, setShowPw] = useState(false);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const data = await login(email, password);
            toast.success(`Welcome back, ${data.user.name}!`);
            if (data.user.role === 'salon_owner') navigate('/dashboard');
            else navigate('/');
        } catch (err) {
            toast.error(err.response?.data?.error || 'Login failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-card glass">
                <div className="auth-header">
                    <Scissors size={32} className="gradient-text" />
                    <h1>Welcome Back</h1>
                    <p>Sign in to continue to Saloon</p>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Email Address</label>
                        <input type="email" className="form-control" placeholder="you@example.com" value={email} onChange={(e) => setEmail(e.target.value)} required />
                    </div>
                    <div className="form-group">
                        <label>Password</label>
                        <div style={{ position: 'relative' }}>
                            <input type={showPw ? 'text' : 'password'} className="form-control" placeholder="••••••••" value={password} onChange={(e) => setPassword(e.target.value)} required />
                            <button type="button" onClick={() => setShowPw(!showPw)} style={{ position: 'absolute', right: 12, top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}>
                                {showPw ? <EyeOff size={18} /> : <Eye size={18} />}
                            </button>
                        </div>
                    </div>
                    <button type="submit" className="btn btn-primary btn-lg" style={{ width: '100%' }} disabled={loading}>
                        {loading ? <div className="spinner" style={{ width: 20, height: 20, marginBottom: 0 }} /> : <><LogIn size={18} /> Sign In</>}
                    </button>
                </form>

                <div className="auth-divider"><span>OR</span></div>

                <button className="btn btn-secondary" style={{ width: '100%' }} onClick={() => toast('Google OAuth requires configuration')}>
                    <img src="https://www.google.com/favicon.ico" alt="Google" width={16} /> Continue with Google
                </button>

                <p style={{ textAlign: 'center', marginTop: 24, color: 'var(--text-secondary)' }}>
                    Don't have an account? <Link to="/register" style={{ fontWeight: 600 }}>Sign Up</Link>
                </p>

                <div className="auth-demo">
                    <p><strong>Demo Accounts:</strong></p>
                    <button onClick={() => { setEmail('customer@test.com'); setPassword('password123'); }} className="btn btn-sm btn-outline" style={{ marginRight: 8 }}>Customer</button>
                    <button onClick={() => { setEmail('owner1@test.com'); setPassword('password123'); }} className="btn btn-sm btn-outline">Salon Owner</button>
                </div>
            </div>

            <style>{`
        .auth-page {
          min-height: calc(100vh - 80px);
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 40px 20px;
        }
        .auth-card {
          max-width: 440px;
          width: 100%;
          padding: 40px;
          border-radius: var(--radius-xl);
        }
        .auth-header {
          text-align: center;
          margin-bottom: 32px;
        }
        .auth-header h1 {
          font-size: 1.8rem;
          margin: 12px 0 4px;
        }
        .auth-header p {
          color: var(--text-secondary);
        }
        .auth-divider {
          text-align: center;
          margin: 24px 0;
          position: relative;
        }
        .auth-divider::before {
          content: '';
          position: absolute;
          left: 0;
          right: 0;
          top: 50%;
          height: 1px;
          background: var(--border);
        }
        .auth-divider span {
          background: var(--bg-glass);
          padding: 0 16px;
          color: var(--text-muted);
          font-size: 0.85rem;
          position: relative;
        }
        .auth-demo {
          margin-top: 24px;
          padding: 16px;
          background: var(--bg-elevated);
          border-radius: var(--radius-md);
          text-align: center;
        }
        .auth-demo p {
          font-size: 0.8rem;
          color: var(--text-muted);
          margin-bottom: 8px;
        }
      `}</style>
        </div>
    );
}
