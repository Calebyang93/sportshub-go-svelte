// Client-side cart kept in localStorage. Prices shown here are for display only;
// the Go API re-prices every line from the products table at checkout.
import type { Product } from './types';

export interface CartLine {
	id: number;
	name: string;
	sport: string;
	price_cents: number;
	stock: number;
	qty: number;
}

const KEY = 'sportshub-cart';

class Cart {
	lines = $state<CartLine[]>([]);
	ready = $state(false);

	count = $derived(this.lines.reduce((n, l) => n + l.qty, 0));
	subtotal = $derived(this.lines.reduce((n, l) => n + l.qty * l.price_cents, 0));

	/** Call once in the browser (from the root layout). */
	load() {
		try {
			const raw = localStorage.getItem(KEY);
			if (raw) this.lines = JSON.parse(raw);
		} catch {
			/* storage blocked or corrupt — start empty */
		}
		this.ready = true;
	}

	private save() {
		try {
			localStorage.setItem(KEY, JSON.stringify(this.lines));
		} catch {
			/* ignore */
		}
	}

	add(p: Product, qty = 1) {
		const line = this.lines.find((l) => l.id === p.id);
		if (line) line.qty = Math.min(line.qty + qty, p.stock);
		else
			this.lines.push({ id: p.id, name: p.name, sport: p.sport, price_cents: p.price_cents, stock: p.stock, qty: Math.min(qty, p.stock) });
		this.save();
	}

	setQty(id: number, qty: number) {
		const line = this.lines.find((l) => l.id === id);
		if (!line) return;
		if (qty <= 0) this.remove(id);
		else {
			line.qty = Math.min(qty, line.stock);
			this.save();
		}
	}

	remove(id: number) {
		this.lines = this.lines.filter((l) => l.id !== id);
		this.save();
	}

	clear() {
		this.lines = [];
		this.save();
	}
}

export const cart = new Cart();
