# 🏪 Saloon - Premium Salon Booking Platform

A full-stack salon appointment web application with Go backend, React frontend, and PostgreSQL database.

## ✨ Features

### Customer
- Browse & search salons (location-based)
- View salon details, services, staff, gallery, reviews
- Book appointments with real-time availability
- Reschedule or cancel appointments
- Favorite salons
- Promo codes & discounts
- Notifications & reminders
- Waitlist system

### Salon Owner
- Multi-tenant dashboard (manage multiple salons)
- Analytics: revenue, popular services, peak hours
- Appointment management (complete, no-show)
- Staff & service management (via API)
- Payment tracking & digital receipts
- Promo code management

## 🏗️ Tech Stack

| Layer      | Technology                    |
|------------|-------------------------------|
| Backend    | Go + Gin Framework            |
| Frontend   | React + Vite                  |
| Database   | PostgreSQL                    |
| Auth       | JWT + bcrypt                  |
| Styling    | Custom CSS (Dark theme)       |
| Deployment | Docker + Docker Compose       |

## 🚀 Quick Start

### With Docker (Recommended)

```bash
docker-compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api
- Database is auto-initialized with schema + seed data

### Without Docker

#### 1. Database Setup
```bash
# Start PostgreSQL and create database
psql -U postgres -c "CREATE DATABASE saloon;"
psql -U postgres -d saloon -f db/migrations/001_schema.up.sql
psql -U postgres -d saloon -f db/seed.sql
```

#### 2. Backend
```bash
cd backend
cp .env.example .env  # Edit with your database URL
go mod tidy
go run main.go
# Server starts on :8080
```

#### 3. Frontend
```bash
cd frontend
npm install
npm run dev
# Dev server starts on :5173
```

## 🔐 Demo Accounts

| Role         | Email              | Password     |
|-------------|--------------------| -------------|
| Customer    | customer@test.com  | password123  |
| Customer    | jane@test.com      | password123  |
| Salon Owner | owner1@test.com    | password123  |
| Salon Owner | owner2@test.com    | password123  |
| Admin       | admin@test.com     | password123  |

## 📡 API Endpoints

### Auth
- `POST /api/auth/register` - Register
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Get profile
- `PUT /api/auth/profile` - Update profile

### Salons (Public)
- `GET /api/salons` - List salons (search, city, pagination)
- `GET /api/salons/:id` - Salon details
- `GET /api/salons/:id/services` - Services list
- `GET /api/salons/:id/staff` - Staff list
- `GET /api/salons/:id/reviews` - Reviews

### Customer (Authenticated)
- `POST /api/appointments` - Book appointment
- `GET /api/appointments` - My appointments
- `PUT /api/appointments/:id/reschedule` - Reschedule
- `PUT /api/appointments/:id/cancel` - Cancel
- `POST /api/favorites/:salon_id` - Toggle favorite
- `GET /api/favorites` - My favorites
- `POST /api/reviews` - Create review
- `GET /api/promos/validate` - Validate promo
- `POST /api/waitlist` - Join waitlist
- `GET /api/notifications` - Notifications

### Dashboard (Salon Owner)
- `POST /api/dashboard/salons` - Create salon
- `GET /api/dashboard/salons` - My salons
- `GET /api/dashboard/salons/:id/analytics` - Analytics
- `GET /api/dashboard/salons/:id/appointments` - Appointments
- CRUD for services, staff, promos, payments

## 📁 Project Structure

```
saloon/
├── backend/
│   ├── config/          # Environment configuration
│   ├── database/        # PostgreSQL connection
│   ├── handlers/        # API request handlers
│   ├── middleware/       # Auth, CORS, multi-tenant
│   ├── models/          # Data models & DTOs
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic (booking, notifications)
│   └── main.go          # Entry point
├── frontend/
│   └── src/
│       ├── api/         # Axios client
│       ├── components/  # Reusable components
│       ├── context/     # Auth context
│       └── pages/       # Page components
├── db/
│   ├── migrations/      # SQL schema
│   └── seed.sql         # Sample data
├── docker-compose.yml
└── README.md
```

## 🔒 Security

- JWT-based authentication with role-based access control
- Bcrypt password hashing
- Multi-tenant data isolation (salon_id scoping)
- Concurrency-safe booking with PostgreSQL `SELECT ... FOR UPDATE`
- CORS middleware
- SQL injection prevention via parameterized queries

## 📝 License

MIT
