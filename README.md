# Fleetify Logistics - Invoice System Technical Test

Sistem manajemen invoice logistik yang dibangun dengan **Go (Golang)** untuk Backend dan **Next.js (TypeScript)** untuk Frontend. Project ini mengimplementasikan sistem *multi-step wizard* dengan state management yang persisten, keamanan *Zero-Trust*, dan fitur pencetakan resi resmi.

## 🚀 Fitur Utama

- **Zero-Setup Dockerization**: Menjalankan seluruh stack (FE, BE, Database) hanya dengan satu perintah.
- **Role-Based Access Control**:
  - **Admin**: Akses penuh termasuk rincian harga dan total finansial.
  - **Kerani**: Akses terbatas hanya pada pengelolaan item/kuantitas (Harga disanitasi dari payload & UI).
- **Multi-Step Wizard**: Form interaktif dengan validasi dan state yang tidak hilang saat refresh (Zustand Persist).
- **Zero-Trust Backend**: Backend menghitung ulang seluruh total harga berdasarkan data master di database, mengabaikan manipulasi harga dari sisi client.
- **Print-Ready Layout**: Layout resi resmi yang otomatis menyembunyikan elemen UI web saat dicetak (CTRL+P).

## 🛠️ Tech Stack

- **Frontend**: Next.js 15, Tailwind CSS, Zustand, Axios.
- **Backend**: Golang 1.22, GORM, Gin Framework, JWT.
- **Database**: PostgreSQL 16.

---

## 📦 Instruksi Menjalankan (Docker)

Pastikan Anda sudah menginstal **Docker** dan **Docker Compose** di mesin Anda.

### 1. Jalankan Project
Cukup jalankan perintah berikut di direktori utama project:

```bash
docker-compose up --build
