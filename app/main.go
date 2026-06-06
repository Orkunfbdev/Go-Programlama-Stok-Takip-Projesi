package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	Stock    int
	Category string
	Image    string
	Desc     string
}
type Category struct {
	ID   int
	Name string
}
type CartItem struct {
	Product  Product
	Qty      int
	Subtotal float64
}
type CartPageData struct {
	Items   []CartItem
	Total   float64
	Success bool
	Error   string
}
type HomeData struct {
	Featured     []Product
	Categories   []Category
	MusteriAd    string
	MusteriGiris bool
}
type SystemsData struct {
	Products      []Product
	Categories    []Category
	AktifKategori string
}
type AdminMusteriRow struct {
	ID         int
	Ad         string
	Soyad      string
	Email      string
	Tel        string
	KayitTarih string
}
type AdminSiparisRow struct {
	ID         int
	MusteriAd  string
	Email      string
	Toplam     float64
	Durum      string
	AdresBilgi string
	Tarih      string
}
type AdminPageData struct {
	Products   []Product
	Categories []Category
	Musteriler []AdminMusteriRow
	Siparisler []AdminSiparisRow
}
type Musteri struct {
	ID    int
	Ad    string
	Soyad string
	Email string
	Tel   string
}
type Adres struct {
	ID         int
	MusteriID  int
	Baslik     string
	Il         string
	Ilce       string
	Mahalle    string
	AdresSatir string
	PostaKodu  string
	Varsayilan bool
}
type SiparisOzet struct {
	ID         int
	Toplam     float64
	Durum      string
	Tarih      string
	UrunOzeti  string
	AdresBilgi string
}
type CheckoutData struct {
	Musteri  Musteri
	Adresler []Adres
	Items    []CartItem
	Total    float64
	Error    string
}
type ProfilData struct {
	Musteri    Musteri
	Siparisler []SiparisOzet
}
type MusteriGirisPageData struct {
	Error    string
	Info     string
	Redirect string
}
type MusteriKayitPageData struct {
	Error string
}

func hashSifre(sifre string) string {
	h := sha256.Sum256([]byte("techshop_salt_" + sifre))
	return hex.EncodeToString(h[:])
}

func getMusteriSession(r *http.Request) (Musteri, bool) {
	c, err := r.Cookie("musteri_id")
	if err != nil || c.Value == "" {
		return Musteri{}, false
	}
	id, err := strconv.Atoi(c.Value)
	if err != nil {
		return Musteri{}, false
	}
	var m Musteri
	err = db.QueryRow(`SELECT id, ad, soyad, email, COALESCE(tel,'') FROM public.musteriler WHERE id=$1`, id).Scan(&m.ID, &m.Ad, &m.Soyad, &m.Email, &m.Tel)
	if err != nil {
		return Musteri{}, false
	}
	return m, true
}

func getMusteriAdresler(musteriID int) []Adres {
	rows, err := db.Query(`SELECT id, musteri_id, baslik, il, ilce, mahalle, adres_satir, COALESCE(posta_kodu,''), varsayilan FROM public.adresler WHERE musteri_id=$1 ORDER BY varsayilan DESC, id ASC`, musteriID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []Adres
	for rows.Next() {
		var a Adres
		rows.Scan(&a.ID, &a.MusteriID, &a.Baslik, &a.Il, &a.Ilce, &a.Mahalle, &a.AdresSatir, &a.PostaKodu, &a.Varsayilan)
		list = append(list, a)
	}
	return list
}

var db *sql.DB

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
	"add":   func(a, b int) int { return a + b },
}).ParseGlob(templatePattern()))

func templatePattern() string {
	candidates := []string{
		filepath.Join("templates", "*.html"),
		filepath.Join("app", "templates", "*.html"),
	}
	for _, pattern := range candidates {
		if matches, err := filepath.Glob(pattern); err == nil && len(matches) > 0 {
			return pattern
		}
	}
	return filepath.Join("templates", "*.html")
}

func isAdmin(r *http.Request) bool {
	c, err := r.Cookie("admin_session")
	return err == nil && c.Value == "logged_in"
}

func getCart(r *http.Request) map[int]int {
	cart := make(map[int]int)
	c, err := r.Cookie("cart")
	if err != nil || c.Value == "" {
		return cart
	}
	for _, item := range strings.Split(c.Value, ",") {
		p := strings.Split(item, ":")
		if len(p) == 2 {
			id, _ := strconv.Atoi(p[0])
			qty, _ := strconv.Atoi(p[1])
			if id > 0 && qty > 0 {
				cart[id] = qty
			}
		}
	}
	return cart
}

func saveCart(w http.ResponseWriter, cart map[int]int) {
	var parts []string
	for id, qty := range cart {
		if qty > 0 {
			parts = append(parts, fmt.Sprintf("%d:%d", id, qty))
		}
	}
	http.SetCookie(w, &http.Cookie{Name: "cart", Value: strings.Join(parts, ","), Path: "/"})
}

func getProducts() ([]Product, error) {
	rows, err := db.Query(`SELECT id, isim, fiyat, stok, kategori, COALESCE(resim, ''), COALESCE("tanım", '') FROM public.products ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category, &p.Image, &p.Desc); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func getCategories() ([]Category, error) {
	rows, err := db.Query(`SELECT id, isim FROM public.categories ORDER BY isim ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func getCartItems(r *http.Request) ([]CartItem, float64) {
	var items []CartItem
	var total float64
	list, err := getProducts()
	if err != nil {
		return nil, 0
	}
	for id, qty := range getCart(r) {
		for _, p := range list {
			if p.ID == id {
				sub := p.Price * float64(qty)
				items = append(items, CartItem{p, qty, sub})
				total += sub
				break
			}
		}
	}
	return items, total
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	list, err := getProducts()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cats, _ := getCategories()
	featured := list
	if len(list) > 4 {
		featured = list[:4]
	}
	data := HomeData{Featured: featured, Categories: cats}
	if m, ok := getMusteriSession(r); ok {
		data.MusteriAd = m.Ad + " " + m.Soyad
		data.MusteriGiris = true
	}
	tmpl.ExecuteTemplate(w, "index.html", data)
}

func systemsHandler(w http.ResponseWriter, r *http.Request) {
	list, err := getProducts()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cats, _ := getCategories()
	aktif := r.URL.Query().Get("kategori")
	if aktif != "" {
		var filtered []Product
		for _, p := range list {
			if p.Category == aktif {
				filtered = append(filtered, p)
			}
		}
		list = filtered
	}
	tmpl.ExecuteTemplate(w, "systems.html", SystemsData{Products: list, Categories: cats, AktifKategori: aktif})
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	items, total := getCartItems(r)
	tmpl.ExecuteTemplate(w, "cart.html", CartPageData{items, total, r.URL.Query().Get("success") == "1", r.URL.Query().Get("error")})
}

func addToCartHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	cart := getCart(r)

	var stock int
	err := db.QueryRow("SELECT stok FROM public.products WHERE id = $1", id).Scan(&stock)
	if err == nil && stock > 0 {
		if cart[id] < stock {
			cart[id]++
		}
	}
	saveCart(w, cart)
	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/systems"
	}
	http.Redirect(w, r, ref, http.StatusFound)
}

func removeFromCartHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	cart := getCart(r)
	delete(cart, id)
	saveCart(w, cart)
	http.Redirect(w, r, "/cart", http.StatusFound)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	cart := getCart(r)
	if len(cart) == 0 {
		http.Redirect(w, r, "/cart?error=empty", http.StatusFound)
		return
	}
	// Müşteri giriş yapmamışsa checkout'a yönlendir (oradan login/kayıt yaptırır)
	_, ok := getMusteriSession(r)
	if !ok {
		http.Redirect(w, r, "/musteri/giris?redirect=/checkout", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/checkout", http.StatusFound)
}

// ─── Müşteri Kayıt ──────────────────────────────────────────────────────────
func musteriKayitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{})
		return
	}
	ad := strings.TrimSpace(r.FormValue("ad"))
	soyad := strings.TrimSpace(r.FormValue("soyad"))
	email := strings.TrimSpace(r.FormValue("email"))
	tel := strings.TrimSpace(r.FormValue("tel"))
	sifre := r.FormValue("sifre")
	sifre2 := r.FormValue("sifre2")

	if ad == "" || soyad == "" || email == "" || sifre == "" {
		tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{Error: "Tüm alanları doldurun."})
		return
	}
	if sifre != sifre2 {
		tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{Error: "Şifreler eşleşmiyor!"})
		return
	}
	if len(sifre) < 6 {
		tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{Error: "Şifre en az 6 karakter olmalı."})
		return
	}
	var newID int
	err := db.QueryRow(`INSERT INTO public.musteriler (ad,soyad,email,sifre,tel) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		ad, soyad, email, hashSifre(sifre), tel).Scan(&newID)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{Error: "Bu e-posta zaten kayıtlı!"})
			return
		}
		log.Println("Kayıt hata:", err)
		tmpl.ExecuteTemplate(w, "musteri_kayit.html", MusteriKayitPageData{Error: "Kayıt sırasında hata oluştu."})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "musteri_id", Value: strconv.Itoa(newID), Path: "/", MaxAge: 86400 * 7})
	http.Redirect(w, r, "/checkout", http.StatusFound)
}

// ─── Müşteri Giriş ──────────────────────────────────────────────────────────
func musteriGirisHandler(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/musteri/profil"
	}
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "musteri_giris.html", MusteriGirisPageData{Redirect: redirect})
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	sifre := r.FormValue("sifre")
	redirectPost := r.FormValue("redirect")
	if redirectPost != "" {
		redirect = redirectPost
	}
	var id int
	err := db.QueryRow(`SELECT id FROM public.musteriler WHERE email=$1 AND sifre=$2`, email, hashSifre(sifre)).Scan(&id)
	if err != nil {
		tmpl.ExecuteTemplate(w, "musteri_giris.html", MusteriGirisPageData{Error: "E-posta veya şifre hatalı!", Redirect: redirect})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "musteri_id", Value: strconv.Itoa(id), Path: "/", MaxAge: 86400 * 7})
	http.Redirect(w, r, redirect, http.StatusFound)
}

// ─── Müşteri Çıkış ──────────────────────────────────────────────────────────
func musteriCikisHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "musteri_id", Value: "", MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/", http.StatusFound)
}

// ─── Müşteri Profil ─────────────────────────────────────────────────────────
func musteriProfilHandler(w http.ResponseWriter, r *http.Request) {
	m, ok := getMusteriSession(r)
	if !ok {
		http.Redirect(w, r, "/musteri/giris", http.StatusFound)
		return
	}
	rows, err := db.Query(`
		SELECT s.id, s.toplam_tutar, s.durum, s.olusturulma,
			COALESCE(a.il,'') || ' ' || COALESCE(a.ilce,'') as adres
		FROM public.siparisler s
		LEFT JOIN public.adresler a ON a.id = s.adres_id
		WHERE s.musteri_id=$1 ORDER BY s.olusturulma DESC`, m.ID)
	var siparisler []SiparisOzet
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s SiparisOzet
			var t time.Time
			rows.Scan(&s.ID, &s.Toplam, &s.Durum, &t, &s.AdresBilgi)
			s.Tarih = t.Format("02.01.2006 15:04")
			// Sipariş ürünlerini çek
			var ozet []string
			uRows, _ := db.Query(`SELECT urun_isim, adet FROM public.siparis_urunleri WHERE siparis_id=$1`, s.ID)
			if uRows != nil {
				for uRows.Next() {
					var isim string
					var adet int
					uRows.Scan(&isim, &adet)
					ozet = append(ozet, fmt.Sprintf("%s x%d", isim, adet))
				}
				uRows.Close()
			}
			s.UrunOzeti = strings.Join(ozet, ", ")
			siparisler = append(siparisler, s)
		}
	}
	tmpl.ExecuteTemplate(w, "profil.html", ProfilData{Musteri: m, Siparisler: siparisler})
}

// ─── Checkout ───────────────────────────────────────────────────────────────
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	m, ok := getMusteriSession(r)
	if !ok {
		http.Redirect(w, r, "/musteri/giris?redirect=/checkout", http.StatusFound)
		return
	}
	cart := getCart(r)
	if len(cart) == 0 {
		http.Redirect(w, r, "/cart", http.StatusFound)
		return
	}
	items, total := getCartItems(r)
	adresler := getMusteriAdresler(m.ID)

	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total})
		return
	}

	// POST: siparişi kaydet
	adresID := r.FormValue("adres_id")
	var finalAdresID int

	if adresID == "new" || adresID == "" {
		// Yeni adres kaydet
		bas := strings.TrimSpace(r.FormValue("baslik"))
		il := strings.TrimSpace(r.FormValue("il"))
		ilce := strings.TrimSpace(r.FormValue("ilce"))
		mah := strings.TrimSpace(r.FormValue("mahalle"))
		asat := strings.TrimSpace(r.FormValue("adres_satir"))
		pk := strings.TrimSpace(r.FormValue("posta_kodu"))
		vsyl := r.FormValue("varsayilan") == "1"
		if il == "" || ilce == "" || asat == "" {
			tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Lütfen il, ilçe ve açık adresi doldurun."})
			return
		}
		if bas == "" {
			bas = "Ev"
		}
		if vsyl {
			db.Exec(`UPDATE public.adresler SET varsayilan=FALSE WHERE musteri_id=$1`, m.ID)
		}
		err := db.QueryRow(`INSERT INTO public.adresler (musteri_id,baslik,il,ilce,mahalle,adres_satir,posta_kodu,varsayilan) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
			m.ID, bas, il, ilce, mah, asat, pk, vsyl).Scan(&finalAdresID)
		if err != nil {
			log.Println("Adres kayıt hata:", err)
			tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Adres kaydedilemedi."})
			return
		}
	} else {
		finalAdresID, _ = strconv.Atoi(adresID)
	}

	// Transaction ile stok düş + sipariş oluştur
	tx, err := db.Begin()
	if err != nil {
		tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Veritabanı hatası."})
		return
	}
	defer tx.Rollback()

	for id, qty := range cart {
		var stock int
		err := tx.QueryRow("SELECT stok FROM public.products WHERE id=$1 FOR UPDATE", id).Scan(&stock)
		if err != nil || stock < qty {
			tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Yetersiz stok! Sepetinizi kontrol edin."})
			return
		}
		tx.Exec("UPDATE public.products SET stok=stok-$1 WHERE id=$2", qty, id)
	}

	var sipID int
	err = tx.QueryRow(`INSERT INTO public.siparisler (musteri_id,adres_id,toplam_tutar,durum) VALUES ($1,$2,$3,'tamamlandi') RETURNING id`,
		m.ID, finalAdresID, total).Scan(&sipID)
	if err != nil {
		log.Println("Sipariş kayıt hata:", err)
		tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Sipariş oluşturulamadı."})
		return
	}

	for _, item := range items {
		tx.Exec(`INSERT INTO public.siparis_urunleri (siparis_id,urun_id,urun_isim,adet,birim_fiyat) VALUES ($1,$2,$3,$4,$5)`,
			sipID, item.Product.ID, item.Product.Name, item.Qty, item.Product.Price)
	}

	if err := tx.Commit(); err != nil {
		tmpl.ExecuteTemplate(w, "checkout.html", CheckoutData{Musteri: m, Adresler: adresler, Items: items, Total: total, Error: "Sipariş tamamlanamadı."})
		return
	}

	saveCart(w, make(map[int]int))
	http.Redirect(w, r, "/cart?success=1", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		err := db.QueryRow(`SELECT "şifre" FROM public.admin WHERE "kullanıcıadı" = $1`, username).Scan(&dbPassword)
		if err == nil && dbPassword == password {
			http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "logged_in", Path: "/"})
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		tmpl.ExecuteTemplate(w, "login.html", "Hatalı kullanıcı adı veya şifre!")
		return
	}
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	products, err := getProducts()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	categories, err := getCategories()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Müşteriler
	var musteriler []AdminMusteriRow
	mRows, _ := db.Query(`SELECT id, ad, soyad, email, COALESCE(tel,''), TO_CHAR(olusturulma,'DD.MM.YYYY HH24:MI') FROM public.musteriler ORDER BY id DESC`)
	if mRows != nil {
		defer mRows.Close()
		for mRows.Next() {
			var m AdminMusteriRow
			mRows.Scan(&m.ID, &m.Ad, &m.Soyad, &m.Email, &m.Tel, &m.KayitTarih)
			musteriler = append(musteriler, m)
		}
	}
	// Siparişler
	var siparisler []AdminSiparisRow
	sRows, _ := db.Query(`
		SELECT s.id, m.ad||' '||m.soyad, m.email, s.toplam_tutar, s.durum,
			COALESCE(a.il,'')||' '||COALESCE(a.ilce,''),
			TO_CHAR(s.olusturulma,'DD.MM.YYYY HH24:MI')
		FROM public.siparisler s
		LEFT JOIN public.musteriler m ON m.id=s.musteri_id
		LEFT JOIN public.adresler a ON a.id=s.adres_id
		ORDER BY s.olusturulma DESC LIMIT 50`)
	if sRows != nil {
		defer sRows.Close()
		for sRows.Next() {
			var s AdminSiparisRow
			sRows.Scan(&s.ID, &s.MusteriAd, &s.Email, &s.Toplam, &s.Durum, &s.AdresBilgi, &s.Tarih)
			siparisler = append(siparisler, s)
		}
	}
	tmpl.ExecuteTemplate(w, "admin.html", AdminPageData{Products: products, Categories: categories, Musteriler: musteriler, Siparisler: siparisler})
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		name := strings.TrimSpace(r.FormValue("name"))
		category := r.FormValue("category")
		_, err := db.Exec("UPDATE public.products SET isim = $1, stok = $2, fiyat = $3, kategori = $4 WHERE id = $5", name, stock, price, category, id)
		if err != nil {
			log.Println("Update error:", err)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		_, err := db.Exec(`INSERT INTO public.products (isim, fiyat, stok, kategori, resim, "tanım") VALUES ($1, $2, $3, $4, $5, $6)`,
			r.FormValue("name"), price, stock, r.FormValue("category"), r.FormValue("image"), r.FormValue("desc"))
		if err != nil {
			log.Println("Insert error:", err)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		_, err := db.Exec("DELETE FROM public.products WHERE id = $1", id)
		if err != nil {
			log.Println("Delete error:", err)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "", MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/", http.StatusFound)
}

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		name := strings.TrimSpace(r.FormValue("name"))
		if name != "" {
			_, err := db.Exec(`INSERT INTO public.categories (isim) VALUES ($1) ON CONFLICT (isim) DO NOTHING`, name)
			if err != nil {
				log.Println("Category insert error:", err)
			}
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		_, err := db.Exec("DELETE FROM public.categories WHERE id = $1", id)
		if err != nil {
			log.Println("Category delete error:", err)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func editCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := strings.TrimSpace(r.FormValue("name"))
		if name != "" {
			_, err := db.Exec("UPDATE public.categories SET isim = $1 WHERE id = $2", name, id)
			if err != nil {
				log.Println("Category update error:", err)
			}
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func ensureDatabaseSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS public.products (
			id SERIAL PRIMARY KEY, isim TEXT NOT NULL,
			fiyat NUMERIC(10,2) NOT NULL DEFAULT 0,
			stok INTEGER NOT NULL DEFAULT 0,
			kategori TEXT NOT NULL DEFAULT '',
			resim TEXT DEFAULT '', "tanım" TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS public.admin (
			id SERIAL PRIMARY KEY,
			"kullanıcıadı" TEXT NOT NULL UNIQUE,
			"şifre" TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS public.categories (
			id SERIAL PRIMARY KEY, isim TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS public.musteriler (
			id SERIAL PRIMARY KEY, ad TEXT NOT NULL, soyad TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE, sifre TEXT NOT NULL,
			tel TEXT DEFAULT '', olusturulma TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS public.adresler (
			id SERIAL PRIMARY KEY,
			musteri_id INTEGER NOT NULL REFERENCES public.musteriler(id) ON DELETE CASCADE,
			baslik TEXT NOT NULL DEFAULT 'Ev',
			il TEXT NOT NULL DEFAULT '', ilce TEXT NOT NULL DEFAULT '',
			mahalle TEXT NOT NULL DEFAULT '', adres_satir TEXT NOT NULL DEFAULT '',
			posta_kodu TEXT DEFAULT '', varsayilan BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS public.siparisler (
			id SERIAL PRIMARY KEY,
			musteri_id INTEGER REFERENCES public.musteriler(id) ON DELETE SET NULL,
			adres_id INTEGER REFERENCES public.adresler(id) ON DELETE SET NULL,
			toplam_tutar NUMERIC(12,2) NOT NULL,
			durum TEXT NOT NULL DEFAULT 'tamamlandi',
			olusturulma TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS public.siparis_urunleri (
			id SERIAL PRIMARY KEY,
			siparis_id INTEGER NOT NULL REFERENCES public.siparisler(id) ON DELETE CASCADE,
			urun_id INTEGER REFERENCES public.products(id) ON DELETE SET NULL,
			urun_isim TEXT NOT NULL, adet INTEGER NOT NULL, birim_fiyat NUMERIC(10,2) NOT NULL
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var err error
	connStr := "host=127.0.0.1 port=5432 user=postgres password=5757 dbname=stok_takip sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB Open Error:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("DB Ping Error: ", err)
	}

	if err := ensureDatabaseSchema(); err != nil {
		log.Fatal("DB Schema Error: ", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/systems", systemsHandler)
	http.HandleFunc("/cart", cartHandler)
	http.HandleFunc("/cart/add", addToCartHandler)
	http.HandleFunc("/cart/remove", removeFromCartHandler)
	http.HandleFunc("/cart/buy", buyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/update", updateHandler)
	http.HandleFunc("/admin/add", addProductHandler)
	http.HandleFunc("/admin/delete", deleteProductHandler)
	http.HandleFunc("/admin/categories/add", addCategoryHandler)
	http.HandleFunc("/admin/categories/delete", deleteCategoryHandler)
	http.HandleFunc("/admin/categories/edit", editCategoryHandler)
	http.HandleFunc("/logout", logoutHandler)
	// Müşteri
	http.HandleFunc("/musteri/kayit", musteriKayitHandler)
	http.HandleFunc("/musteri/giris", musteriGirisHandler)
	http.HandleFunc("/musteri/cikis", musteriCikisHandler)
	http.HandleFunc("/musteri/profil", musteriProfilHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Stok takip sunucusu: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
