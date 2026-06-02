# 📦 Go Stok Takip Projesi

Go ile yazılmış, PostgreSQL destekli web tabanlı stok yönetim sistemi.

## 🚀 Kurulum

### 1. Gereksinimler
- [Go](https://go.dev/dl/) (1.21+)
- [PostgreSQL](https://www.postgresql.org/download/) (18)

### 2. Projeyi İndir
```bash
git clone https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi.git
cd Go-Programlama-Stok-Takip-Projesi
```

### 3. Veritabanı Kurulumu

Eğer `stok_takip` database'i yoksa oluştur:
```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

> Eğer database zaten varsa bu adımı atla.

Backup'ı import et:
```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "stok_takip_backup.sql"
```

Veya pgAdmin kullanıyorsan:
> `stok_takip` → Sağ tık → **Restore** → `stok_takip_backup.sql` dosyasını seç → Restore

### 4. Bağlantı Ayarı

`main.go` dosyasının 63. satırındaki şifreyi kendi PostgreSQL şifrenle değiştir:
```go
connStr := "host=127.0.0.1 port=5432 user=postgres password=ŞİFREN dbname=stok_takip sslmode=disable"
```

### 5. Çalıştır
```powershell
go run main.go
```

Tarayıcıda aç: **http://localhost:8080**

---

## 🔑 Admin Paneli

- URL: `http://localhost:8080/login`
- Admin hesabı veritabanındaki `public.admin` tablosunda tutulur

---

## 📁 Proje Yapısı

```
├── main.go                 # Ana Go dosyası
├── go.mod / go.sum         # Bağımlılıklar
├── stok_takip_backup.sql   # Veritabanı yedeği
└── templates/
    ├── index.html          # Ana sayfa
    ├── systems.html        # Ürünler sayfası
    ├── cart.html           # Sepet
    ├── login.html          # Admin girişi
    └── admin.html          # Admin paneli
```

---

## 🛠️ Özellikler

- ✅ Ürün listeleme ve arama
- ✅ Sepet sistemi
- ✅ Admin paneli (ürün ekle/düzenle/sil)
- ✅ Kategori yönetimi (DB'den)
- ✅ Stok takibi ve uyarıları
- ✅ PostgreSQL entegrasyonu
