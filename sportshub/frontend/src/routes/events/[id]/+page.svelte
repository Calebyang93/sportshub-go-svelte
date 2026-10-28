<script lang="ts">
	import { enhance } from '$app/forms';
	import { durationMins, formatDay, formatPrice, formatRange, kindLabel, sportGlyph } from '#lib/format.ts';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();
	const e = $derived(data.event);
	const full = $derived(e.spots_left === 0);
	const fillPct = $derived(Math.min(100, Math.round((e.registered / e.capacity) * 100)));

	let submitting = $state(false);
</script>

<svelte:head><title>{e.title} · SportsHub</title></svelte:head>

<div class="container">
	<nav class="crumbs"><a href="/events">Train &amp; events</a> / {e.title}</nav>

	<div class="layout">
		<article>
			<div class="banner kind-{e.kind}">
				<span class="glyph" aria-hidden="true">{sportGlyph(e.sport)}</span>
				<div>
					<p class="kind">{kindLabel[e.kind]} · {e.sport}</p>
					<h1>{e.title}</h1>
				</div>
			</div>

			<dl class="facts">
				<div><dt>Date</dt><dd>{formatDay(e.starts_at)}</dd></div>
				<div><dt>Time</dt><dd>{formatRange(e.starts_at, e.ends_at)} <small>({durationMins(e.starts_at, e.ends_at)} min)</small></dd></div>
				<div><dt>Venue</dt><dd>{e.venue}</dd></div>
				{#if e.coach}<div><dt>Coach</dt><dd>{e.coach}</dd></div>{/if}
				<div><dt>Level</dt><dd>{e.level}</dd></div>
				<div><dt>Price</dt><dd>{formatPrice(e.price_cents)}</dd></div>
			</dl>

			<h2>About this {e.kind === 'event' ? 'event' : 'session'}</h2>
			<p class="desc">{e.description}</p>
		</article>

		<aside class="panel">
			<div class="cap">
				<div class="cap-row">
					<strong>{full ? 'Fully booked' : `${e.spots_left} spots left`}</strong>
					<span>{e.registered}/{e.capacity}</span>
				</div>
				<div class="bar" class:full><span style="width: {fillPct}%"></span></div>
			</div>

			{#if form?.success}
				<div class="notice ok" role="status">
					<strong>You're in, {form.name}!</strong><br />
					Registration #{form.registrationId} confirmed. Details will go to {form.email}.
				</div>
				<a class="btn btn-ghost btn-block" href="/events">Browse more sessions</a>
			{:else if full}
				<p class="notice error">This session is full. Check other dates on the <a href="/events?sport={e.sport}">{e.sport} schedule</a>.</p>
			{:else}
				<form
					method="POST"
					action="?/register"
					use:enhance={() => {
						submitting = true;
						return async ({ update }) => {
							await update({ reset: false });
							submitting = false;
						};
					}}
				>
					<h2 class="form-title">Register</h2>
					{#if form?.message}<p class="notice error" role="alert">{form.message}</p>{/if}
					<label class="field">
						Full name
						<input name="name" required autocomplete="name" value={form?.values?.name ?? ''} />
					</label>
					<label class="field">
						Email
						<input name="email" type="email" required autocomplete="email" value={form?.values?.email ?? ''} />
					</label>
					<label class="field">
						<span>Phone <span class="opt">(optional)</span></span>
						<input name="phone" type="tel" autocomplete="tel" value={form?.values?.phone ?? ''} />
					</label>
					<label class="field">
						<span>Notes <span class="opt">(injuries, experience…)</span></span>
						<textarea name="notes" rows="2">{form?.values?.notes ?? ''}</textarea>
					</label>
					<button class="btn btn-accent btn-block" disabled={submitting}>
						{submitting ? 'Registering…' : e.price_cents ? `Register · ${formatPrice(e.price_cents)}` : 'Register for free'}
					</button>
					{#if e.price_cents}<p class="fine">Pay at the venue. Payment integration can be added later.</p>{/if}
				</form>
			{/if}
		</aside>
	</div>
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
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 360px;
		gap: 32px;
		align-items: start;
	}
	.banner {
		display: flex;
		align-items: center;
		gap: 18px;
		padding: 28px;
		border-radius: var(--radius);
		background: linear-gradient(135deg, var(--court), var(--court-2));
		color: #fff;
	}
	.glyph {
		font-size: 3.5rem;
	}
	.kind {
		margin: 0 0 6px;
		font-family: var(--display);
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: var(--lime);
	}
	.banner h1 {
		margin: 0;
		font-size: clamp(2rem, 5vw, 3rem);
	}
	.facts {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
		gap: 1px;
		margin: 20px 0;
		background: var(--line);
		border: 1px solid var(--line);
		border-radius: var(--radius);
		overflow: hidden;
	}
	.facts div {
		background: var(--surface);
		padding: 12px 16px;
	}
	.facts dt {
		font-size: 0.75rem;
		color: var(--muted);
		text-transform: uppercase;
		letter-spacing: 0.06em;
	}
	.facts dd {
		margin: 2px 0 0;
		font-weight: 600;
	}
	.facts small {
		color: var(--muted);
		font-weight: 400;
	}
	article h2 {
		font-size: 1.5rem;
		margin: 28px 0 8px;
	}
	.desc {
		color: var(--muted);
		max-width: 60ch;
	}
	.panel {
		position: sticky;
		top: 80px;
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 20px;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
	}
	.cap-row {
		display: flex;
		justify-content: space-between;
		margin-bottom: 6px;
		font-size: 0.9rem;
	}
	.cap-row span {
		color: var(--muted);
	}
	.bar {
		height: 8px;
		border-radius: 4px;
		background: var(--chip);
		overflow: hidden;
	}
	.bar span {
		display: block;
		height: 100%;
		background: var(--court-2);
	}
	.bar.full span {
		background: var(--clay);
	}
	form {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	.form-title {
		margin: 0;
		font-size: 1.6rem;
	}
	.opt {
		color: var(--muted);
		font-weight: 400;
	}
	.fine {
		margin: 0;
		font-size: 0.78rem;
		color: var(--muted);
		text-align: center;
	}
	@media (max-width: 860px) {
		.layout {
			grid-template-columns: 1fr;
		}
		.panel {
			position: static;
		}
		.banner {
			padding: 20px;
		}
	}
</style>
