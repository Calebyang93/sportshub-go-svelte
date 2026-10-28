import { error } from '@sveltejs/kit';
import { api, ApiError } from '#lib/server/api.ts';
import type { Product, SportsEvent } from '#lib/types.ts';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const [events, products] = await Promise.all([
			api.get<SportsEvent[]>(fetch, '/api/events'),
			api.get<Product[]>(fetch, '/api/products?in_stock=1')
		]);
		return {
			upcoming: events.slice(0, 4),
			featured: products.slice(0, 4),
			counts: { events: events.length, products: products.length }
		};
	} catch (e) {
		if (e instanceof ApiError) error(e.status, e.message);
		throw e;
	}
};
