import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';

// Pages
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import SalonList from './pages/SalonList';
import SalonDetail from './pages/SalonDetail';
import BookAppointment from './pages/BookAppointment';
import MyAppointments from './pages/MyAppointments';
import Favorites from './pages/Favorites';
import Notifications from './pages/Notifications';
import Dashboard from './pages/Dashboard';

import './index.css';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Navbar />
        <main style={{ flex: 1 }}>
          <Routes>
            {/* Public */}
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/salons" element={<SalonList />} />
            <Route path="/salons/:id" element={<SalonDetail />} />

            {/* Authenticated */}
            <Route path="/book/:salonId" element={
              <ProtectedRoute><BookAppointment /></ProtectedRoute>
            } />
            <Route path="/appointments" element={
              <ProtectedRoute><MyAppointments /></ProtectedRoute>
            } />
            <Route path="/favorites" element={
              <ProtectedRoute><Favorites /></ProtectedRoute>
            } />
            <Route path="/notifications" element={
              <ProtectedRoute><Notifications /></ProtectedRoute>
            } />

            {/* Salon Owner */}
            <Route path="/dashboard" element={
              <ProtectedRoute roles={['salon_owner', 'admin']}><Dashboard /></ProtectedRoute>
            } />

            {/* 404 */}
            <Route path="*" element={
              <div className="page" style={{ textAlign: 'center', padding: '80px 20px' }}>
                <h1 style={{ fontSize: '6rem', fontWeight: 800, opacity: 0.1 }}>404</h1>
                <h2>Page Not Found</h2>
                <p style={{ color: 'var(--text-secondary)', marginTop: 8 }}>The page you're looking for doesn't exist.</p>
              </div>
            } />
          </Routes>
        </main>

        {/* Footer */}
        <footer style={{ borderTop: '1px solid var(--border)', padding: '32px 0', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.85rem' }}>
          <p>© 2026 Saloon — Premium Salon Booking Platform</p>
        </footer>

        <Toaster
          position="top-right"
          toastOptions={{
            style: {
              background: 'var(--bg-card)',
              color: 'var(--text)',
              border: '1px solid var(--border)',
              borderRadius: 'var(--radius-md)',
            },
          }}
        />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
