<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { cart } from '#lib/cart.svelte.ts';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	onMount(() => cart.load());

	const links = [
		{ href: '/shop', label: 'Shop' },
		{ href: '/events', label: 'Train & Events' }
	];
	const active = (href: string) => page.url.pathname.startsWith(href);
</script>

<svelte:head>
	<link rel="icon" href="/favicon.svg" />
	<title>SportsHub</title>
</svelte:head>

<header>
	<div class="container bar">
		<a href="/" class="logo"><span class="dot"></span>SportsHub</a>
		<nav>
			{#each links as l}
				<a href={l.href} class:active={active(l.href)}>{l.label}</a>
			{/each}
			<a href="/cart" class="cart" class:active={active('/cart')} aria-label="Cart, {cart.count} items">
				<svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"
					><path
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M3 4h2l2.4 11.2a2 2 0 0 0 2 1.6h7.7a2 2 0 0 0 2-1.5L21 8H6.2M10 21h.01M17 21h.01"
					/></svg
				>
				{#if cart.count > 0}<span class="badge">{cart.count}</span>{/if}
			</a>
		</nav>
	</div>
</header>

<main>
	{@render children()}
</main>

<footer>
	<div class="container">SportsHub · Gear, coaching and events in Singapore</div>
</footer>

<style>
	header {
		position: sticky;
		top: 0;
		z-index: 10;
		background: var(--court);
		color: #fff;
	}
	.bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		height: 60px;
	}
	.logo {
		display: flex;
		align-items: center;
		gap: 8px;
		font-family: var(--display);
		font-weight: 800;
		font-size: 1.5rem;
		text-transform: uppercase;
		color: #fff;
		text-decoration: none;
		letter-spacing: 0.02em;
	}
	.dot {
		width: 14px;
		height: 14px;
		border-radius: 50%;
		background: var(--lime);
	}
	nav {
		display: flex;
		align-items: center;
		gap: 4px;
	}
	nav a {
		color: rgb(255 255 255 / 0.8);
		text-decoration: none;
		font-weight: 500;
		font-size: 0.92rem;
		padding: 8px 12px;
		border-radius: 999px;
		white-space: nowrap;
	}
	nav a:hover {
		color: #fff;
		background: rgb(255 255 255 / 0.08);
	}
	nav a.active {
		color: var(--court);
		background: var(--lime);
	}
	.cart {
		position: relative;
		display: grid;
		place-items: center;
	}
	.badge {
		position: absolute;
		top: 0;
		right: 0;
		min-width: 18px;
		height: 18px;
		padding: 0 4px;
		border-radius: 9px;
		background: var(--clay);
		color: #fff;
		font-size: 0.7rem;
		font-weight: 700;
		display: grid;
		place-items: center;
	}
	main {
		min-height: calc(100vh - 140px);
		padding-bottom: 60px;
	}
	footer {
		border-top: 1px solid var(--line);
		padding: 24px 0;
		color: var(--muted);
		font-size: 0.85rem;
	}
	@media (max-width: 480px) {
		.logo {
			font-size: 1.25rem;
		}
		nav a {
			padding: 8px 9px;
			font-size: 0.85rem;
		}
	}
</style>
