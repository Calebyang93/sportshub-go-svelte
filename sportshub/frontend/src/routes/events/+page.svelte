<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import EventCard from '#lib/components/EventCard.svelte';
	import { formatDay, kindLabel } from '#lib/format.ts';
	import type { SportsEvent } from '#lib/types.ts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Group sessions under a heading per calendar day (Singapore time).
	const groups = $derived.by(() => {
		const map = new Map<string, SportsEvent[]>();
		for (const e of data.events) {
			const key = formatDay(e.starts_at);
			map.set(key, [...(map.get(key) ?? []), e]);
		}
		return [...map];
	});

	function tabHref(kind: string) {
		const p = new URLSearchParams(page.url.search);
		if (kind) p.set('kind', kind);
		else p.delete('kind');
		const s = p.toString();
		return s ? `?${s}` : '/events';
	}

	function submit(e: Event) {
		const form = (e.currentTarget as HTMLElement).closest('form')!;
		const params = new URLSearchParams();
		for (const [k, v] of new FormData(form)) if (v) params.set(k, String(v));
		goto(`?${params}`, { reset: false, replace: true });
	}
</script>

<svelte:head><title>Training, coaching & events · SportsHub</title></svelte:head>

<div class="container">
	<div class="page-head">
		<h1>Train &amp; events</h1>
		<p>Group training, 1-on-1 coaching and community events. Register in a few seconds.</p>
	</div>

	<div class="tabs" role="tablist">
		<a href={tabHref('')} class:on={!data.filters.kind} role="tab" aria-selected={!data.filters.kind}>All</a>
		{#each data.facets.event_kinds as k}
			<a href={tabHref(k)} class:on={data.filters.kind === k} role="tab" aria-selected={data.filters.kind === k}
				>{kindLabel[k]}</a
			>
		{/each}
	</div>

	<form class="filters" method="GET" onsubmit={(e) => { e.preventDefault(); submit(e); }}>
		{#if data.filters.kind}<input type="hidden" name="kind" value={data.filters.kind} />{/if}
		<label class="field grow">
			Search
			<input type="search" name="q" value={data.filters.q} placeholder="Title, coach or venue" />
		</label>
		<label class="field">
			Sport
			<select name="sport" value={data.filters.sport} onchange={submit}>
				<option value="">All sports</option>
				{#each data.facets.event_sports as s}<option>{s}</option>{/each}
			</select>
		</label>
		<button class="btn">Filter</button>
	</form>

	{#if groups.length}
		{#each groups as [day, events] (day)}
			<section class="day">
				<h2>{day}</h2>
				<div class="list">
					{#each events as event (event.id)}
						<EventCard {event} />
					{/each}
				</div>
			</section>
		{/each}
	{:else}
		<p class="empty">No upcoming sessions match. <a href="/events">See everything</a></p>
	{/if}
</div>

<style>
	.tabs {
		display: flex;
		gap: 6px;
		margin-bottom: 16px;
		overflow-x: auto;
	}
	.tabs a {
		padding: 8px 16px;
		border-radius: 999px;
		border: 1px solid var(--line);
		color: var(--ink);
		text-decoration: none;
		font-weight: 600;
		font-size: 0.9rem;
		white-space: nowrap;
	}
	.tabs a.on {
		background: var(--court);
		border-color: var(--court);
		color: #fff;
	}
	.day {
		margin-bottom: 28px;
	}
	.day h2 {
		font-size: 1.25rem;
		letter-spacing: 0.06em;
		color: var(--muted);
		margin: 0 0 10px;
	}
	.list {
		display: grid;
		gap: 12px;
	}
</style>
