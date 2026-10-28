// Server-only client for the Go API. Pages call this from +page.server.ts, so the
// browser never talks to the Go service directly (no CORS, API_URL stays private).
import { API_URL } from '$app/env/private';

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

type Fetch = typeof fetch;

async function request<T>(fetch: Fetch, path: string, init?: RequestInit): Promise<T> {
	let res: Response;
	try {
		res = await fetch(`${API_URL}${path}`, {
			...init,
			headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) }
		});
	} catch {
		throw new ApiError(503, `Can't reach the API at ${API_URL}. Is the Go server running?`);
	}
	const body = await res.json().catch(() => ({}));
	if (!res.ok) throw new ApiError(res.status, body.error ?? res.statusText);
	return body as T;
}

export const api = {
	get: <T>(fetch: Fetch, path: string) => request<T>(fetch, path),
	post: <T>(fetch: Fetch, path: string, data: unknown) =>
		request<T>(fetch, path, { method: 'POST', body: JSON.stringify(data) })
};

/** Build a query string from only the non-empty params. */
export function qs(params: Record<string, string | null | undefined>) {
	const s = new URLSearchParams();
	for (const [k, v] of Object.entries(params)) if (v) s.set(k, v);
	const out = s.toString();
	return out ? `?${out}` : '';
}
