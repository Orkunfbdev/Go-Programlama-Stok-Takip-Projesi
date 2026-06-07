# Go Programlama Stok Takip Projesi

Bu proje Go dili ile yazılmış, PostgreSQL veritabanı kullanan web tabanlı bir stok takip ve sipariş uygulamasıdır.

Projede ana sayfa, ürün listeleme, kategori filtreleme, sepet, stok kontrolü, satın alma, müşteri kayıt/giriş, adres kaydı, sipariş geçmişi ve admin paneli bulunur. Admin panelinden ürün ve kategori eklenebilir; eklenen ürünler listeleme sayfasına ve filtrelere otomatik gelir.

## GitHub

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

Bu projede PostgreSQL 18 kullanıldı. Başka sürüm varsa komutlardaki `18` kısmını kendi PostgreSQL sürümüne göre değiştir.

Örnek PostgreSQL yolu:

```text
C:\Program Files\PostgreSQL\18\bin
```

PostgreSQL 17 kuruluysa:

```text
C:\Program Files\PostgreSQL\17\bin
```

## Hızlı Kurulum

PowerShell aç ve şu adımları sırayla uygula.

### 1. Projeyi İndir

```powershell
git clone https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi.git
cd Go-Programlama-Stok-Takip-Projesi
```

### 2. Veritabanını Oluştur

PostgreSQL şifren sorulursa kurulumda verdiğin `postgres` kullanıcısının şifresini yaz.

```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

Eğer `database "stok_takip" already exists` hatası gelirse ve temiz kurulum yapmak istiyorsan dikkatli şekilde şu komutları kullanabilirsin:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\dropdb.exe" -U postgres --if-exists stok_takip
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

Bu işlem eski `stok_takip` veritabanını siler. İçinde önemli veri varsa silme.

### 3. SQL Yedeğini İçeri Aktar

Tam çalışan proje verilerini yüklemek için bu dosyayı kullan:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

Bu yedek dosyası tablo yapısını ve mevcut örnek verileri birlikte kurar.

Yedekte bulunan ana tablolar:

- `admin`
- `products`
- `categories`
- `musteriler`
- `adresler`
- `siparisler`
- `siparis_urunleri`

Sadece boş tablo yapısı istenirse şu dosya kullanılabilir:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\schema.sql"
```

Normal kurulum için önerilen dosya `database\stok_takip_backup.sql` dosyasıdır.

### 4. Veritabanı Şifresini Ayarla

Uygulama varsayılan olarak şu bağlantı bilgileriyle çalışır:

```text
host=127.0.0.1
port=5432
user=postgres
password=5757
dbname=stok_takip
sslmode=disable
```

Arkadaşının PostgreSQL şifresi `5757` değilse kodu değiştirmeden PowerShell'de şu komutu çalıştırabilir:

```powershell
$env:DB_PASSWORD="POSTGRES_SIFRESI"
```

Örnek:

```powershell
$env:DB_PASSWORD="1234"
```

İstersen diğer ayarlar da ortam değişkeniyle verilebilir:

```powershell
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="5757"
$env:DB_NAME="stok_takip"
$env:DB_SSLMODE="disable"
```

Tek satır bağlantı adresi kullanmak istersen:

```powershell
$env:DATABASE_URL="postgres://postgres:5757@127.0.0.1:5432/stok_takip?sslmode=disable"
```

### 5. Go Paketlerini İndir

```powershell
go mod download
```

### 6. Uygulamayı Çalıştır

Proje klasörünün içindeyken:

```powershell
go run ./app
```

Başarılı olursa terminalde şu çıktı görünür:

```text
Stok takip sunucusu: http://localhost:8080
```

Tarayıcıdan aç:

```text
http://localhost:8080
```

## Giriş Bilgileri

Admin paneli:

```text
http://localhost:8080/login
```

SQL yedeğinde gelen admin hesabı:

```text
Kullanıcı adı: orkunorkun
Şifre: 2005
```

Müşteri tarafında yeni hesap oluşturmak için:

```text
http://localhost:8080/musteri/kayit
```

## Projede Ne Var?

- Ana sayfada öne çıkan ürünler gösterilir.
- `Sistemler` sayfasında ürünler listelenir.
- Kategoriler veritabanından gelir.
- Admin panelden eklenen ürünler kullanıcı tarafında görünür.
- Sepete ekleme ve sepetten çıkarma işlemlerinde bildirim çıkar.
- Stokta olmayan ürün sepete eklenmez.
- Satın alma tamamlanınca ürün stoku düşer.
- Müşteri kayıt, giriş, adres ve sipariş geçmişi vardır.
- SQL yedeği başka bilgisayarda aynı verilerle projeyi ayağa kaldırır.

## SQL Dosyaları

SQL dosyaları `database` klasöründedir.

```text
database/schema.sql
database/stok_takip_backup.sql
database/README.md
```

`schema.sql` sadece tablo yapısını içerir.

`stok_takip_backup.sql` tablo yapısını ve mevcut verileri içerir. Arkadaşın projeyi kendi bilgisayarında tam haliyle açacaksa bu dosyayı import etmelidir.

## Antigravity / AI Agent İçin Kurulum Talimatı

Aşağıdaki metin direkt Antigravity veya başka bir yapay zeka kod editörüne verilebilir:

```text
Bu Go + PostgreSQL stok takip projesini Windows bilgisayarda çalıştır.

Repo:
https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi

Yapılacaklar:

1. Repoyu indir:
   git clone https://github.com/Orkunfbdev/Go-Programlama-Stok-Takip-Projesi.git
   cd Go-Programlama-Stok-Takip-Projesi

2. PostgreSQL'in bin klasörünü bul.
   Örnek:
   C:\Program Files\PostgreSQL\18\bin

3. PostgreSQL içinde stok_takip adında veritabanı oluştur:
   & "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip

4. Eğer veritabanı zaten varsa ve temiz kurulum isteniyorsa:
   & "C:\Program Files\PostgreSQL\18\bin\dropdb.exe" -U postgres --if-exists stok_takip
   & "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip

5. Tam SQL yedeğini import et:
   & "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"

6. PostgreSQL şifresi 5757 değilse uygulamayı çalıştırmadan önce:
   $env:DB_PASSWORD="POSTGRES_SIFRESI"

7. Go bağımlılıklarını indir:
   go mod download

8. Uygulamayı çalıştır:
   go run ./app

9. Tarayıcıda aç:
   http://localhost:8080

10. Admin paneli için:
    http://localhost:8080/login
    Kullanıcı adı: orkunorkun
    Şifre: 2005

Beklenen sonuç:
- Ana sayfa açılmalı.
- Sistemler sayfasında ürünler görünmeli.
- Kategori filtreleri çalışmalı.
- Sepete ekleme ve çıkarma çalışmalı.
- Admin panelinden ürün/kategori yönetimi çalışmalı.
```

## Kontrol Komutları

Kodun derlenip derlenmediğini kontrol etmek için:

```powershell
go test ./...
```

Veritabanında ürünlerin geldiğini kontrol etmek için:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -c "SELECT id, isim, stok, kategori FROM public.products ORDER BY id;"
```

Tablo sayısını kontrol etmek için:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -c "SELECT COUNT(*) FROM public.products;"
```

## Sık Karşılaşılan Hatalar

### `psql` veya `createdb` bulunamadı

PostgreSQL sürüm yolu farklı olabilir. Şu klasörü kontrol et:

```text
C:\Program Files\PostgreSQL
```

Komutlardaki `18` kısmını bilgisayarda kurulu olan sürümle değiştir.

### `password authentication failed`

PostgreSQL şifresi uygulamadaki varsayılan `5757` değildir.

Çözüm:

```powershell
$env:DB_PASSWORD="KENDI_POSTGRES_SIFREN"
go run ./app
```

### `database "stok_takip" does not exist`

Veritabanı oluşturulmamıştır.

Çözüm:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
```

### `relation "public.products" does not exist`

SQL yedeği import edilmemiştir.

Çözüm:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

### `listen tcp :8080: bind`

8080 portu başka bir uygulama tarafından kullanılıyordur. O uygulamayı kapat veya `app/main.go` içindeki portu değiştir:

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

## Notlar

- Ana Go kodu `app/main.go` dosyasındadır.
- HTML dosyaları `app/templates` klasöründedir.
- SQL dosyaları `database` klasöründedir.
- Derlenmiş `.exe`, log ve geçici dosyalar GitHub'a yüklenmez.
- Arkadaşının bilgisayarında tam verili kurulum için `database/stok_takip_backup.sql` kullanılmalıdır.
