import { fail } from '@sveltejs/kit';
import { api, ApiError } from '#lib/server/api.ts';
import type { Actions } from './$types';

export const actions: Actions = {
	checkout: async ({ fetch, request }) => {
		const form = await request.formData();
		const values = {
			customer_name: String(form.get('customer_name') ?? '').trim(),
			email: String(form.get('email') ?? '').trim()
		};

		let items: { product_id: number; qty: number }[] = [];
		try {
			items = JSON.parse(String(form.get('items') ?? '[]'));
		} catch {
			return fail(400, { values, message: 'Your cart could not be read. Please refresh and try again.' });
		}
		if (!items.length) return fail(422, { values, message: 'Your cart is empty.' });
		if (!values.customer_name) return fail(422, { values, message: 'Please enter your name.' });
		if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.email))
			return fail(422, { values, message: 'Please enter a valid email address.' });

		try {
			const order = await api.post<{ id: number; total_cents: number }>(fetch, '/api/orders', { ...values, items });
			return { success: true as const, orderId: order.id, total: order.total_cents, email: values.email };
		} catch (e) {
			if (e instanceof ApiError) return fail(e.status, { values, message: e.message });
			throw e;
		}
	}
};
