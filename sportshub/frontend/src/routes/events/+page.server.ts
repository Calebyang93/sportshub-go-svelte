import { error } from '@sveltejs/kit';
import { api, ApiError, qs } from '#lib/server/api.ts';
import type { Facets, SportsEvent } from '#lib/types.ts';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
	const filters = {
		kind: url.searchParams.get('kind') ?? '',
		sport: url.searchParams.get('sport') ?? '',
		q: url.searchParams.get('q') ?? ''
	};
	try {
		const [events, facets] = await Promise.all([
			api.get<SportsEvent[]>(fetch, `/api/events${qs(filters)}`),
			api.get<Facets>(fetch, '/api/facets')
		]);
		return { events, facets, filters };
	} catch (e) {
		if (e instanceof ApiError) error(e.status, e.message);
		throw e;
	}
};
