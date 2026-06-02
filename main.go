package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
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
type HomeData struct{ Featured []Product }

var db *sql.DB

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
	"add":   func(a, b int) int { return a + b },
}).ParseGlob("templates/*.html"))

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
	featured := list
	if len(list) > 4 {
		featured = list[:4]
	}
	tmpl.ExecuteTemplate(w, "index.html", HomeData{Featured: featured})
}

func systemsHandler(w http.ResponseWriter, r *http.Request) {
	list, err := getProducts()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "systems.html", list)
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

	tx, err := db.Begin()
	if err != nil {
		http.Redirect(w, r, "/cart?error=db", http.StatusFound)
		return
	}
	defer tx.Rollback()

	for id, qty := range cart {
		var stock int
		err := tx.QueryRow("SELECT stok FROM public.products WHERE id = $1 FOR UPDATE", id).Scan(&stock)
		if err != nil {
			http.Redirect(w, r, "/cart?error=db", http.StatusFound)
			return
		}
		if stock < qty {
			http.Redirect(w, r, "/cart?error=stock", http.StatusFound)
			return
		}
		_, err = tx.Exec("UPDATE public.products SET stok = stok - $1 WHERE id = $2", qty, id)
		if err != nil {
			http.Redirect(w, r, "/cart?error=db", http.StatusFound)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Redirect(w, r, "/cart?error=db", http.StatusFound)
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
	list, err := getProducts()
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "admin.html", list)
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
		_, err := db.Exec("UPDATE public.products SET stok = $1, fiyat = $2 WHERE id = $3", stock, price, id)
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

func ensureDatabaseSchema() error {
	if err := renameTrimmedTable("products"); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS public.products (
			id SERIAL PRIMARY KEY,
			isim TEXT NOT NULL,
			fiyat NUMERIC(10, 2) NOT NULL DEFAULT 0,
			stok INTEGER NOT NULL DEFAULT 0,
			kategori TEXT NOT NULL DEFAULT '',
			resim TEXT DEFAULT '',
			"tanım" TEXT DEFAULT ''
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS public.admin (
			id SERIAL PRIMARY KEY,
			"kullanıcıadı" TEXT NOT NULL UNIQUE,
			"şifre" TEXT NOT NULL
		)
	`)
	return err
}

func renameTrimmedTable(tableName string) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)
	`, tableName).Scan(&exists)
	if err != nil || exists {
		return err
	}

	var currentName string
	err = db.QueryRow(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
			AND btrim(table_name) = $1
			AND table_name <> $1
		ORDER BY length(table_name)
		LIMIT 1
	`, tableName).Scan(&currentName)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = db.Exec(`ALTER TABLE public.` + pq.QuoteIdentifier(currentName) + ` RENAME TO ` + pq.QuoteIdentifier(tableName))
	return err
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
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Stok takip sunucusu: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
