<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head><title>{page.status} · SportsHub</title></svelte:head>

<div class="container wrap">
	<p class="code">{page.status}</p>
	<h1>{page.status === 404 ? 'Out of bounds' : 'Something went wrong'}</h1>
	<p class="msg">{page.error?.message}</p>
	<a class="btn" href="/">Back to home</a>
</div>

<style>
	.wrap {
		padding: 80px 16px;
		text-align: center;
	}
	.code {
		font-family: var(--display);
		font-weight: 800;
		font-size: 5rem;
		margin: 0;
		color: var(--court-2);
		line-height: 1;
	}
	h1 {
		font-size: 2.2rem;
		margin: 8px 0;
	}
	.msg {
		color: var(--muted);
		margin-bottom: 24px;
	}
</style>
