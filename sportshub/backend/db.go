package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS products (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    category     TEXT    NOT NULL,          -- e.g. Footwear, Rackets, Balls
    sport        TEXT    NOT NULL,          -- e.g. Running, Tennis
    brand        TEXT    NOT NULL DEFAULT '',
    description  TEXT    NOT NULL DEFAULT '',
    price_cents  INTEGER NOT NULL CHECK (price_cents >= 0),
    stock        INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    image_url    TEXT,
    created_at   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now'))
);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_sport    ON products(sport);

CREATE TABLE IF NOT EXISTS events (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        TEXT    NOT NULL,
    kind         TEXT    NOT NULL CHECK (kind IN ('training','coaching','event')),
    sport        TEXT    NOT NULL,
    description  TEXT    NOT NULL DEFAULT '',
    coach        TEXT    NOT NULL DEFAULT '',
    venue        TEXT    NOT NULL DEFAULT '',
    level        TEXT    NOT NULL DEFAULT 'All levels',
    starts_at    TEXT    NOT NULL,          -- RFC3339 UTC
    ends_at      TEXT    NOT NULL,
    capacity     INTEGER NOT NULL CHECK (capacity > 0),
    price_cents  INTEGER NOT NULL DEFAULT 0 CHECK (price_cents >= 0),
    image_url    TEXT
);
CREATE INDEX IF NOT EXISTS idx_events_starts ON events(starts_at);
CREATE INDEX IF NOT EXISTS idx_events_kind   ON events(kind);

CREATE TABLE IF NOT EXISTS registrations (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id    INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name        TEXT    NOT NULL,
    email       TEXT    NOT NULL,
    phone       TEXT    NOT NULL DEFAULT '',
    notes       TEXT    NOT NULL DEFAULT '',
    created_at  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
    UNIQUE (event_id, email)
);
CREATE INDEX IF NOT EXISTS idx_registrations_event ON registrations(event_id);

CREATE TABLE IF NOT EXISTS orders (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_name  TEXT    NOT NULL,
    email          TEXT    NOT NULL,
    total_cents    INTEGER NOT NULL,
    created_at     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now'))
);

CREATE TABLE IF NOT EXISTS order_items (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id          INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id        INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    qty               INTEGER NOT NULL CHECK (qty > 0),
    unit_price_cents  INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
`

func openDB(path string) (*sql.DB, error) {
	// _foreign_keys and _busy_timeout apply to every pooled connection.
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite: one writer; keeps capacity/stock checks race-free
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return db, nil
}

// seed inserts demo data only when the tables are empty.
func seed(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		products := []struct {
			name, cat, sport, brand, desc string
			price, stock                  int
		}{
			{"Velocity Road Runner 3", "Footwear", "Running", "Stride", "Neutral daily trainer with a responsive foam midsole and breathable knit upper.", 15900, 24},
			{"Trail Grip GTX", "Footwear", "Running", "Ridgeback", "Waterproof trail shoe with 5mm lugs for wet MacRitchie-style trails.", 18900, 12},
			{"Pro Tour 98 Racket", "Rackets", "Tennis", "Apex", "305g control frame, 16x19 pattern, strung with multifilament.", 26900, 8},
			{"Championship Tennis Balls (4-pack)", "Balls", "Tennis", "Apex", "ITF-approved pressurised balls for hard courts.", 1290, 120},
			{"Featherlight Badminton Racket", "Rackets", "Badminton", "Shuttle Co", "Head-light 78g frame for fast doubles play.", 13900, 15},
			{"Feather Shuttlecocks (tube of 12)", "Balls", "Badminton", "Shuttle Co", "Tournament-grade goose feather shuttles, speed 77.", 3490, 60},
			{"Match Football Size 5", "Balls", "Football", "Kickline", "Thermally bonded, FIFA Quality match ball.", 5990, 40},
			{"Shin Guards with Ankle Sock", "Protective", "Football", "Kickline", "Lightweight shell with EVA backing, sizes S–L.", 2490, 35},
			{"Indoor Basketball Size 7", "Balls", "Basketball", "Courtside", "Composite leather with deep channels for grip.", 4590, 30},
			{"Adjustable Dumbbell 24kg", "Strength", "Fitness", "IronWorks", "Dial-select 2.5–24kg, replaces 15 pairs of dumbbells.", 34900, 6},
			{"Resistance Band Set", "Strength", "Fitness", "IronWorks", "Five latex bands, 5–50kg, with door anchor and handles.", 3990, 50},
			{"Yoga Mat 6mm", "Mats", "Fitness", "Flowform", "Non-slip natural rubber mat with alignment lines.", 6900, 25},
			{"Swim Goggles Anti-Fog", "Swim", "Swimming", "Aqualine", "UV-protected mirrored lenses, three nose bridges included.", 2990, 45},
			{"GPS Running Watch", "Wearables", "Running", "Pacer", "Heart rate, GPS, interval workouts, 14-day battery.", 32900, 0},
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT INTO products (name, category, sport, brand, description, price_cents, stock) VALUES (?,?,?,?,?,?,?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, p := range products {
			if _, err := stmt.Exec(p.name, p.cat, p.sport, p.brand, p.desc, p.price, p.stock); err != nil {
				tx.Rollback()
				return err
			}
		}
		stmt.Close()
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		// Dates are relative to today so the demo always has upcoming sessions.
		day := time.Now().UTC().Truncate(24 * time.Hour)
		at := func(days, hour, min int) time.Time {
			// hour given in Singapore time (UTC+8)
			return day.AddDate(0, 0, days).Add(time.Duration(hour-8)*time.Hour + time.Duration(min)*time.Minute)
		}
		events := []struct {
			title, kind, sport, desc, coach, venue, level string
			start                                         time.Time
			mins, capacity, price                         int
		}{
			{"Couch to 5K – Week 1", "training", "Running", "Run-walk intervals, form drills and a gentle progression plan for new runners.", "Coach Mei Lin", "East Coast Park, Area C", "Beginner", at(3, 7, 0), 75, 20, 2500},
			{"Track Speed Session", "training", "Running", "400m repeats and strides to sharpen your 10K pace.", "Coach Daniel Tan", "Kallang Practice Track", "Intermediate", at(5, 19, 0), 90, 16, 3000},
			{"1-on-1 Tennis Technique Clinic", "coaching", "Tennis", "Video-assisted stroke analysis with a personalised practice plan.", "Coach Arjun Rao", "Kallang Tennis Centre, Court 4", "All levels", at(4, 9, 0), 60, 1, 9500},
			{"Junior Badminton Squad", "coaching", "Badminton", "Footwork, clears and smash technique for ages 10–15.", "Coach Siti Rahman", "Bishan Sports Hall", "Beginner", at(6, 16, 0), 90, 12, 4000},
			{"Strength Foundations Bootcamp", "training", "Fitness", "Compound lifts, kettlebell basics and mobility in small groups.", "Coach Ryan Lee", "Studio 8, Bukit Timah", "All levels", at(2, 18, 30), 60, 10, 3500},
			{"Sunrise Yoga Flow", "training", "Fitness", "Vinyasa flow to open the hips and reset after a week of training.", "Coach Priya Nair", "Marina Barrage Rooftop", "All levels", at(7, 7, 0), 60, 25, 0},
			{"Community 10K Fun Run", "event", "Running", "Chip-timed 10K loop with finisher medal and post-race brunch.", "", "Gardens by the Bay East", "All levels", at(14, 6, 30), 180, 300, 3800},
			{"3v3 Street Basketball Tournament", "event", "Basketball", "Register as an individual; teams are balanced on the day.", "", "Our Tampines Hub", "Intermediate", at(10, 14, 0), 300, 48, 1500},
			{"Open Water Swim Clinic", "coaching", "Swimming", "Sighting, drafting and safety for your first open-water race.", "Coach Hannah Ong", "Sentosa Palawan Beach", "Intermediate", at(9, 7, 30), 90, 8, 6000},
			{"5-a-side Football Social", "event", "Football", "Friendly mixed games, bibs and balls provided.", "", "The Cage, Kallang", "All levels", at(1, 20, 0), 120, 20, 1200},
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT INTO events (title, kind, sport, description, coach, venue, level, starts_at, ends_at, capacity, price_cents) VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, e := range events {
			end := e.start.Add(time.Duration(e.mins) * time.Minute)
			if _, err := stmt.Exec(e.title, e.kind, e.sport, e.desc, e.coach, e.venue, e.level,
				e.start.Format(time.RFC3339), end.Format(time.RFC3339), e.capacity, e.price); err != nil {
				tx.Rollback()
				return err
			}
		}
		stmt.Close()
		// Pre-fill a few registrations so "spots left" is visibly non-trivial.
		tx.Exec(`INSERT INTO registrations (event_id, name, email) VALUES
			(1,'Alex Goh','alex@example.com'),(1,'Bea Lim','bea@example.com'),(1,'Chen Wei','chen@example.com'),
			(3,'Dana Koh','dana@example.com'),
			(5,'Eli Chua','eli@example.com'),(5,'Fiona Ho','fiona@example.com'),(5,'Gabe Ng','gabe@example.com'),
			(5,'Hui Min','huimin@example.com'),(5,'Ivan Teo','ivan@example.com'),(5,'Jas Yeo','jas@example.com'),
			(5,'Kai Lau','kai@example.com'),(5,'Lena Sim','lena@example.com')`)
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
