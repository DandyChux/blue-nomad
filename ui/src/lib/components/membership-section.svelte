<script lang="ts">
	import BookButton from "./book-button.svelte";
	import Picture from "./picture.svelte";
	import { Button, buttonVariants } from "./ui/button";

	const membershipRows = [["Existing members", "60 min – $185"]];
</script>

<section class="bn-section">
	<div class="bn-inner grid">
		<img
			class="photo"
			src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/membership/membership.jpg"
			alt="Blue Nomad client after a Facial ST"
			loading="lazy"
		/>
		<div class="copy">
			<p class="bn-kicker">Membership</p>
			<h2 class="bn-title">Facial ST Membership</h2>
			<p class="bn-body bn-brief">
				Monthly access to our core Facial ST at a preferred rate, plus
				ongoing skin health support and wellness perks. Existing members
				may book below. For new membership enquiries, please email
				memberships@bluenomad.nyc.
			</p>
			<dl class="rows">
				{#each membershipRows as [label, value]}
					<div class="row">
						<dt>{label}</dt>
						<dd>{value}</dd>
					</div>
				{/each}
			</dl>
			<div class="cta">
				<BookButton
					href="https://bluenomadworld.com/booking/2BZBFRDWP6J42MLFCG7GIUOE"
					label="Member Facial ST"
					event="book_membership_click"
					location="membership"
				/>
			</div>
		</div>
	</div>
</section>

<style>
	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 64px;
		align-items: center;
	}
	.photo {
		width: 100%;
		height: 560px;
		object-fit: cover;
		border-radius: var(--bn-radius);
	}
	.copy {
		display: flex;
		flex-direction: column;
		gap: 22px;
	}
	.rows {
		margin: 0;
		border-top: 1px solid var(--color-gray-500);
	}
	.row {
		display: flex;
		justify-content: space-between;
		padding: 14px 0;
		border-bottom: 1px solid var(--color-gray-500);
		font-family: var(--font-source-code-pro);
		font-weight: 500;
		font-size: 14px;
		letter-spacing: 1px;
		text-transform: uppercase;
		font-variant-numeric: tabular-nums;
	}
	dt,
	dd {
		margin: 0;
	}
	@media (max-width: 1280px) {
		.grid {
			grid-template-columns: 1fr;
			gap: 32px;
		}
		.photo {
			height: 420px;
		}
		.row {
			font-size: 12px;
		}
		.cta {
			align-self: center;
		}
	}
</style>
