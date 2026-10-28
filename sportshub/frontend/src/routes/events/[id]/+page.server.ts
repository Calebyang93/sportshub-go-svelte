import { error, fail } from '@sveltejs/kit';
import { api, ApiError } from '#lib/server/api.ts';
import type { SportsEvent } from '#lib/types.ts';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
	try {
		const event = await api.get<SportsEvent>(fetch, `/api/events/${encodeURIComponent(params.id)}`);
		return { event };
	} catch (e) {
		if (e instanceof ApiError) error(e.status === 400 ? 404 : e.status, e.message);
		throw e;
	}
};

export const actions: Actions = {
	register: async ({ fetch, params, request }) => {
		const form = await request.formData();
		const values = {
			name: String(form.get('name') ?? '').trim(),
			email: String(form.get('email') ?? '').trim(),
			phone: String(form.get('phone') ?? '').trim(),
			notes: String(form.get('notes') ?? '').trim()
		};

		if (!values.name) return fail(422, { values, message: 'Please enter your name.' });
		if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.email))
			return fail(422, { values, message: 'Please enter a valid email address.' });

		try {
			const res = await api.post<{ id: number; spots_left: number }>(
				fetch,
				`/api/events/${encodeURIComponent(params.id)}/registrations`,
				values
			);
			return { success: true as const, name: values.name, email: values.email, registrationId: res.id };
		} catch (e) {
			if (e instanceof ApiError) return fail(e.status, { values, message: e.message });
			throw e;
		}
	}
};
