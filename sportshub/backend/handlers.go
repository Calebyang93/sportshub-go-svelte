package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

// ---------- models ----------

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Sport       string  `json:"sport"`
	Brand       string  `json:"brand"`
	Description string  `json:"description"`
	PriceCents  int64   `json:"price_cents"`
	Stock       int64   `json:"stock"`
	ImageURL    *string `json:"image_url"`
}

type Event struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Kind        string  `json:"kind"`
	Sport       string  `json:"sport"`
	Description string  `json:"description"`
	Coach       string  `json:"coach"`
	Venue       string  `json:"venue"`
	Level       string  `json:"level"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      string  `json:"ends_at"`
	Capacity    int64   `json:"capacity"`
	Registered  int64   `json:"registered"`
	SpotsLeft   int64   `json:"spots_left"`
	PriceCents  int64   `json:"price_cents"`
	ImageURL    *string `json:"image_url"`
}

type Server struct{ db *sql.DB }

// ---------- helpers ----------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func pathID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}

func decode(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// ---------- products ----------

const productCols = `id, name, category, sport, brand, description, price_cents, stock, image_url`

func scanProduct(s interface{ Scan(...any) error }) (Product, error) {
	var p Product
	err := s.Scan(&p.ID, &p.Name, &p.Category, &p.Sport, &p.Brand, &p.Description, &p.PriceCents, &p.Stock, &p.ImageURL)
	return p, err
}

// GET /api/products?q=&category=&sport=&in_stock=1&sort=price_asc|price_desc|name
func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	where := []string{"1=1"}
	args := []any{}
	if v := strings.TrimSpace(q.Get("q")); v != "" {
		where = append(where, "(name LIKE ? OR brand LIKE ? OR description LIKE ?)")
		like := "%" + v + "%"
		args = append(args, like, like, like)
	}
	if v := q.Get("category"); v != "" {
		where = append(where, "category = ?")
		args = append(args, v)
	}
	if v := q.Get("sport"); v != "" {
		where = append(where, "sport = ?")
		args = append(args, v)
	}
	if q.Get("in_stock") == "1" {
		where = append(where, "stock > 0")
	}
	order := "name ASC"
	switch q.Get("sort") {
	case "price_asc":
		order = "price_cents ASC"
	case "price_desc":
		order = "price_cents DESC"
	}

	rows, err := s.db.QueryContext(r.Context(),
		`SELECT `+productCols+` FROM products WHERE `+strings.Join(where, " AND ")+` ORDER BY `+order, args...)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		out = append(out, p)
	}
	writeJSON(w, 200, out)
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, 400, "invalid id")
		return
	}
	p, err := scanProduct(s.db.QueryRowContext(r.Context(), `SELECT `+productCols+` FROM products WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		writeErr(w, 404, "product not found")
		return
	} else if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, p)
}

// GET /api/facets — distinct values for filter dropdowns.
func (s *Server) facets(w http.ResponseWriter, r *http.Request) {
	get := func(query string) ([]string, error) {
		rows, err := s.db.QueryContext(r.Context(), query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		out := []string{}
		for rows.Next() {
			var v string
			rows.Scan(&v)
			out = append(out, v)
		}
		return out, rows.Err()
	}
	cats, err1 := get(`SELECT DISTINCT category FROM products ORDER BY 1`)
	pSports, err2 := get(`SELECT DISTINCT sport FROM products ORDER BY 1`)
	eSports, err3 := get(`SELECT DISTINCT sport FROM events ORDER BY 1`)
	if err := errors.Join(err1, err2, err3); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{
		"product_categories": cats,
		"product_sports":     pSports,
		"event_sports":       eSports,
		"event_kinds":        []string{"training", "coaching", "event"},
	})
}

// ---------- events ----------

const eventSelect = `
SELECT e.id, e.title, e.kind, e.sport, e.description, e.coach, e.venue, e.level,
       e.starts_at, e.ends_at, e.capacity, e.price_cents, e.image_url,
       (SELECT COUNT(*) FROM registrations r WHERE r.event_id = e.id) AS registered
FROM events e`

func scanEvent(s interface{ Scan(...any) error }) (Event, error) {
	var e Event
	err := s.Scan(&e.ID, &e.Title, &e.Kind, &e.Sport, &e.Description, &e.Coach, &e.Venue, &e.Level,
		&e.StartsAt, &e.EndsAt, &e.Capacity, &e.PriceCents, &e.ImageURL, &e.Registered)
	e.SpotsLeft = max(e.Capacity-e.Registered, 0)
	return e, err
}

// GET /api/events?kind=&sport=&q=&include_past=1
func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	where := []string{"1=1"}
	args := []any{}
	if q.Get("include_past") != "1" {
		where = append(where, "e.ends_at >= ?")
		args = append(args, time.Now().UTC().Format(time.RFC3339))
	}
	if v := q.Get("kind"); v != "" {
		where = append(where, "e.kind = ?")
		args = append(args, v)
	}
	if v := q.Get("sport"); v != "" {
		where = append(where, "e.sport = ?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q.Get("q")); v != "" {
		where = append(where, "(e.title LIKE ? OR e.coach LIKE ? OR e.venue LIKE ?)")
		like := "%" + v + "%"
		args = append(args, like, like, like)
	}
	rows, err := s.db.QueryContext(r.Context(),
		eventSelect+` WHERE `+strings.Join(where, " AND ")+` ORDER BY e.starts_at ASC`, args...)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer rows.Close()
	out := []Event{}
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		out = append(out, e)
	}
	writeJSON(w, 200, out)
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, 400, "invalid id")
		return
	}
	e, err := scanEvent(s.db.QueryRowContext(r.Context(), eventSelect+` WHERE e.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		writeErr(w, 404, "event not found")
		return
	} else if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, e)
}

type registrationReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Notes string `json:"notes"`
}

// POST /api/events/{id}/registrations
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, 400, "invalid id")
		return
	}
	var req registrationReq
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, "invalid JSON body")
		return
	}
	req.Name, req.Email = strings.TrimSpace(req.Name), strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" {
		writeErr(w, 422, "name is required")
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		writeErr(w, 422, "a valid email is required")
		return
	}

	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer tx.Rollback()

	var capacity, registered int64
	var endsAt string
	err = tx.QueryRow(`SELECT capacity, ends_at, (SELECT COUNT(*) FROM registrations WHERE event_id = ?) FROM events WHERE id = ?`, id, id).
		Scan(&capacity, &endsAt, &registered)
	if errors.Is(err, sql.ErrNoRows) {
		writeErr(w, 404, "event not found")
		return
	} else if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	if end, _ := time.Parse(time.RFC3339, endsAt); end.Before(time.Now()) {
		writeErr(w, 409, "this session has already finished")
		return
	}
	if registered >= capacity {
		writeErr(w, 409, "this session is fully booked")
		return
	}
	res, err := tx.Exec(`INSERT INTO registrations (event_id, name, email, phone, notes) VALUES (?,?,?,?,?)`,
		id, req.Name, req.Email, strings.TrimSpace(req.Phone), strings.TrimSpace(req.Notes))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeErr(w, 409, "this email is already registered for this session")
			return
		}
		writeErr(w, 500, err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	regID, _ := res.LastInsertId()
	writeJSON(w, 201, map[string]any{"id": regID, "event_id": id, "spots_left": capacity - registered - 1})
}

// ---------- orders ----------

type orderReq struct {
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
	Items        []struct {
		ProductID int64 `json:"product_id"`
		Qty       int64 `json:"qty"`
	} `json:"items"`
}

// POST /api/orders — prices come from the DB, never the client; stock is decremented atomically.
func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var req orderReq
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, "invalid JSON body")
		return
	}
	req.CustomerName, req.Email = strings.TrimSpace(req.CustomerName), strings.ToLower(strings.TrimSpace(req.Email))
	if req.CustomerName == "" {
		writeErr(w, 422, "name is required")
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		writeErr(w, 422, "a valid email is required")
		return
	}
	if len(req.Items) == 0 {
		writeErr(w, 422, "cart is empty")
		return
	}

	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer tx.Rollback()

	type line struct{ pid, qty, price int64 }
	lines := []line{}
	var total int64
	for _, it := range req.Items {
		if it.Qty <= 0 {
			writeErr(w, 422, "quantities must be positive")
			return
		}
		var name string
		var price, stock int64
		err := tx.QueryRow(`SELECT name, price_cents, stock FROM products WHERE id = ?`, it.ProductID).Scan(&name, &price, &stock)
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, 422, "product "+strconv.FormatInt(it.ProductID, 10)+" does not exist")
			return
		} else if err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		if stock < it.Qty {
			writeErr(w, 409, "only "+strconv.FormatInt(stock, 10)+" of "+name+" left in stock")
			return
		}
		lines = append(lines, line{it.ProductID, it.Qty, price})
		total += price * it.Qty
	}

	res, err := tx.Exec(`INSERT INTO orders (customer_name, email, total_cents) VALUES (?,?,?)`, req.CustomerName, req.Email, total)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	orderID, _ := res.LastInsertId()
	for _, l := range lines {
		if _, err := tx.Exec(`INSERT INTO order_items (order_id, product_id, qty, unit_price_cents) VALUES (?,?,?,?)`,
			orderID, l.pid, l.qty, l.price); err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		if _, err := tx.Exec(`UPDATE products SET stock = stock - ? WHERE id = ?`, l.qty, l.pid); err != nil {
			writeErr(w, 500, err.Error())
			return
		}
	}
	if err := tx.Commit(); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]any{"id": orderID, "total_cents": total})
}
