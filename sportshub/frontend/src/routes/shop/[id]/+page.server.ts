import { error } from '@sveltejs/kit';
import { api, ApiError } from '#lib/server/api.ts';
import type { Product } from '#lib/types.ts';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
	try {
		const product = await api.get<Product>(fetch, `/api/products/${encodeURIComponent(params.id)}`);
		const related = (await api.get<Product[]>(fetch, `/api/products?sport=${encodeURIComponent(product.sport)}`))
			.filter((p) => p.id !== product.id)
			.slice(0, 4);
		return { product, related };
	} catch (e) {
		if (e instanceof ApiError) error(e.status === 400 ? 404 : e.status, e.message);
		throw e;
	}
};
