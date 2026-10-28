<script lang="ts">
	import { goto } from '$app/navigation';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Submitting the GET form updates the URL, which re-runs the server load.
	// Selects auto-submit when JS is on; the button covers no-JS.
	function submit(e: Event) {
		const form = (e.currentTarget as HTMLElement).closest('form')!;
		const params = new URLSearchParams();
		for (const [k, v] of new FormData(form)) if (v) params.set(k, String(v));
		goto(`?${params}`, { reset: false, replace: true });
	}
</script>

<svelte:head><title>Shop · SportsHub</title></svelte:head>

<div class="container">
	<div class="page-head">
		<h1>Shop gear</h1>
		<p>{data.products.length} {data.products.length === 1 ? 'item' : 'items'}</p>
	</div>

	<form class="filters" method="GET" onsubmit={(e) => { e.preventDefault(); submit(e); }}>
		<label class="field grow">
			Search
			<input type="search" name="q" value={data.filters.q} placeholder="Rackets, shoes, brand…" />
		</label>
		<label class="field">
			Category
			<select name="category" value={data.filters.category} onchange={submit}>
				<option value="">All</option>
				{#each data.facets.product_categories as c}<option>{c}</option>{/each}
			</select>
		</label>
		<label class="field">
			Sport
			<select name="sport" value={data.filters.sport} onchange={submit}>
				<option value="">All</option>
				{#each data.facets.product_sports as s}<option>{s}</option>{/each}
			</select>
		</label>
		<label class="field">
			Sort
			<select name="sort" value={data.filters.sort} onchange={submit}>
				<option value="">Name</option>
				<option value="price_asc">Price: low → high</option>
				<option value="price_desc">Price: high → low</option>
			</select>
		</label>
		<label class="check">
			<input type="checkbox" name="in_stock" value="1" checked={data.filters.in_stock === '1'} onchange={submit} />
			In stock only
		</label>
		<button class="btn">Filter</button>
	</form>

	{#if data.products.length}
		<div class="grid">
			{#each data.products as product (product.id)}
				<ProductCard {product} />
			{/each}
		</div>
	{:else}
		<p class="empty">No gear matches those filters. <a href="/shop">Clear filters</a></p>
	{/if}
</div>

<style>
	.check {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 0.88rem;
		padding: 10px 4px;
		white-space: nowrap;
	}
	.check input {
		width: 18px;
		height: 18px;
		accent-color: var(--court-2);
	}
</style>
