<script lang="ts">
	import { enhance } from '$app/forms';
	import { cart } from '#lib/cart.svelte.ts';
	import { formatPrice, sportGlyph } from '#lib/format.ts';
	import type { PageProps } from './$types';

	let { form }: PageProps = $props();
	let submitting = $state(false);

	const itemsJson = $derived(JSON.stringify(cart.lines.map((l) => ({ product_id: l.id, qty: l.qty }))));
</script>

<svelte:head><title>Cart · SportsHub</title></svelte:head>

<div class="container">
	<div class="page-head">
		<h1>Your cart</h1>
	</div>

	{#if form?.success}
		<div class="done">
			<p class="big">Order #{form.orderId} placed</p>
			<p>Total charged: <strong>{formatPrice(form.total)}</strong>. A confirmation will go to {form.email}.</p>
			<div class="actions">
				<a class="btn" href="/shop">Keep shopping</a>
				<a class="btn btn-ghost" href="/events">Book a session</a>
			</div>
		</div>
	{:else if !cart.ready}
		<p class="empty">Loading cart…</p>
	{:else if cart.lines.length === 0}
		<div class="empty">
			<p>Your cart is empty.</p>
			<a class="btn" href="/shop">Browse gear</a>
		</div>
	{:else}
		<div class="layout">
			<ul class="lines">
				{#each cart.lines as line (line.id)}
					<li>
						<span class="thumb" aria-hidden="true">{sportGlyph(line.sport)}</span>
						<div class="name">
							<a href="/shop/{line.id}">{line.name}</a>
							<span>{formatPrice(line.price_cents)} each</span>
						</div>
						<div class="qty">
							<button class="step" onclick={() => cart.setQty(line.id, line.qty - 1)} aria-label="Decrease quantity">−</button>
							<span aria-live="polite">{line.qty}</span>
							<button
								class="step"
								onclick={() => cart.setQty(line.id, line.qty + 1)}
								disabled={line.qty >= line.stock}
								aria-label="Increase quantity">+</button
							>
						</div>
						<span class="line-total">{formatPrice(line.price_cents * line.qty)}</span>
						<button class="remove" onclick={() => cart.remove(line.id)} aria-label="Remove {line.name}">✕</button>
					</li>
				{/each}
			</ul>

			<form
				class="summary"
				method="POST"
				action="?/checkout"
				use:enhance={() => {
					submitting = true;
					return async ({ result, update }) => {
						if (result.type === 'success') cart.clear();
						await update({ reset: false });
						submitting = false;
					};
				}}
			>
				<h2>Checkout</h2>
				<div class="sum-row"><span>Subtotal</span><strong>{formatPrice(cart.subtotal)}</strong></div>
				<p class="fine">Final prices and stock are confirmed by the store when you place the order.</p>

				{#if form?.message}<p class="notice error" role="alert">{form.message}</p>{/if}

				<input type="hidden" name="items" value={itemsJson} />
				<label class="field">
					Full name
					<input name="customer_name" required autocomplete="name" value={form?.values?.customer_name ?? ''} />
				</label>
				<label class="field">
					Email
					<input name="email" type="email" required autocomplete="email" value={form?.values?.email ?? ''} />
				</label>
				<button class="btn btn-accent btn-block" disabled={submitting}>
					{submitting ? 'Placing order…' : `Place order · ${formatPrice(cart.subtotal)}`}
				</button>
			</form>
		</div>
	{/if}
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 340px;
		gap: 28px;
		align-items: start;
	}
	.lines {
		list-style: none;
		margin: 0;
		padding: 0;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
	}
	li {
		display: grid;
		grid-template-columns: 52px minmax(0, 1fr) auto auto auto;
		align-items: center;
		gap: 14px;
		padding: 14px 16px;
	}
	li + li {
		border-top: 1px solid var(--line);
	}
	.thumb {
		width: 52px;
		height: 52px;
		border-radius: 10px;
		background: linear-gradient(135deg, var(--court), var(--court-2));
		display: grid;
		place-items: center;
		font-size: 1.6rem;
	}
	.name {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}
	.name a {
		color: var(--ink);
		font-weight: 600;
		text-decoration: none;
	}
	.name span {
		font-size: 0.8rem;
		color: var(--muted);
	}
	.qty {
		display: flex;
		align-items: center;
		gap: 8px;
		font-weight: 600;
	}
	.step,
	.remove {
		width: 30px;
		height: 30px;
		border-radius: 50%;
		border: 1px solid var(--line);
		background: var(--surface);
		color: var(--ink);
		cursor: pointer;
		font-size: 1rem;
		line-height: 1;
	}
	.step:disabled {
		opacity: 0.35;
		cursor: not-allowed;
	}
	.remove {
		border: 0;
		color: var(--muted);
	}
	.remove:hover {
		color: var(--clay);
	}
	.line-total {
		font-family: var(--display);
		font-weight: 700;
		font-size: 1.15rem;
		min-width: 80px;
		text-align: right;
	}
	.summary {
		position: sticky;
		top: 80px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 20px;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
	}
	.summary h2 {
		margin: 0;
		font-size: 1.6rem;
	}
	.sum-row {
		display: flex;
		justify-content: space-between;
		font-size: 1.05rem;
	}
	.fine {
		margin: 0;
		font-size: 0.78rem;
		color: var(--muted);
	}
	.done {
		padding: 32px;
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--radius);
	}
	.big {
		font-family: var(--display);
		font-size: 2.2rem;
		font-weight: 800;
		text-transform: uppercase;
		margin: 0 0 8px;
		color: var(--court-2);
	}
	.actions {
		display: flex;
		gap: 10px;
		margin-top: 18px;
		flex-wrap: wrap;
	}
	.empty .btn {
		margin-top: 12px;
	}
	@media (max-width: 860px) {
		.layout {
			grid-template-columns: 1fr;
		}
		.summary {
			position: static;
		}
	}
	@media (max-width: 520px) {
		li {
			grid-template-columns: 44px minmax(0, 1fr) auto;
			grid-template-areas: 'thumb name remove' 'thumb qty total';
			row-gap: 8px;
		}
		.thumb {
			grid-area: thumb;
			width: 44px;
			height: 44px;
		}
		.name {
			grid-area: name;
		}
		.qty {
			grid-area: qty;
		}
		.line-total {
			grid-area: total;
		}
		.remove {
			grid-area: remove;
		}
	}
</style>
