-- ===============================
-- COMPLETE DATABASE SCHEMA
-- BIOSKOP ONLINE TICKETING SYSTEM
-- ===============================

-- ===============================
-- DROP EXISTING TABLES (CLEAN SLATE)
-- ===============================
DROP TABLE IF EXISTS refunds CASCADE;
DROP TABLE IF EXISTS seat_locks CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS showtimes CASCADE;
DROP TABLE IF EXISTS movies CASCADE;
DROP TABLE IF EXISTS branches CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ===============================
-- TABLE: users
-- ===============================
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(120) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    role        VARCHAR(20) DEFAULT 'user', -- 'user' or 'admin'
    balance     BIGINT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- ===============================
-- TABLE: branches
-- ===============================
CREATE TABLE branches (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(150) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    address     VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_branches_city ON branches(city);

-- ===============================
-- TABLE: movies
-- ===============================
CREATE TABLE movies (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    genre       VARCHAR(100),
    duration    INT, -- dalam menit
    synopsis    TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: showtimes (jadwal tayang)
-- ===============================
CREATE TABLE showtimes (
    id          SERIAL PRIMARY KEY,
    movie_id    INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    branch_id   INT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    studio      VARCHAR(50) NOT NULL,
    show_time   TIMESTAMP NOT NULL,
    seats_total INT NOT NULL DEFAULT 50,
    seats_left  INT NOT NULL DEFAULT 50,
    price       BIGINT NOT NULL DEFAULT 50000, -- harga per tiket dalam rupiah
    status      VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, CANCELLED
    created_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_seats_left CHECK (seats_left >= 0),
    CONSTRAINT chk_seats_total CHECK (seats_total > 0)
);

CREATE INDEX idx_showtimes_movie_id ON showtimes(movie_id);
CREATE INDEX idx_showtimes_branch_id ON showtimes(branch_id);
CREATE INDEX idx_showtimes_show_time ON showtimes(show_time);
CREATE INDEX idx_showtimes_status ON showtimes(status);

-- ===============================
-- TABLE: bookings
-- ===============================
CREATE TABLE bookings (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    showtime_id INT NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    seats       INT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PAID, CANCELLED, REFUNDED
    expires_at  TIMESTAMP, -- untuk timeout 10 menit
    created_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_seats_positive CHECK (seats > 0),
    CONSTRAINT chk_booking_status CHECK (status IN ('PENDING', 'PAID', 'CANCELLED', 'REFUNDED'))
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_showtime_id ON bookings(showtime_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
CREATE INDEX idx_bookings_showtime_status ON bookings(showtime_id, status);

-- ===============================
-- TABLE: seat_locks
-- ===============================
CREATE TABLE seat_locks (
    id          SERIAL PRIMARY KEY,
    showtime_id INT NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    booking_id  INT REFERENCES bookings(id) ON DELETE CASCADE,
    seats_count INT NOT NULL DEFAULT 1,
    locked_at   TIMESTAMP DEFAULT NOW(),
    expires_at  TIMESTAMP NOT NULL,
    CONSTRAINT chk_seats_count_positive CHECK (seats_count > 0)
);

CREATE INDEX idx_seat_locks_showtime_id ON seat_locks(showtime_id);
CREATE INDEX idx_seat_locks_booking_id ON seat_locks(booking_id);
CREATE INDEX idx_seat_locks_expires_at ON seat_locks(expires_at);

-- ===============================
-- TABLE: transactions
-- ===============================
CREATE TABLE transactions (
    id          SERIAL PRIMARY KEY,
    booking_id  INT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    amount      BIGINT NOT NULL,
    status      VARCHAR(20) NOT NULL, -- SUCCESS, FAILED, REFUND
    created_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_transaction_status CHECK (status IN ('SUCCESS', 'FAILED', 'REFUND'))
);

CREATE INDEX idx_transactions_booking_id ON transactions(booking_id);
CREATE INDEX idx_transactions_status ON transactions(status);

-- ===============================
-- TABLE: refunds
-- ===============================
CREATE TABLE refunds (
    id          SERIAL PRIMARY KEY,
    booking_id  INT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_cents BIGINT NOT NULL,
    reason      TEXT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'REQUESTED', -- REQUESTED, APPROVED, FAILED
    created_at  TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP,
    CONSTRAINT chk_refund_amount_positive CHECK (amount_cents > 0),
    CONSTRAINT chk_refund_status CHECK (status IN ('REQUESTED', 'APPROVED', 'FAILED'))
);

CREATE INDEX idx_refunds_booking_id ON refunds(booking_id);
CREATE INDEX idx_refunds_user_id ON refunds(user_id);
CREATE INDEX idx_refunds_status ON refunds(status);

-- ===============================
-- SAMPLE DATA
-- ===============================

-- Insert sample movies
INSERT INTO movies (title, genre, duration, synopsis) VALUES
('Avengers: Endgame', 'Action', 181, 'Setelah peristiwa Infinity War, para pahlawan yang tersisa harus bekerja sama untuk mengembalikan keseimbangan alam semesta.'),
('The Dark Knight', 'Action', 152, 'Batman menghadapi Joker yang mengancam Gotham City dengan kekacauan dan teror.'),
('Inception', 'Sci-Fi', 148, 'Seorang pencuri yang masuk ke dalam mimpi orang lain untuk mencuri rahasia dari alam bawah sadar mereka.'),
('Parasite', 'Thriller', 132, 'Keluarga miskin yang menyusup ke dalam kehidupan keluarga kaya dengan cara yang tidak terduga.'),
('Interstellar', 'Sci-Fi', 169, 'Tim penjelajah antariksa melakukan perjalanan melalui lubang cacing untuk mencari planet baru yang bisa dihuni.');

-- Insert sample branches
INSERT INTO branches (name, city, address) VALUES
('XXI Bandung PVJ', 'Bandung', 'Jl. Sukajadi No. 123, Bandung'),
('XXI Jakarta Kota Kasablanka', 'Jakarta', 'Jl. Casablanca Raya No. 88, Jakarta Selatan'),
('CGV Grand Indonesia', 'Jakarta', 'Jl. MH Thamrin No. 1, Jakarta Pusat'),
('XXI Surabaya Tunjungan Plaza', 'Surabaya', 'Jl. Basuki Rahmat No. 8-12, Surabaya'),
('Cinema 21 Yogyakarta Malioboro', 'Yogyakarta', 'Jl. Malioboro No. 52-58, Yogyakarta');

-- Insert sample showtimes
INSERT INTO showtimes (movie_id, branch_id, studio, show_time, seats_total, seats_left, price, status) VALUES
(1, 1, 'Studio 1', NOW() + INTERVAL '2 days' + INTERVAL '14 hours', 50, 50, 50000, 'ACTIVE'),
(1, 1, 'Studio 2', NOW() + INTERVAL '2 days' + INTERVAL '17 hours', 50, 50, 50000, 'ACTIVE'),
(2, 2, 'Studio 1', NOW() + INTERVAL '1 day' + INTERVAL '19 hours', 50, 50, 45000, 'ACTIVE'),
(3, 2, 'Studio 3', NOW() + INTERVAL '3 days' + INTERVAL '15 hours', 50, 50, 55000, 'ACTIVE'),
(4, 1, 'Studio 3', NOW() + INTERVAL '1 day' + INTERVAL '20 hours', 50, 50, 40000, 'ACTIVE');

-- Insert sample admin user (password: admin123)
-- Password hash untuk "admin123" menggunakan bcrypt
INSERT INTO users (name, email, password, role) VALUES
('Admin Bioskop', 'admin@bioskop.com', '$2a$10$rKqXqXqXqXqXqXqXqXqXeKqXqXqXqXqXqXqXqXqXqXqXqXqXqXq', 'admin');

-- ===============================
-- FUNCTION: Auto cleanup expired bookings
-- ===============================
CREATE OR REPLACE FUNCTION cleanup_expired_bookings()
RETURNS void AS $$
BEGIN
    -- Update expired PENDING bookings to CANCELLED
    UPDATE bookings
    SET status = 'CANCELLED'
    WHERE status = 'PENDING'
      AND expires_at < NOW();
    
    -- Release seats for expired bookings
    UPDATE showtimes s
    SET seats_left = s.seats_left + b.seats
    FROM bookings b
    WHERE b.showtime_id = s.id
      AND b.status = 'CANCELLED'
      AND b.expires_at < NOW()
      AND NOT EXISTS (
          SELECT 1 FROM bookings b2
          WHERE b2.id = b.id AND b2.status != 'CANCELLED'
      );
    
    -- Delete expired seat locks
    DELETE FROM seat_locks
    WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- ===============================
-- TRIGGER: Update seats_left when booking status changes
-- ===============================
CREATE OR REPLACE FUNCTION update_seats_on_booking_change()
RETURNS TRIGGER AS $$
BEGIN
    -- If booking is cancelled or refunded, restore seats
    IF NEW.status IN ('CANCELLED', 'REFUNDED') AND OLD.status = 'PAID' THEN
        UPDATE showtimes
        SET seats_left = seats_left + NEW.seats
        WHERE id = NEW.showtime_id;
    END IF;
    
    -- If booking becomes PAID, ensure seats are already reduced (handled in application)
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_seats_on_booking
AFTER UPDATE OF status ON bookings
FOR EACH ROW
EXECUTE FUNCTION update_seats_on_booking_change();

-- ===============================
-- VIEW: Available showtimes with movie and branch info
-- ===============================
CREATE OR REPLACE VIEW v_available_showtimes AS
SELECT 
    s.id,
    s.studio,
    s.show_time,
    s.seats_left,
    s.seats_total,
    s.price,
    s.status,
    m.title as movie_title,
    m.genre,
    m.duration,
    m.synopsis,
    b.name as branch_name,
    b.city,
    b.address
FROM showtimes s
JOIN movies m ON s.movie_id = m.id
JOIN branches b ON s.branch_id = b.id
WHERE s.status = 'ACTIVE'
  AND s.seats_left > 0
  AND s.show_time > NOW()
ORDER BY s.show_time ASC;

-- ===============================
-- COMMENTS
-- ===============================
COMMENT ON TABLE users IS 'Tabel untuk menyimpan data user dan admin';
COMMENT ON TABLE movies IS 'Tabel untuk menyimpan data film';
COMMENT ON TABLE branches IS 'Tabel untuk menyimpan data cabang bioskop';
COMMENT ON TABLE showtimes IS 'Tabel untuk menyimpan jadwal tayang film';
COMMENT ON TABLE bookings IS 'Tabel untuk menyimpan data pemesanan tiket';
COMMENT ON TABLE seat_locks IS 'Tabel untuk tracking kursi yang sedang di-hold';
COMMENT ON TABLE transactions IS 'Tabel untuk menyimpan data transaksi pembayaran';
COMMENT ON TABLE refunds IS 'Tabel untuk menyimpan data refund/pembatalan';

