<script lang="ts">
	import EventCard from '#lib/components/EventCard.svelte';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head><title>SportsHub · Gear, coaching & events</title></svelte:head>

<section class="hero">
	<div class="container inner">
		<p class="kicker">Gear up · Train · Compete</p>
		<h1>Everything you need<br />for your <span>next session</span></h1>
		<p class="lede">
			Shop sports equipment and book training, coaching and community events across Singapore — all in one place.
		</p>
		<div class="cta">
			<a class="btn btn-accent" href="/events">Book a session</a>
			<a class="btn btn-ghost light" href="/shop">Browse gear</a>
		</div>
		<dl class="stats">
			<div><dt>Upcoming sessions</dt><dd>{data.counts.events}</dd></div>
			<div><dt>Items in stock</dt><dd>{data.counts.products}</dd></div>
		</dl>
	</div>
</section>

<section class="container block">
	<div class="head">
		<h2>Coming up</h2>
		<a href="/events">All sessions →</a>
	</div>
	{#if data.upcoming.length}
		<div class="events">
			{#each data.upcoming as event (event.id)}
				<EventCard {event} />
			{/each}
		</div>
	{:else}
		<p class="empty">No upcoming sessions yet.</p>
	{/if}
</section>

<section class="container block">
	<div class="head">
		<h2>Featured gear</h2>
		<a href="/shop">Shop all →</a>
	</div>
	<div class="grid">
		{#each data.featured as product (product.id)}
			<ProductCard {product} />
		{/each}
	</div>
</section>

<style>
	.hero {
		background:
			radial-gradient(circle at 85% 20%, rgb(198 244 50 / 0.18), transparent 40%),
			linear-gradient(160deg, var(--court) 0%, #0a2a20 100%);
		color: #fff;
		overflow: hidden;
	}
	.inner {
		padding-top: 56px;
		padding-bottom: 48px;
	}
	.kicker {
		margin: 0 0 12px;
		font-family: var(--display);
		font-weight: 700;
		letter-spacing: 0.14em;
		text-transform: uppercase;
		color: var(--lime);
	}
	h1 {
		margin: 0;
		font-size: clamp(2.6rem, 8vw, 5rem);
	}
	h1 span {
		color: var(--lime);
	}
	.lede {
		max-width: 520px;
		color: rgb(255 255 255 / 0.78);
		font-size: 1.05rem;
		margin: 18px 0 24px;
	}
	.cta {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}
	.light {
		color: #fff;
		border-color: rgb(255 255 255 / 0.35);
	}
	.light:hover {
		background: rgb(255 255 255 / 0.08) !important;
	}
	.stats {
		display: flex;
		gap: 36px;
		margin: 36px 0 0;
	}
	.stats div {
		display: flex;
		flex-direction: column-reverse;
	}
	.stats dt {
		font-size: 0.8rem;
		color: rgb(255 255 255 / 0.65);
	}
	.stats dd {
		margin: 0;
		font-family: var(--display);
		font-size: 2.6rem;
		font-weight: 800;
		line-height: 1;
	}
	.block {
		margin-top: 44px;
	}
	.head {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		margin-bottom: 16px;
	}
	.head h2 {
		margin: 0;
		font-size: 2rem;
	}
	.head a {
		font-weight: 600;
		text-decoration: none;
	}
	.events {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(min(100%, 440px), 1fr));
		gap: 14px;
	}
</style>
