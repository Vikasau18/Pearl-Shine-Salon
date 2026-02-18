-- =============================================
-- SEED DATA FOR TESTING
-- =============================================

-- Password for all test users: "password123" (bcrypt hash)
-- $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

-- USERS
INSERT INTO users (id, email, password_hash, name, phone, role) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'customer@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Alice Johnson', '+1234567890', 'customer'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'customer2@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Bob Smith', '+1234567891', 'customer'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'owner1@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Maria Garcia', '+1234567892', 'salon_owner'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'owner2@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'James Wilson', '+1234567893', 'salon_owner'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'owner3@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Sarah Lee', '+1234567894', 'salon_owner'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'staff1@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Emma Davis', '+1234567895', 'staff'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'staff2@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Liam Brown', '+1234567896', 'staff'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'staff3@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Olivia Taylor', '+1234567897', 'staff');

-- SALONS
INSERT INTO salons (id, owner_id, name, address, city, state, lat, lng, phone, email, description, rating, total_reviews, image_url) VALUES
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Glamour Studio', '123 Main Street', 'New York', 'NY', 40.7128, -74.0060, '+1555000001', 'info@glamourstudio.com', 'Premium hair and beauty salon in the heart of Manhattan. We specialize in cuts, coloring, and styling for all hair types.', 4.80, 124, 'https://images.unsplash.com/photo-1560066984-138dadb4c035?w=800'),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'The Barber Shop', '456 Oak Avenue', 'Los Angeles', 'CA', 34.0522, -118.2437, '+1555000002', 'hello@thebarbershop.com', 'Classic barbershop with modern amenities. Specializing in mens grooming, fades, and hot towel shaves.', 4.60, 89, 'https://images.unsplash.com/photo-1585747860019-8587e tried-and-true?w=800'),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Zen Spa & Nails', '789 Elm Boulevard', 'Chicago', 'IL', 41.8781, -87.6298, '+1555000003', 'contact@zenspanails.com', 'Relax and rejuvenate at our full-service spa. Offering manicures, pedicures, facials, and massage therapy.', 4.90, 203, 'https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800');

-- SALON GALLERY
INSERT INTO salon_gallery (salon_id, image_url, caption, sort_order) VALUES
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://images.unsplash.com/photo-1522337360788-8b13dee7a37e?w=600', 'Our stylish interior', 1),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://images.unsplash.com/photo-1562322140-8baeececf3df?w=600', 'Hair coloring station', 2),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://images.unsplash.com/photo-1595476108010-b4d1f102b1b1?w=600', 'Premium products', 3),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'https://images.unsplash.com/photo-1503951914875-452162b0f3f1?w=600', 'Classic barber chair', 1),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'https://images.unsplash.com/photo-1599351431202-1e0f0137899a?w=600', 'Hot towel shave', 2),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=600', 'Relaxation area', 1),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'https://images.unsplash.com/photo-1519014816548-bf5fe059798b?w=600', 'Nail art samples', 2);

-- STAFF
INSERT INTO staff (id, salon_id, user_id, name, role, specialization, bio, avatar_url) VALUES
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Emma Davis', 'senior_stylist', 'Hair Coloring', 'Over 10 years experience in balayage and highlights.', 'https://randomuser.me/api/portraits/women/44.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'Sophie Martinez', 'stylist', 'Cutting & Styling', '5 years experience specializing in modern cuts.', 'https://randomuser.me/api/portraits/women/65.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'Ryan Chen', 'stylist', 'Mens Grooming', 'Expert in fades, pompadours, and beard grooming.', 'https://randomuser.me/api/portraits/men/32.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Liam Brown', 'barber', 'Classic Cuts', 'Master barber with 15 years in the business.', 'https://randomuser.me/api/portraits/men/75.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', NULL, 'Marcus Johnson', 'barber', 'Fades & Designs', 'Known for precision fades and hair art.', 'https://randomuser.me/api/portraits/men/22.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Olivia Taylor', 'nail_tech', 'Nail Art', '8 years specializing in gel nails and nail art.', 'https://randomuser.me/api/portraits/women/17.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', NULL, 'Mia Patel', 'therapist', 'Facials & Massage', 'Licensed esthetician and massage therapist.', 'https://randomuser.me/api/portraits/women/90.jpg'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', NULL, 'Nina Kim', 'nail_tech', 'Pedicures', 'Luxury pedicure specialist with 6 years experience.', 'https://randomuser.me/api/portraits/women/52.jpg');

-- STAFF WORKING HOURS (Mon=1 to Fri=5 for most staff)
INSERT INTO staff_working_hours (staff_id, day_of_week, start_time, end_time, is_off) VALUES
-- Emma Davis: Mon-Sat
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, '09:00', '18:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, '09:00', '18:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, '09:00', '18:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, '09:00', '18:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, '09:00', '18:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 6, '10:00', '16:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, '00:00', '00:00', true),
-- Liam Brown: Tue-Sat
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 1, '00:00', '00:00', true),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 2, '10:00', '19:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 3, '10:00', '19:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 4, '10:00', '19:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 5, '10:00', '19:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 6, '10:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 0, '00:00', '00:00', true),
-- Olivia Taylor: Mon-Fri
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 1, '09:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 2, '09:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 3, '09:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 4, '09:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 5, '09:00', '17:00', false),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 6, '00:00', '00:00', true),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 0, '00:00', '00:00', true);

-- SERVICES
INSERT INTO services (id, salon_id, name, description, price, duration_minutes, buffer_minutes, category) VALUES
-- Glamour Studio
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Womens Haircut', 'Precision cut with wash and blowout', 65.00, 45, 10, 'Hair'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Balayage Highlights', 'Natural-looking hand-painted highlights', 180.00, 120, 15, 'Hair'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Full Color', 'Root to tip single-process color', 120.00, 90, 10, 'Hair'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Blowout & Style', 'Wash, blowout and style', 45.00, 30, 5, 'Hair'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Mens Haircut', 'Classic cut with hot towel finish', 35.00, 30, 5, 'Hair'),
-- The Barber Shop
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Classic Cut', 'Traditional barber cut', 30.00, 30, 5, 'Barber'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Skin Fade', 'Precision skin fade with detailing', 40.00, 40, 5, 'Barber'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Hot Towel Shave', 'Luxury straight razor shave', 35.00, 30, 10, 'Barber'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a09', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Beard Trim & Shape', 'Precision beard grooming', 25.00, 20, 5, 'Barber'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Cut + Shave Combo', 'Haircut with hot towel shave', 55.00, 60, 10, 'Barber'),
-- Zen Spa & Nails
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Classic Manicure', 'Nail shaping, cuticle care, and polish', 30.00, 30, 5, 'Nails'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Gel Manicure', 'Long-lasting gel polish manicure', 45.00, 45, 5, 'Nails'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Luxury Pedicure', 'Foot soak, scrub, massage, and polish', 55.00, 60, 10, 'Nails'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Classic Facial', 'Deep cleansing facial with mask', 75.00, 60, 10, 'Spa'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Swedish Massage', 'Full body relaxation massage', 90.00, 60, 15, 'Spa'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Nail Art', 'Custom nail art design per nail', 10.00, 15, 0, 'Nails');

-- STAFF SERVICES (which staff can perform which services)
INSERT INTO staff_services (staff_id, service_id) VALUES
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a09'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13');

-- PROMO CODES
INSERT INTO promo_codes (id, salon_id, code, discount_percent, valid_from, valid_until, max_uses, is_active) VALUES
('10eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'WELCOME20', 20.00, '2025-01-01', '2026-12-31', 100, true),
('10eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'BARBER10', 10.00, '2025-06-01', '2026-06-30', 50, true),
('10eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'SPADAY25', 25.00, '2025-01-01', '2026-12-31', 200, true);

-- SAMPLE APPOINTMENTS
INSERT INTO appointments (id, customer_id, salon_id, staff_id, service_id, appointment_date, start_time, end_time, status) VALUES
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '2026-02-15', '10:00', '10:45', 'confirmed'),
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '2026-02-18', '14:00', '14:45', 'pending'),
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06', '2026-02-12', '11:00', '11:30', 'completed'),
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', '2026-01-20', '13:00', '15:00', 'completed');

-- SAMPLE PAYMENTS
INSERT INTO payments (appointment_id, amount, discount, tax, total, method, status, receipt_number) VALUES
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 30.00, 0, 2.40, 32.40, 'card', 'completed', 'RCP-2026-0001'),
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 180.00, 0, 14.40, 194.40, 'card', 'completed', 'RCP-2026-0002');

-- SAMPLE REVIEWS
INSERT INTO reviews (customer_id, salon_id, staff_id, appointment_id, rating, comment) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '20eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 5, 'Best haircut I have ever had! Liam is amazing.'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '20eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 5, 'Emma did a fantastic balayage. Absolutely love the result!');

-- SAMPLE FAVORITES
INSERT INTO favorites (user_id, salon_id) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12');
