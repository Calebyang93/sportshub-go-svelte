<script lang="ts">
	import ProductCard from '#lib/components/ProductCard.svelte';
	import { cart } from '#lib/cart.svelte.ts';
	import { formatPrice, sportGlyph } from '#lib/format.ts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const p = $derived(data.product);

	let qty = $state(1);
	let added = $state(false);
	const inCart = $derived(cart.lines.find((l) => l.id === p.id)?.qty ?? 0);
	const maxAddable = $derived(Math.max(0, p.stock - inCart));

	function add() {
		cart.add(p, qty);
		added = true;
		qty = 1;
		setTimeout(() => (added = false), 1500);
	}
</script>

<svelte:head><title>{p.name} · SportsHub</title></svelte:head>

<div class="container">
	<nav class="crumbs"><a href="/shop">Shop</a> / <a href="/shop?sport={p.sport}">{p.sport}</a> / {p.name}</nav>

	<div class="detail">
		<div class="visual">
			{#if p.image_url}
				<img src={p.image_url} alt={p.name} />
			{:else}
				<span aria-hidden="true">{sportGlyph(p.sport)}</span>
			{/if}
		</div>

		<div class="info">
			<p class="meta">{p.brand} · {p.category}</p>
			<h1>{p.name}</h1>
			<p class="price">{formatPrice(p.price_cents)}</p>
			<p class="desc">{p.description}</p>

			<p class="stock">
				{#if p.stock === 0}
					<span class="pill out">Sold out</span>
				{:else if p.stock <= 8}
					<span class="pill low">Only {p.stock} left</span>
				{:else}
					<span class="ok">● In stock</span>
				{/if}
				{#if inCart}<span class="in-cart">{inCart} in your cart</span>{/if}
			</p>

			{#if p.stock > 0}
				<div class="buy">
					<label class="field qty">
						<span class="visually-hidden">Quantity</span>
						<input type="number" min="1" max={maxAddable} bind:value={qty} disabled={maxAddable === 0} />
					</label>
					<button class="btn btn-accent" onclick={add} disabled={maxAddable === 0 || qty < 1 || qty > maxAddable}>
						{added ? 'Added to cart ✓' : maxAddable === 0 ? 'All stock in cart' : 'Add to cart'}
					</button>
					{#if inCart}<a class="btn btn-ghost" href="/cart">View cart</a>{/if}
				</div>
			{/if}
		</div>
	</div>

	{#if data.related.length}
		<h2 class="more">More for {p.sport}</h2>
		<div class="grid">
			{#each data.related as product (product.id)}
				<ProductCard {product} />
			{/each}
		</div>
	{/if}
</div>

<style>
	.crumbs {
		padding: 24px 0 16px;
		font-size: 0.85rem;
		color: var(--muted);
	}
	.crumbs a {
		text-decoration: none;
	}
	.detail {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
		gap: 36px;
		align-items: start;
	}
	.visual {
		aspect-ratio: 1;
		border-radius: var(--radius);
		background: linear-gradient(135deg, var(--court), var(--court-2));
		display: grid;
		place-items: center;
		font-size: 7rem;
		overflow: hidden;
	}
	.visual img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.meta {
		color: var(--muted);
		margin: 0 0 8px;
		font-size: 0.9rem;
	}
	h1 {
		font-size: clamp(2rem, 5vw, 3rem);
		margin: 0;
	}
	.price {
		font-family: var(--display);
		font-size: 2.2rem;
		font-weight: 800;
		margin: 14px 0 6px;
	}
	.desc {
		color: var(--muted);
		max-width: 52ch;
	}
	.stock {
		display: flex;
		align-items: center;
		gap: 12px;
		font-size: 0.88rem;
	}
	.ok {
		color: var(--court-2);
		font-weight: 600;
	}
	.in-cart {
		color: var(--muted);
	}
	.buy {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
		margin-top: 18px;
	}
	.qty input {
		width: 80px;
	}
	.more {
		margin: 56px 0 16px;
		font-size: 1.8rem;
	}
	@media (max-width: 720px) {
		.detail {
			grid-template-columns: 1fr;
			gap: 20px;
		}
		.visual {
			aspect-ratio: 16 / 10;
			font-size: 5rem;
		}
	}
</style>
