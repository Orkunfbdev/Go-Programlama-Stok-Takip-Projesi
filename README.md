# Go Stok Takip Projesi

PostgreSQL destekli web tabanli stok takip ve siparis uygulamasi.

## Klasor Yapisi

```text
.
├── app/
│   ├── main.go
│   └── templates/
│       ├── index.html
│       ├── systems.html
│       ├── cart.html
│       ├── checkout.html
│       ├── login.html
│       ├── admin.html
│       ├── musteri_giris.html
│       ├── musteri_kayit.html
│       └── profil.html
├── database/
│   ├── schema.sql
│   └── stok_takip_backup.sql
├── go.mod
├── go.sum
└── README.md
```

## Gereksinimler

- Go 1.21 veya ustu
- PostgreSQL 18

## Veritabani Kurulumu

Yeni bilgisayarda once veritabanini olustur:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

Tum tablo ve verileri yuklemek icin:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

Sadece tablo yapisini kurmak istersen:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\schema.sql"
```

## Baglanti Ayari

PostgreSQL sifren farkliysa `app/main.go` icindeki `connStr` satirinda `password=5757` kismini kendi sifrenle degistir.

## Calistirma

Repo kokunden:

```powershell
go run ./app
```

Ya da uygulama klasorunden:

```powershell
cd app
go run .
```

Tarayicida ac:

```text
http://localhost:8080
```

## Ozellikler

- Urun listeleme
- Sepet sistemi
- Admin paneli
- Kategori yonetimi
- Musteri kayit ve giris
- Adres kaydi
- Checkout ve siparis olusturma
- Siparis gecmisi
- PostgreSQL tablo kurulumu ve SQL yedegi
