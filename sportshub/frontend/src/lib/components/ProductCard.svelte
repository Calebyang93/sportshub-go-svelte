<script lang="ts">
	import type { Product } from '#lib/types.ts';
	import { formatPrice, sportGlyph } from '#lib/format.ts';
	import { cart } from '#lib/cart.svelte.ts';

	let { product }: { product: Product } = $props();
	let added = $state(false);

	function add() {
		cart.add(product);
		added = true;
		setTimeout(() => (added = false), 1200);
	}
</script>

<article class="card">
	<a href="/shop/{product.id}" class="visual" aria-label={product.name}>
		{#if product.image_url}
			<img src={product.image_url} alt="" />
		{:else}
			<span class="glyph" aria-hidden="true">{sportGlyph(product.sport)}</span>
			<span class="cat">{product.category}</span>
		{/if}
	</a>
	<div class="body">
		<p class="meta">{product.brand} · {product.sport}</p>
		<h3><a href="/shop/{product.id}">{product.name}</a></h3>
		<div class="row">
			<span class="price">{formatPrice(product.price_cents)}</span>
			{#if product.stock === 0}
				<span class="pill out">Sold out</span>
			{:else if product.stock <= 8}
				<span class="pill low">{product.stock} left</span>
			{/if}
		</div>
		<button class="btn btn-block" onclick={add} disabled={product.stock === 0}>
			{added ? 'Added ✓' : product.stock === 0 ? 'Out of stock' : 'Add to cart'}
		</button>
	</div>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
		overflow: hidden;
		transition: transform 0.15s ease, box-shadow 0.15s ease;
	}
	.card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow);
	}
	.visual {
		position: relative;
		display: grid;
		place-items: center;
		aspect-ratio: 4 / 3;
		background: linear-gradient(135deg, var(--court) 0%, var(--court-2) 100%);
		text-decoration: none;
	}
	.visual img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.glyph {
		font-size: 3.4rem;
		filter: drop-shadow(0 6px 10px rgb(0 0 0 / 0.25));
	}
	.cat {
		position: absolute;
		left: 12px;
		bottom: 10px;
		font-family: var(--display);
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		font-size: 0.8rem;
		color: var(--lime);
	}
	.body {
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 14px 16px 16px;
		flex: 1;
	}
	.meta {
		margin: 0;
		font-size: 0.78rem;
		color: var(--muted);
	}
	h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		line-height: 1.3;
		flex: 1;
	}
	h3 a {
		color: inherit;
		text-decoration: none;
	}
	h3 a:hover {
		text-decoration: underline;
	}
	.row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin: 4px 0 8px;
	}
	.price {
		font-family: var(--display);
		font-size: 1.35rem;
		font-weight: 700;
	}
</style>
