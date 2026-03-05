import { useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../api/client';
import { Mail, ArrowLeft, Scissors } from 'lucide-react';
import toast from 'react-hot-toast';

export default function ForgotPassword() {
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const [submitted, setSubmitted] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            await api.post('/auth/forgot-password', { email });
            setSubmitted(true);
            toast.success('Reset code sent to your email');
        } catch (err) {
            const errorMsg = err.response?.data?.error || 'Something went wrong. Please try again.';
            toast.error(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-card glass">
                <div className="auth-header">
                    <Scissors size={32} className="gradient-text" />
                    <h1>Forgot Password?</h1>
                    <p>Enter your email to receive a 6-digit reset code</p>
                </div>

                {!submitted ? (
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
                        <button type="submit" className="btn btn-primary btn-lg" style={{ width: '100%' }} disabled={loading}>
                            {loading ? <div className="spinner" style={{ width: 20, height: 20, marginBottom: 0 }} /> : <><Mail size={18} /> Send Reset Code</>}
                        </button>
                    </form>
                ) : (
                    <div style={{ textAlign: 'center', padding: '20px 0' }}>
                        <div style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e', padding: '16px', borderRadius: '12px', marginBottom: '24px' }}>
                            <p>We've sent a 6-digit reset code to <strong>{email}</strong> if it exists in our system.</p>
                        </div>
                        <Link to={`/reset-password?email=${email}`} className="btn btn-primary" style={{ width: '100%', marginBottom: 16 }}>
                            Enter Reset Code
                        </Link>
                        <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
                            Didn't receive the email? Check your spam folder or try again.
                        </p>
                    </div>
                )}

                <p style={{ textAlign: 'center', marginTop: 24 }}>
                    <Link to="/login" style={{ fontWeight: 600, display: 'inline-flex', alignItems: 'center', gap: 6 }}>
                        <ArrowLeft size={16} /> Back to Sign In
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
