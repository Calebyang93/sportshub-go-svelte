import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	API_URL: {
		description: 'Base URL of the Go API, e.g. http://localhost:8080',
		schema: (value) => (value && value.trim() ? value.trim().replace(/\/$/, '') : 'http://localhost:8080')
	}
});
