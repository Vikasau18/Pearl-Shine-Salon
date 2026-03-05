import { useState, useEffect } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import api from '../api/client';
import { Lock, Eye, EyeOff, Scissors, CheckCircle } from 'lucide-react';
import toast from 'react-hot-toast';

export default function ResetPassword() {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const [email, setEmail] = useState(searchParams.get('email') || '');
    const [code, setCode] = useState(searchParams.get('token') || ''); // Allow 'token' param to fill code for legacy/convenience
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [showPw, setShowPw] = useState(false);
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            return toast.error('Passwords do not match');
        }
        if (password.length < 6) {
            return toast.error('Password must be at least 6 characters');
        }
        if (!code) {
            return toast.error('Please enter the reset code');
        }

        setLoading(true);
        try {
            await api.post('/auth/reset-password', {
                email,
                code,
                new_password: password
            });
            setSuccess(true);
            toast.success('Password reset successful!');
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to reset password. Code may be expired or incorrect.');
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div className="auth-page">
                <div className="auth-card glass" style={{ textAlign: 'center' }}>
                    <CheckCircle size={64} color="#22c55e" style={{ marginBottom: 24 }} />
                    <h1>Success!</h1>
                    <p style={{ color: 'var(--text-secondary)', marginBottom: 32 }}>
                        Your password has been successfully reset. You can now use your new password to sign in.
                    </p>
                    <Link to="/login" className="btn btn-primary btn-lg" style={{ width: '100%' }}>
                        Go to Sign In
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="auth-page">
            <div className="auth-card glass">
                <div className="auth-header">
                    <Scissors size={32} className="gradient-text" />
                    <h1>Reset Password</h1>
                    <p>Enter the 6-digit code sent to your email</p>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Email Address</label>
                        <input
                            type="email"
                            className="form-control"
                            placeholder="you@example.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Reset Code</label>
                        <input
                            type="text"
                            className="form-control"
                            placeholder="123456"
                            value={code}
                            onChange={(e) => setCode(e.target.value)}
                            required
                            maxLength={32} // Keep flexible if we ever change length
                        />
                    </div>

                    <div className="form-group">
                        <label>New Password</label>
                        <div style={{ position: 'relative' }}>
                            <input
                                type={showPw ? 'text' : 'password'}
                                className="form-control"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                            <button
                                type="button"
                                onClick={() => setShowPw(!showPw)}
                                style={{ position: 'absolute', right: 12, top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer' }}
                            >
                                {showPw ? <EyeOff size={18} /> : <Eye size={18} />}
                            </button>
                        </div>
                    </div>

                    <div className="form-group">
                        <label>Confirm New Password</label>
                        <input
                            type="password"
                            className="form-control"
                            placeholder="••••••••"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                        />
                    </div>

                    <button type="submit" className="btn btn-primary btn-lg" style={{ width: '100%' }} disabled={loading}>
                        {loading ? <div className="spinner" style={{ width: 20, height: 20, marginBottom: 0 }} /> : <><Lock size={18} /> Reset Password</>}
                    </button>
                </form>

                <p style={{ textAlign: 'center', marginTop: 24 }}>
                    <Link to="/forgot-password" style={{ fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
                        Request a new code
                    </Link>
                </p>
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
      `}</style>
        </div>
    );
}
