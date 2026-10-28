import { error } from '@sveltejs/kit';
import { api, ApiError, qs } from '#lib/server/api.ts';
import type { Facets, Product } from '#lib/types.ts';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
	const filters = {
		q: url.searchParams.get('q') ?? '',
		category: url.searchParams.get('category') ?? '',
		sport: url.searchParams.get('sport') ?? '',
		sort: url.searchParams.get('sort') ?? '',
		in_stock: url.searchParams.get('in_stock') ?? ''
	};
	try {
		const [products, facets] = await Promise.all([
			api.get<Product[]>(fetch, `/api/products${qs(filters)}`),
			api.get<Facets>(fetch, '/api/facets')
		]);
		return { products, facets, filters };
	} catch (e) {
		if (e instanceof ApiError) error(e.status, e.message);
		throw e;
	}
};
