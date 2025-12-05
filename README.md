Jessica Giovanna Chandra

## 🔗 Repository
**GitHub:** https://github.com/zalg2261/gVII

## Deskripsi Proyek

Platform pembelian tiket bioskop online yang memungkinkan customer melakukan transaksi kapanpun secara online. Sistem ini memastikan nomor kursi tidak akan digunakan oleh orang lain dengan sistem locking dan timeout pembayaran 10 menit.

## Fitur Utama

1. **Browse Film Tanpa Login** - User dapat melihat daftar film dan sinopsis tanpa perlu login
2. **Login/Register** - Sistem autentikasi dengan JWT
3. **Pemesanan Tiket** - User dapat memesan tiket setelah login
4. **Sistem Locking Kursi** - Kursi di-hold selama 10 menit saat booking
5. **Pembayaran dengan Timer** - Timer 10 menit untuk menyelesaikan pembayaran
6. **Auto Release Kursi** - Kursi otomatis dikembalikan jika pembayaran tidak selesai
7. **Refund Otomatis** - Sistem refund otomatis jika bioskop membatalkan showtime
8. **Admin CRUD** - Admin dapat mengelola jadwal tayang dan film

## Teknologi yang Digunakan

### Backend
- **Golang** 1.23
- **Fiber** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **JWT** - Authentication

### Frontend
- **Next.js** 16
- **React** 19
- **TypeScript**
- **Tailwind CSS** 4

## Struktur Proyek

```
bioskop/
├── backend/
│   ├── internal/
│   │   ├── controllers/    # API handlers
│   │   ├── models/         # Database models
│   │   ├── routes/         # Route definitions
│   │   ├── middleware/     # Auth middleware
│   │   ├── services/       # Business logic
│   │   └── db/             # Database connection
│   ├── migrations/         # SQL migration files
│   └── main.go            # Entry point
├── frontend/
│   ├── app/               # Next.js app directory
│   ├── src/
│   │   ├── pages/         # Page components
│   │   └── services/      # API service
│   └── package.json
├── SYSTEM_DESIGN.md       # Dokumentasi System Design
├── DATABASE_DESIGN.md     # Dokumentasi Database Design
└── README.md              # File ini
```

## Instalasi dan Setup

### Prerequisites

- Go 1.23 atau lebih baru
- Node.js 18+ dan npm
- PostgreSQL 12+

### 1. Setup Database

```bash
# Buat database PostgreSQL
createdb db_bioskop

# Atau menggunakan psql
psql -U postgres
CREATE DATABASE db_bioskop;

# Import schema
psql -U postgres -d db_bioskop -f backend/migrations/001_complete_schema.sql
```

### 2. Setup Backend

```bash
cd backend

# Install dependencies
go mod download

# Buat file .env
cat > .env << EOF
PORT=4000
JWT_SECRET=your-secret-key-change-this
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=db_bioskop
DB_PORT=5432
EOF

# Update connection string di internal/db/connection.go sesuai konfigurasi Anda
# Atau gunakan environment variables

# Run server
go run main.go
```

Backend akan berjalan di `http://localhost:4000`

### 3. Setup Frontend

```bash
cd frontend

# Install dependencies
npm install

# Buat file .env.local (optional)
echo "NEXT_PUBLIC_API_URL=http://localhost:4000" > .env.local

# Run development server
npm run dev
```

Frontend akan berjalan di `http://localhost:3000`

## API Endpoints

### Public Endpoints

- `GET /movies` - Daftar semua film
- `GET /movies/:id` - Detail film
- `GET /schedule` - Daftar jadwal tayang (dapat filter dengan `?city=`, `?branch_id=`, `?movie_id=`)
- `GET /schedule/:id` - Detail jadwal tayang
- `GET /branches` - Daftar cabang bioskop
- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login user

### Protected Endpoints (Require JWT Token)

- `POST /book` - Buat booking (reserve kursi)
- `POST /payment/:bookingId` - Selesaikan pembayaran
- `POST /payment/failed/:bookingId` - Batalkan pembayaran
- `GET /my-bookings` - Daftar booking user
- `POST /wallet/topup` - Top up wallet (optional)

### Admin Endpoints (Require Admin Role)

- `POST /admin/schedule` - Buat jadwal tayang
- `PUT /admin/schedule/:id` - Update jadwal tayang
- `DELETE /admin/schedule/:id` - Hapus jadwal tayang
- `POST /admin/schedule/:id/cancel` - Batalkan showtime dan refund semua booking
- `POST /admin/movies` - Buat film baru
- `PUT /admin/movies/:id` - Update film
- `DELETE /admin/movies/:id` - Hapus film
- `POST /admin/refund/:bookingId` - Refund booking tertentu
- `GET /admin/refunds` - Daftar semua refund

## Cara Menggunakan

### Sebagai User

1. Buka `http://localhost:3000`
2. Browse film tanpa perlu login
3. Klik "Pesan Tiket" pada film yang diinginkan
4. Login atau register jika belum
5. Pilih jadwal tayang
6. Pilih jumlah kursi
7. Selesaikan pembayaran dalam 10 menit
8. Lihat booking di "My Bookings"

### Sebagai Admin

1. Login dengan akun admin (default: `admin@bioskop.com` / `admin123`)
2. Gunakan Postman atau API client untuk akses endpoint admin
3. CRUD jadwal tayang dan film melalui API

## Testing dengan Postman

Import file `Bioskop_API.postman_collection.json` ke Postman untuk mendapatkan semua endpoint yang sudah dikonfigurasi.

### Setup Postman Environment

1. Buat environment baru di Postman
2. Tambahkan variable:
   - `base_url`: `http://localhost:4000`
   - `token`: (akan diisi otomatis setelah login)

## Alur Sistem

### 1. Booking Flow

```
User Browse Film → Pilih Film → Pilih Jadwal → Login/Register → 
Pilih Jumlah Kursi → Reserve (Lock 10 menit) → Halaman Pembayaran → 
Selesaikan Pembayaran → Kursi Permanen Locked
```

### 2. Timeout Flow

```
Booking Created → Timer 10 Menit → Jika Tidak Bayar → 
Auto Cancel → Release Kursi → Status: CANCELLED
```

### 3. Refund Flow (Cinema Issue)

```
Admin Cancel Showtime → Sistem Cari Semua PAID Booking → 
Update Status: REFUNDED → Release Kursi → 
Buat Refund Record → Notifikasi User
```

## Background Jobs

Sistem memiliki background job yang berjalan setiap 1 menit untuk:
- Cleanup expired bookings (PENDING yang sudah melewati 10 menit)
- Release kursi yang expired
- Hapus seat locks yang expired

## Database Schema

Lihat `DATABASE_DESIGN.md` untuk detail lengkap schema database.

Tabel utama:
- `users` - Data user dan admin
- `movies` - Data film
- `branches` - Data cabang bioskop
- `showtimes` - Jadwal tayang
- `bookings` - Data pemesanan
- `seat_locks` - Tracking kursi yang di-hold
- `transactions` - Data transaksi
- `refunds` - Data refund

## System Design

Lihat `SYSTEM_DESIGN.md` untuk:
- Flowchart sistem
- Solusi pemilihan tempat duduk
- Sistem restok tiket
- Alur refund/pembatalan
- Arsitektur sistem

## Catatan Penting

1. **JWT Secret**: Pastikan mengubah JWT_SECRET di production
2. **Database Connection**: Update connection string sesuai konfigurasi PostgreSQL Anda
3. **CORS**: Jika frontend dan backend di domain berbeda, perlu setup CORS
4. **Environment Variables**: Gunakan .env untuk konfigurasi sensitif

## Troubleshooting

### Backend tidak bisa connect ke database
- Pastikan PostgreSQL running
- Check connection string di `internal/db/connection.go`
- Pastikan database `db_bioskop` sudah dibuat

### Frontend tidak bisa connect ke backend
- Pastikan backend running di port 4000
- Check `NEXT_PUBLIC_API_URL` di frontend
- Check CORS settings di backend (jika perlu)

### Booking tidak otomatis expired
- Pastikan background job running (check logs)
- Manual cleanup bisa dilakukan dengan call function `CleanupExpiredBookings()`

## Kontribusi

Proyek ini dibuat untuk keperluan skill test. Untuk production, perlu penambahan:
- Unit tests
- Integration tests
- Error handling yang lebih robust
- Logging system
- Rate limiting
- CORS configuration
- Security enhancements

## License

Proyek ini dibuat untuk keperluan skill test.

