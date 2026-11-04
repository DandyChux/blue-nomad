<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import * as Card from "$lib/components/ui/card";
	import Picture from "$lib/components/picture.svelte";
	import { getCart } from "$lib/context/cart.svelte";
	import { generateSrcSet } from "$lib/utils";
	import { fade } from "svelte/transition";

	let { data } = $props();
	const cart = getCart();

	const product = $derived(data.product);
	const itemData = $derived(product.item_data);
	const variations = $derived(itemData.variations || []);
	const imageUrls = $derived(data.imageUrls);
	const categoryName = $derived(data.categoryName);

	// Per-variation sold-out map
	const variationSoldOut = $derived(
		variations.map((v) =>
			(v.item_variation_data?.location_overrides || []).some(
				(o) => o.sold_out === true,
			),
		),
	);

	// Default to the first in-stock variation; fall back to index 0
	let selectedVariationIndex = $state(0);
	$effect(() => {
		const firstAvailable = variationSoldOut.findIndex((s) => !s);
		if (firstAvailable !== -1 && variationSoldOut[selectedVariationIndex]) {
			selectedVariationIndex = firstAvailable;
		}
	});

	const selectedVariation = $derived(variations[selectedVariationIndex]);
	const priceCents = $derived(
		selectedVariation?.item_variation_data?.price_money?.amount || 0,
	);
	const price = $derived((priceCents / 100).toFixed(2));
	const isSoldOut = $derived(
		variationSoldOut[selectedVariationIndex] ?? false,
	);

	// Track which image is currently displayed
	let activeImageIndex = $state(0);
</script>

<svelte:head>
	<title>{itemData.name} | Blue Nomad Shop</title>
</svelte:head>

<section class="min-h-screen pt-28 lg:pt-36 px-6 md:px-12 lg:px-16 pb-20">
	<!-- Back Link -->
	<a
		href="/shop"
		class="inline-block font-source-code-pro text-[10px] uppercase tracking-widest hover:text-foreground transition-colors mb-10"
	>
		&larr; Back to Shop
	</a>

	<div
		class="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20"
	>
		<!-- Left: Images -->
		<div class="flex flex-col gap-4">
			<!-- Main Image -->
			<div class="relative w-full aspect-[3/4] bg-muted overflow-hidden">
				{#key activeImageIndex}
					<div in:fade={{ duration: 300 }} class="absolute inset-0">
						<Picture
							src={imageUrls[activeImageIndex] ||
								product.image_url}
							alt="{itemData.name} - Image {activeImageIndex + 1}"
							class="w-full h-full object-cover"
							loading="eager"
							sizes="(max-width: 1024px) 100vw, 50vw"
							sources={imageUrls[activeImageIndex]
								? [
										{
											type: "image/webp",
											srcset: generateSrcSet(
												imageUrls[activeImageIndex],
												[300, 600, 900],
												"webp",
												85,
											),
										},
									]
								: []}
						/>
					</div>
				{/key}
			</div>

			<!-- Thumbnail Strip -->
			{#if imageUrls.length > 1}
				<div class="flex gap-2 overflow-x-auto no-scrollbar">
					{#each imageUrls as url, i}
						<button
							class="shrink-0 w-20 h-20 border-2 overflow-hidden transition-colors {activeImageIndex ===
							i
								? 'border-foreground'
								: 'border-transparent hover:border-border'}"
							onclick={() => (activeImageIndex = i)}
						>
							<Picture
								src={url}
								alt="{itemData.name} thumbnail {i + 1}"
								class="w-full h-full object-cover"
								loading="lazy"
							/>
						</button>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Right: Product Details -->
		<div class="flex flex-col gap-8 lg:py-8">
			<!-- Category Badge -->
			{#if categoryName}
				<span
					class="font-source-code-pro text-[12px] uppercase tracking-widest"
				>
					{categoryName}
				</span>
			{/if}

			<!-- Name -->
			<h1
				class="uppercase text-xl lg:text-3xl font-medium leading-[0.95]"
			>
				{itemData.name}
			</h1>

			<!-- Price -->
			<p class="font-source-code-pro text-2xl">${price}</p>

			<!-- Variation Selector (if more than one) -->
			{#if variations.length > 1}
				<div class="space-y-3">
					<span
						class="font-source-code-pro text-[10px] uppercase tracking-widest block"
					>
						Options
					</span>
					<div class="flex flex-wrap gap-2">
						{#each variations as variation, i}
							{@const variationIsSoldOut = variationSoldOut[i]}
							<button
								type="button"
								disabled={variationIsSoldOut}
								aria-pressed={selectedVariationIndex === i}
								class="px-5 py-3 border font-source-code-pro text-[11px] uppercase tracking-wider transition-colors
									{selectedVariationIndex === i
									? 'bg-foreground text-background border-foreground'
									: 'border-border hover:border-foreground'}
									{variationIsSoldOut
									? 'opacity-50 cursor-not-allowed line-through hover:border-border'
									: ''}"
								onclick={() => (selectedVariationIndex = i)}
							>
								{variation.item_variation_data.name}
								{#if variation.item_variation_data.price_money}
									— ${(
										variation.item_variation_data
											.price_money.amount / 100
									).toFixed(2)}
								{/if}
								{#if variationIsSoldOut}
									<span class="ml-2">· Sold out</span>
								{/if}
							</button>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Description -->
			{#if itemData.description}
				<div class="border-t border-border pt-8">
					<p
						class="text-base leading-relaxed font-spectral lg:text-2xl"
					>
						{itemData.description}
					</p>
				</div>
			{/if}

			<Button
				class="w-full uppercase rounded-none h-16 hover:bg-transparent hover:text-foreground border border-black text-primary-foreground text-sm font-source-code-pro tracking-widest mt-4 disabled:opacity-50"
				disabled={isSoldOut || !selectedVariation}
				onclick={() =>
					selectedVariation && cart.add(product, selectedVariation)}
			>
				{isSoldOut ? "Sold Out" : `Add to Bag — $${price}`}
			</Button>

			<!-- Extra Details -->
			<div
				class="border-t border-border pt-8 space-y-4 font-source-code-pro text-[11px] uppercase tracking-widest"
			>
				{#if selectedVariation?.item_variation_data?.sku}
					<p>SKU: {selectedVariation.item_variation_data.sku}</p>
				{/if}
				<p>Available for in-studio pickup or NYC delivery.</p>
				<p>Complimentary delivery on orders over $150.</p>
			</div>
		</div>
	</div>
</section>
