# Go Stok Takip Projesi

Bu proje Go dili ile yazılmış, PostgreSQL veritabanı kullanan web tabanlı bir stok takip ve sipariş uygulamasıdır.

Uygulamada ürün listeleme, sepet, admin paneli, kategori yönetimi, müşteri kaydı, müşteri girişi, adres kaydı, checkout ve sipariş geçmişi bulunur.

## Proje Linki

```text
https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi
```

## Klasör Yapısı

```text
.
|-- app/
|   |-- main.go
|   `-- templates/
|       |-- index.html
|       |-- systems.html
|       |-- cart.html
|       |-- checkout.html
|       |-- login.html
|       |-- admin.html
|       |-- musteri_giris.html
|       |-- musteri_kayit.html
|       `-- profil.html
|-- database/
|   |-- schema.sql
|   |-- stok_takip_backup.sql
|   `-- README.md
|-- go.mod
|-- go.sum
`-- README.md
```

## Gereksinimler

Yeni bilgisayarda şunlar kurulu olmalı:

- Go 1.21 veya üstü
- PostgreSQL
- Git

PostgreSQL 18 kullanıldı. Farklı PostgreSQL sürümü kullanıyorsan aşağıdaki komutlarda `18` yazan klasör adını kendi sürümüne göre değiştir.

Örnek:

```text
C:\Program Files\PostgreSQL\18\bin
```

PostgreSQL 17 ise:

```text
C:\Program Files\PostgreSQL\17\bin
```

## 1. Projeyi İndir

Terminal veya PowerShell aç:

```powershell
git clone https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi.git
cd Go-Programlama-Stok-Takip-Projesi
```

## 2. PostgreSQL Veritabanını Oluştur

Önce boş veritabanı oluştur:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

Komut şifre isterse PostgreSQL kurarken verdiğin `postgres` şifresini yaz.

Eğer `database "stok_takip" already exists` benzeri bir hata görürsen bu adımı geçebilirsin.

## 3. SQL Dosyasını İçeri Aktar

Tüm tabloları ve mevcut örnek verileri yüklemek için:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

Bu dosyanın içinde şu tablolar bulunur:

- `admin`
- `products`
- `categories`
- `musteriler`
- `adresler`
- `siparisler`
- `siparis_urunleri`

Sadece tablo yapısını kurmak istersen, veri olmadan şu dosyayı kullan:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\schema.sql"
```

Normal kullanım için önerilen dosya:

```text
database\stok_takip_backup.sql
```

## 4. Veritabanı Şifresini Kontrol Et

Uygulama PostgreSQL bağlantısını `app/main.go` içindeki `connStr` satırından alır.

Varsayılan bağlantı:

```go
connStr := "host=127.0.0.1 port=5432 user=postgres password=5757 dbname=stok_takip sslmode=disable"
```

Eğer PostgreSQL şifren `5757` değilse, `password=5757` kısmını kendi şifrenle değiştir.

Örnek:

```go
password=BENIM_SIFREM
```

Veritabanı adını değiştirme. Kod şu veritabanını bekliyor:

```text
stok_takip
```

## 5. Uygulamayı Çalıştır

Proje klasörünün içindeyken:

```powershell
go run ./app
```

Alternatif olarak:

```powershell
cd app
go run .
```

Başarılı olursa terminalde buna benzer çıktı görürsün:

```text
Stok takip sunucusu: http://localhost:8080
```

Sonra tarayıcıdan aç:

```text
http://localhost:8080
```

## Admin Paneli

Admin giriş sayfası:

```text
http://localhost:8080/login
```

Admin kullanıcı bilgileri SQL yedeğindeki `public.admin` tablosundan gelir.

## Sık Karşılaşılan Hatalar

### `psql` veya `createdb` bulunamadı

PostgreSQL yolu farklı olabilir. Şu klasörü kontrol et:

```text
C:\Program Files\PostgreSQL
```

Hangi sürüm varsa komuttaki `18` kısmını ona göre değiştir.

### `password authentication failed`

`app/main.go` içindeki `password=5757` kısmı PostgreSQL şifrenle aynı değil demektir. Kendi PostgreSQL şifreni yaz.

### `database "stok_takip" does not exist`

Veritabanı oluşturulmamış demektir. Şunu çalıştır:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

### `relation public.products does not exist`

SQL dosyası içeri aktarılmamış demektir. Şunu çalıştır:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

### `listen tcp :8080: bind`

8080 portu başka bir program tarafından kullanılıyor olabilir. O programı kapat veya `app/main.go` içindeki şu satırdan portu değiştir:

```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

Örnek:

```go
log.Fatal(http.ListenAndServe(":8081", nil))
```

Sonra şu adresten aç:

```text
http://localhost:8081
```

## Yapay Zekaya Verilecek Kurulum Görevi

Bu projeyi başka bir bilgisayarda yapay zekaya kurdurmak istersen aşağıdaki metni direkt verebilirsin:

```text
Bu Go + PostgreSQL stok takip projesini kur.

1. Repoyu indir:
   https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi

2. PostgreSQL'de stok_takip adında veritabanı oluştur.

3. database/stok_takip_backup.sql dosyasını stok_takip veritabanına import et.

4. app/main.go içindeki connStr satırında PostgreSQL şifresini kontrol et.

5. Proje kökünden go run ./app komutunu çalıştır.

6. Siteyi http://localhost:8080 adresinde aç.

Eğer psql veya createdb bulunamazsa C:\Program Files\PostgreSQL\<surum>\bin klasörünü kullan.
```

## Geliştirme Notları

- Ana Go kodu `app/main.go` dosyasındadır.
- HTML sayfaları `app/templates` klasöründedir.
- SQL dosyaları `database` klasöründedir.
- `database/schema.sql` sadece tablo yapısını içerir.
- `database/stok_takip_backup.sql` tablo yapısı ve mevcut verileri içerir.
- Derlenmiş `.exe` dosyaları GitHub'a yüklenmez.

## Test

Kodun derlenip derlenmediğini kontrol etmek için:

```powershell
go test ./...
```

Ek kontrol için:

```powershell
go vet ./...
```
