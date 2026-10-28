# SportsHub

A SvelteKit storefront for **sports equipment** plus a booking site for **training, coaching sessions and events**, backed by a **Go** JSON API over **SQLite**.

```
sportshub/
├── backend/            Go API (net/http, database/sql, SQLite)
│   ├── main.go         server, routes, CORS, request logging
│   ├── db.go           schema (products, events, registrations, orders, order_items) + demo seed
│   └── handlers.go     JSON handlers
└── frontend/           SvelteKit 3 / Svelte 5 app (adapter-node)
    └── src/
        ├── env.ts                      API_URL definition
        ├── lib/server/api.ts           server-only client for the Go API
        ├── lib/cart.svelte.ts          cart store (runes + localStorage)
        ├── lib/components/             ProductCard, EventCard
        └── routes/
            ├── +page                   home: upcoming sessions + featured gear
            ├── shop/ , shop/[id]       catalogue with search/filter/sort, product detail
            ├── events/ , events/[id]   schedule grouped by day, detail + registration form
            └── cart/                   cart + checkout form action
```

## Run it

**1. API** (Go 1.22+; the `mattn/go-sqlite3` driver needs CGO and a C compiler)

```bash
cd backend
go run .                 # http://localhost:8080, creates + seeds sportshub.db on first start
```

| Env var        | Default                  | Purpose                                  |
| -------------- | ------------------------ | ---------------------------------------- |
| `ADDR`         | `:8080`                  | listen address                           |
| `DB_PATH`      | `sportshub.db`           | SQLite file                              |
| `CORS_ORIGINS` | `http://localhost:5173`  | only needed if a browser calls the API directly |

**2. Frontend** (Node 22+)

```bash
cd frontend
npm install
npm run dev              # http://localhost:5173
```

Set `API_URL` if the API isn't on `http://localhost:8080` (e.g. `API_URL=http://api:8080 npm run dev`).

Production: `npm run build && API_URL=... PORT=3000 node build`. adapter-node assumes **HTTPS** when working out the site's origin, which it uses to reject cross-site form posts, so run it behind a TLS proxy and pass `PROTOCOL_HEADER=x-forwarded-proto` (and `HOST_HEADER=x-forwarded-host` if the proxy rewrites the host).

## How it fits together

- All data loading happens in `+page.server.ts` files, so the browser only talks to SvelteKit; the Go API URL stays private and CORS isn't needed.
- Registration and checkout are SvelteKit **form actions** with `use:enhance`, so they also work with JavaScript disabled.
- The Go API owns the rules: it re-prices every order line from the `products` table, decrements stock in the same transaction, and enforces session capacity and one registration per email per session. The frontend just shows its error messages.

## API

| Method | Path                                | Notes |
| ------ | ----------------------------------- | ----- |
| GET    | `/api/health`                       | |
| GET    | `/api/facets`                       | distinct categories / sports for filter dropdowns |
| GET    | `/api/products`                     | `q`, `category`, `sport`, `in_stock=1`, `sort=price_asc\|price_desc` |
| GET    | `/api/products/{id}`                | |
| GET    | `/api/events`                       | upcoming only; `kind=training\|coaching\|event`, `sport`, `q`, `include_past=1` |
| GET    | `/api/events/{id}`                  | includes `registered` and `spots_left` |
| POST   | `/api/events/{id}/registrations`    | `{name, email, phone?, notes?}` → 201, or 409 when full / already registered |
| POST   | `/api/orders`                       | `{customer_name, email, items:[{product_id, qty}]}` → 201 `{id, total_cents}`, 409 on low stock |

Prices are stored in cents (SGD). Event times are stored as RFC 3339 UTC and displayed in Asia/Singapore time.

## Using your own tables

Point `DB_PATH` at your database. If your column names differ, change the `SELECT` lists in `handlers.go` (`productCols`, `eventSelect`) — the JSON shape the frontend expects is defined in `frontend/src/lib/types.ts`. Seeding only runs when `products` / `events` are empty.

## Not included yet

Payments (sessions say "pay at the venue"; orders are recorded without charging), auth/admin screens, emails, and product images (`image_url` is supported; a sport icon shows when it's empty).
