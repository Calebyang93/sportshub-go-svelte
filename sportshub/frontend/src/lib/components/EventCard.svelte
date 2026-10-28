<script lang="ts">
	import type { SportsEvent } from '#lib/types.ts';
	import { dateParts, formatPrice, formatRange, kindLabel, sportGlyph } from '#lib/format.ts';

	let { event }: { event: SportsEvent } = $props();
	const d = $derived(dateParts(event.starts_at));
	const full = $derived(event.spots_left === 0);
	const fillPct = $derived(Math.min(100, Math.round((event.registered / event.capacity) * 100)));
</script>

<a class="card" href="/events/{event.id}" class:full>
	<div class="date">
		<span class="wd">{d.weekday}</span>
		<span class="dd">{d.day}</span>
		<span class="mm">{d.month}</span>
	</div>
	<div class="body">
		<div class="tags">
			<span class="kind kind-{event.kind}">{kindLabel[event.kind]}</span>
			<span class="sport">{sportGlyph(event.sport)} {event.sport}</span>
			<span class="level">{event.level}</span>
		</div>
		<h3>{event.title}</h3>
		<p class="where">
			{formatRange(event.starts_at, event.ends_at)} · {event.venue}{#if event.coach}&nbsp;· {event.coach}{/if}
		</p>
		<div class="foot">
			<div class="cap" title="{event.registered} of {event.capacity} registered">
				<div class="bar"><span style="width: {fillPct}%"></span></div>
				<span class="spots">
					{#if full}Fully booked{:else}{event.spots_left} of {event.capacity} spots left{/if}
				</span>
			</div>
			<span class="price">{formatPrice(event.price_cents)}</span>
		</div>
	</div>
</a>

<style>
	.card {
		display: grid;
		grid-template-columns: 76px 1fr;
		gap: 0;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
		overflow: hidden;
		color: inherit;
		text-decoration: none;
		transition: transform 0.15s ease, box-shadow 0.15s ease;
	}
	.card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow);
	}
	.date {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: var(--court);
		color: #fff;
		font-family: var(--display);
		text-transform: uppercase;
		line-height: 1;
		padding: 12px 0;
	}
	.full .date {
		background: var(--muted-strong);
	}
	.wd,
	.mm {
		font-size: 0.8rem;
		letter-spacing: 0.08em;
		font-weight: 600;
		color: var(--lime);
	}
	.full .wd,
	.full .mm {
		color: #fff;
	}
	.dd {
		font-size: 2.1rem;
		font-weight: 800;
		margin: 2px 0;
	}
	.body {
		padding: 14px 16px;
		display: flex;
		flex-direction: column;
		gap: 6px;
		min-width: 0;
	}
	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
		font-size: 0.75rem;
	}
	.tags > span {
		padding: 2px 8px;
		border-radius: 999px;
		background: var(--chip);
		color: var(--muted);
		font-weight: 500;
	}
	.tags .kind {
		color: var(--ink-on-accent);
		font-weight: 600;
	}
	.kind-training {
		background: var(--lime) !important;
	}
	.kind-coaching {
		background: var(--sky) !important;
	}
	.kind-event {
		background: var(--clay) !important;
		color: #fff !important;
	}
	h3 {
		margin: 0;
		font-size: 1.05rem;
		font-weight: 600;
	}
	.where {
		margin: 0;
		color: var(--muted);
		font-size: 0.85rem;
		overflow-wrap: anywhere;
	}
	.foot {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 12px;
		margin-top: 4px;
	}
	.cap {
		flex: 1;
		max-width: 260px;
	}
	.bar {
		height: 6px;
		border-radius: 3px;
		background: var(--chip);
		overflow: hidden;
	}
	.bar span {
		display: block;
		height: 100%;
		background: var(--court-2);
	}
	.full .bar span {
		background: var(--clay);
	}
	.spots {
		font-size: 0.75rem;
		color: var(--muted);
	}
	.full .spots {
		color: var(--clay);
		font-weight: 600;
	}
	.price {
		font-family: var(--display);
		font-weight: 700;
		font-size: 1.3rem;
		white-space: nowrap;
	}
</style>
