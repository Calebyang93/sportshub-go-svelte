export interface Product {
	id: number;
	name: string;
	category: string;
	sport: string;
	brand: string;
	description: string;
	price_cents: number;
	stock: number;
	image_url: string | null;
}

export type EventKind = 'training' | 'coaching' | 'event';

export interface SportsEvent {
	id: number;
	title: string;
	kind: EventKind;
	sport: string;
	description: string;
	coach: string;
	venue: string;
	level: string;
	starts_at: string;
	ends_at: string;
	capacity: number;
	registered: number;
	spots_left: number;
	price_cents: number;
	image_url: string | null;
}

export interface Facets {
	product_categories: string[];
	product_sports: string[];
	event_sports: string[];
	event_kinds: EventKind[];
}
