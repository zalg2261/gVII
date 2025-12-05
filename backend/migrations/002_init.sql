-- ===============================
-- 002_ini.sql
-- DATABASE BIOSKOP ONLINE
-- ===============================

-- ===============================
-- DROP TABLE IF EXISTS (biar clean)
-- ===============================
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS showtimes CASCADE;
DROP TABLE IF EXISTS movies CASCADE;
DROP TABLE IF EXISTS branches CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ===============================
-- TABLE: users
-- login + auth
-- ===============================
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(120) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    balance     BIGINT DEFAULT 0,     -- optional (saldo)
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: branches
-- lokasi / kota bioskop
-- ===============================
CREATE TABLE branches (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(150) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    address     VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: movies
-- ===============================
CREATE TABLE movies (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    genre       VARCHAR(100),
    duration    INT,               -- menit
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: showtimes
-- jadwal tayang
-- ===============================
CREATE TABLE showtimes (
    id          SERIAL PRIMARY KEY,
    movie_id    INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    branch_id   INT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    studio      VARCHAR(50) NOT NULL,
    show_time   TIMESTAMP NOT NULL,
    seats_left  INT NOT NULL DEFAULT 50, -- default 50 kursi
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: bookings
-- pending / paid / failed / cancelled / refunded
-- ===============================
CREATE TABLE bookings (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    showtime_id INT NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    seats       INT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    expired_at  TIMESTAMP,                -- auto cancel 10 menit
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- TABLE: transactions
-- catatan pembayaran
-- ===============================
CREATE TABLE transactions (
    id          SERIAL PRIMARY KEY,
    booking_id  INT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    amount      BIGINT NOT NULL,
    status      VARCHAR(20) NOT NULL,     -- SUCCESS / FAILED / REFUND
    created_at  TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- DEFAULT MOVIE PRICE
-- harga flat (bisa kamu isi atau tidak)
-- ===============================
INSERT INTO movies (title, genre, duration)
VALUES 
('Movie Contoh', 'Action', 120),
('Movie Tes', 'Drama', 90);

INSERT INTO branches (name, city, address)
VALUES
('XXI Bandung PVJ', 'Bandung', 'Jl. Sukajadi'),
('XXI Jakarta Kota Kasablanka', 'Jakarta', 'Jl. Casablanca');

