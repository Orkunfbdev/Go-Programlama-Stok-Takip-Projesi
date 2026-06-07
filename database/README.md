# Database Klasörü

Bu klasör PostgreSQL kurulum dosyalarını tutar.

## Dosyalar

- `schema.sql`: Sadece tablo yapısını kurar, örnek veri eklemez.
- `stok_takip_backup.sql`: Tablo yapısını ve mevcut proje verilerini birlikte kurar.

Arkadaşının bilgisayarında tam çalışan proje verileriyle kurulum yapmak için `stok_takip_backup.sql` kullanılmalıdır.

## Temiz Kurulum

PowerShell:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\dropdb.exe" -U postgres --if-exists stok_takip
& "C:\Program Files\PostgreSQL\18\bin\createdb.exe" -U postgres stok_takip
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```

PostgreSQL sürümü farklıysa komuttaki `18` kısmını kendi sürümüne göre değiştir.
