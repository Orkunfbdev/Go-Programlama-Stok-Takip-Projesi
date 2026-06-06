# Database

Bu klasor PostgreSQL dosyalarini tutar.

- `schema.sql`: Sadece tablo yapisi.
- `stok_takip_backup.sql`: Tablo yapisi ve mevcut veriler.

Yeni bilgisayarda tam yedegi yuklemek icin:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d stok_takip -f "database\stok_takip_backup.sql"
```
