<script lang="ts">
	import { Input } from "$lib/components/ui/input";
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Pagination from "$lib/components/ui/pagination";
	import * as Empty from "$lib/components/ui/empty";
	import * as Item from "$lib/components/ui/item";
	import Picture from "$lib/components/picture.svelte";
	import { generateSrcSet, debounce, cn } from "$lib/utils";
	import { page } from "$app/state";
	import { goto } from "$app/navigation";
	import { fade, fly } from "svelte/transition";
	import type { CatalogItem } from "$lib/schemas";
	import Video from "$lib/components/video.svelte";
	import InfiniteMovingCards from "$lib/components/infinite-moving-cards.svelte";
	import { pressBrands } from "$lib";
	import MembershipSection from "$lib/components/membership-section.svelte";
	import UgcCarousel from "$lib/components/ugc-carousel.svelte";
	import Reviews from "$lib/components/reviews.svelte";
	import BookButton from "$lib/components/book-button.svelte";
	import { trackEvent } from "$lib/analytics.svelte.js";
	import SkinHealthGuide from "$lib/components/skin-health-guide.svelte";

	let { data } = $props();

	const services = $derived(data.services as CatalogItem[]);

	let width = $state(0);
	let height = $state(0);

	// --- URL-Based State ---
	let searchQuery = $derived(page.url.searchParams.get("q") ?? "");
	let currentPage = $derived(
		parseInt(page.url.searchParams.get("page") ?? "1"),
	);
	const perPage = 5;

	function updateFilter(key: string, value: string | number | null) {
		const url = new URL(page.url);
		if (value === null || value === "") url.searchParams.delete(key);
		else url.searchParams.set(key, String(value));
		if (key !== "page") url.searchParams.delete("page");
		goto(url, { keepFocus: true, replaceState: true, noScroll: true });
	}

	const debouncedSearch = debounce((q: string) => updateFilter("q", q), 300);

	let filteredServices = $derived(
		services.filter((s) => {
			const name = s.item_data?.name?.toLowerCase() || "";
			return name.includes(searchQuery.toLowerCase());
		}),
	);

	let paginatedServices = $derived(
		filteredServices.slice(
			(currentPage - 1) * perPage,
			currentPage * perPage,
		),
	);

	// Split into featured (first) and the rest
	const coreTreatment = $derived(
		paginatedServices.find((s) => s.id === "TWYSCIC46EIMS3SD2A6UMJ5H") ??
			paginatedServices[0] ??
			null,
	);
	const remaining = $derived(
		filteredServices.filter(
			(s) =>
				s.id !== "TWYSCIC46EIMS3SD2A6UMJ5H" &&
				s.id !== "W3E6HQXJJYYRQXT5LXSCVVSF" &&
				s.id !== "TP44ASSCFLTJZRS4THKLJMTF",
		),
	);

	// Search bar visibility toggle
	let showSearch = $state(false);

	let featuredTreatmentId = $derived(coreTreatment?.id ?? null);

	// Intersection Observer to trigger animation on scroll
	let testimonialsEl = $state<HTMLDivElement>();
	let testimonialsVisible = $state(false);

	$effect(() => {
		if (!testimonialsEl) return;
		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry.isIntersecting) {
					testimonialsVisible = true;
					observer.unobserve(testimonialsEl!);
				}
			},
			{ threshold: 1 },
		);
		observer.observe(testimonialsEl);
		return () => observer.disconnect();
	});

	// Helpers
	const getDuration = (s: CatalogItem) =>
		(
			(s.item_data.variations?.[0]?.item_variation_data
				?.service_duration || 0) / 60000
		).toFixed(0);
	const getPrice = (s: CatalogItem) =>
		(
			(s.item_data.variations?.[0]?.item_variation_data?.price_money
				?.amount || 0) / 100
		).toFixed(0);

	const isNew = (createdAt: string, daysThreshold = 14) => {
		const createdDate = new Date(createdAt);

		const now = new Date();

		const msInADay = 24 * 60 * 60 * 1000;
		const cutoffTime = now.getTime() - daysThreshold * msInADay;

		return createdDate.getTime() > cutoffTime;
	};
</script>

<svelte:head>
	<title>Treatments | Blue Nomad</title>
	<link
		rel="preload"
		href="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/videos/Blue%20Nomad%20-%20Treatment%20V1.webm"
	/>
</svelte:head>

<svelte:window bind:innerWidth={width} bind:innerHeight={height} />

<section class="w-full flex flex-col p-0">
	<!-- ======================== -->
	<!-- HERO: Full-bleed Video   -->
	<!-- ======================== -->
	<div class="relative w-full h-[70vh] overflow-hidden bg-black">
		<!-- Video Background -->
		<Video
			class="absolute inset-0 w-full h-full object-cover opacity-60"
			src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/hero/booking%20hero/hero.mp4"
		/>

		<!-- Gradient Overlay -->
		<div
			class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"
		></div>

		<!-- Hero Content -->
		<div
			class="absolute inset-0 text-brand-white flex flex-col justify-end px-6 md:px-12 lg:px-16 pb-16 lg:pb-24 w-full"
		>
			<p
				class="font-source-code-pro text-[11px] lg:text-[16px] uppercase tracking-[0.3em] text-brand-white/60 mb-4 font-extrabold text-center lg:text-start"
				in:fly={{ y: 10, duration: 600, delay: 200 }}
			>
				Blue Nomad Skin Health Studio
			</p>
			<h1
				class="uppercase text-center lg:text-left text-[32px] lg:text-[72px] tracking-tighter font-light leading-[0.9] mb-6"
				in:fly={{ y: 20, duration: 600, delay: 300 }}
			>
				<span>Your Skin</span>
				<br />
				Our Practice
				<!-- <br />
				<span class="pl-8 md:pl-16">Everything Around</span> -->
			</h1>
			<p
				class="text-brand-white/70 md:text-[18px] max-w-xl leading-relaxed font-light text-center lg:text-start"
				in:fly={{ y: 20, duration: 600, delay: 400 }}
			>
				Our core <em>Facial Skin Therapy (ST)</em> is a personalized treatment
				that improves skin health and clarity for all skin.
			</p>

			<Button
				variant="link"
				href={`/booking/${featuredTreatmentId}`}
				onclick={() =>
					trackEvent("book_facial_st_click", {
						props: { "Book Core Facial ST": "hero" },
					})}
				class={buttonVariants({
					variant: "outline",
					class: "border-brand-white text-brand-white rounded-full mt-10 border self-center px-8 uppercase tracking-wide font-source-code-pro hover:bg-brand-white hover:text-black",
					size: "lg",
				})}
			>
				Book Facial ST
			</Button>
			<!-- <BookButton
				href={`/booking/${featuredTreatmentId}`}
				label="Book Core Facial ST"
				event="book_facial_st_click"
				location="hero"
			/> -->
		</div>
	</div>

	{#if paginatedServices.length === 0}
		<div class="w-full flex justify-center py-32">
			<Empty.Root
				class="border border-dashed w-full max-w-xl min-h-[400px] flex flex-col items-center justify-center bg-transparent"
			>
				<Empty.Header>
					<Empty.Title
						class="uppercase font-source-code-pro font-normal tracking-widest text-sm"
					>
						No treatments found.
					</Empty.Title>
				</Empty.Header>
			</Empty.Root>
		</div>
	{:else}
		<div class="relative overflow-hidden px-6 md:px-10 py-2 lg:py-8">
			<div
				class="mx-auto max-w-[1100px] space-x-4 space-y-4 md:space-x-8 lg:space-x-10 lg:space-y-8"
			>
				<div class="flex items-center justify-between w-full">
					<p class="text-[36px] leading-none lg:text-[48px]">
						Before
					</p>

					<p
						class="text-left font-source-code-pro text-[12px] font-medium uppercase leading-[1.35] tracking-[0.2em] lg:text-right lg:text-[1.5rem]"
					>
						10 days after
						<br />
						first facial ST
					</p>
				</div>

				<div class="flex items-center space-x-4 justify-between w-full">
					<div class="relative group ring-2 ring-card rounded-[15px]">
						<a
							class="absolute inset-0"
							href={`/booking/${coreTreatment?.id}`}
						/>

						<div
							class="aspect-[0.85] w-auto lg:h-[600px] mx-auto overflow-hidden rounded-[15px] shadow-[0_0_0_1px_rgba(255,255,255,0.28)_inset]"
						>
							<Picture
								src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/studio/Treatment%20After%20(Former).webp"
								alt="Before facial skin therapy result"
								class="size-full object-cover object-[center_75%]"
								loading="eager"
								width={720}
								height={950}
								sizes="(max-width: 1024px) 50vw, 50vw"
								sources={[
									{
										type: "image/webp",
										srcset: generateSrcSet(
											"https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/studio/Treatment%20After%20(Former).webp",
											[480, 768, 1024, 1440],
											"webp",
											85,
										),
									},
								]}
							/>
						</div>
						<span
							class="absolute bottom-0 text-[12px] lg:text-base p-2 w-full text-center uppercase font-source-code-pro font-medium bg-card text-brand-white rounded-b-[15px] transition-opacity duration-300"
						>
							{getDuration(coreTreatment)} Min — ${getPrice(
								coreTreatment,
							)}
						</span>
					</div>

					<div class="relative group ring-2 ring-card rounded-[15px]">
						<a
							class="absolute inset-0"
							href={`/booking/${coreTreatment?.id}`}
						/>

						<div
							class="aspect-[0.85] w-auto lg:h-[600px] mx-auto overflow-hidden rounded-[15px] shadow-[0_0_0_1px_rgba(255,255,255,0.28)_inset]"
						>
							<Picture
								src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/studio/Treatment%20(After).webp"
								alt="After 10 days facial skin therapy result"
								class="size-full object-cover"
								loading="eager"
								width={720}
								height={950}
								sizes="(max-width: 1024px) 50vw, 50vw"
								sources={[
									{
										type: "image/webp",
										srcset: generateSrcSet(
											"https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/studio/Treatment%20(After).webp",
											[480, 768, 1024, 1440],
											"webp",
											85,
										),
									},
								]}
							/>
						</div>
						<span
							class="absolute bottom-0 text-[12px] lg:text-base p-2 w-full text-center uppercase font-source-code-pro font-medium bg-card text-brand-white rounded-b-[15px] transition-opacity duration-300"
						>
							Book Now
						</span>
					</div>
				</div>

				<div class="flex items-center justify-between w-full">
					<p
						class="font-source-code-pro text-[12px] font-semibold leading-[1.2em] tracking-[1.6px] md:text-[16px] lg:w-[450px]"
					>
						Perosnalized treatments for all skin, with expertise in
						breakouts, dark spots, and sensitive skin.
					</p>

					<div class="col-start-2 row-start-3 justify-self-end">
						<p class="text-[36px] leading-none lg:text-[48px]">
							After
						</p>
					</div>
				</div>

				<div class="flex items-center justify-around w-full">
					<Button
						variant="link"
						href={`/booking/${coreTreatment.id}`}
						class={buttonVariants({
							size: "lg",
							variant: "outline",
							class: "col-span-1 col-start-1 row-start-4 w-fit min-w-0 max-w-full place-self-center lg:place-self-end rounded-full px-8 sm:px-6 lg:px-14 uppercase font-source-code-pro",
						})}
					>
						Book Facial ST
					</Button>
					<Button
						variant="link"
						href="/shop"
						class={buttonVariants({
							size: "lg",
							variant: "outline",
							class: "w-fit min-w-0 max-w-full rounded-full px-8 sm:px-6 lg:px-14 col-span-1 col-start-2 row-start-4 place-self-start uppercase font-source-code-pro",
						})}
					>
						Shop
					</Button>
				</div>
			</div>
		</div>

		<!-- <span
			class="text-[32px] lg:text-[44px] uppercase text-center font-medium"
			>also available</span
		>
		<div
			class="mx-auto grid w-full max-w-400 auto-rows-fr grid-cols-1 justify-center gap-px px-4 sm:grid-cols-2 lg:gap-4"
			class:lg:grid-cols-1={remaining.length === 1}
			class:lg:grid-cols-2={remaining.length === 2}
			class:lg:grid-cols-3={remaining.length === 3}
			class:lg:grid-cols-4={remaining.length >= 4}
			in:fade
		>
			{#each remaining as service (service.id)}
				<Card.Root
					class="group relative flex h-full w-full min-w-0 flex-col bg-background lg:bg-transparent shadow-none ring-0 py-0 rounded-none gap-0 my-2 lg:rounded-lg"
				>
					<a
						href="/booking/{service.id}"
						aria-label={`Book ${service.item_data.name}`}
						class="absolute inset-0 z-10"
					></a>

					{#if isNew(service.created_at)}
						<span
							class="absolute top-0 left-0 bg-background text-foreground lg:bg-card lg:text-brand-white uppercase rounded-sm z-10 py-2 px-8 font-source-code-pro font-semibold"
						>
							NEW
						</span>
					{/if}

					<Card.Content
						class="relative rounded-lg lg:bg-card p-2 border-2 border-transparent group-hover:border-card"
					>
						<div class="aspect-4/5 overflow-hidden">
							<Picture
								src={service.image_url || ""}
								alt={service.item_data.name}
								class="block size-full object-cover"
								loading="lazy"
								sizes="(min-width: 1024px) 25vw, 50vw"
								sources={service.image_url
									? [
											{
												type: "image/webp",
												srcset: generateSrcSet(
													service.image_url,
													[320, 480, 640, 768],
													"webp",
													80,
												),
											},
										]
									: []}
							/>
						</div>

						<span
							class="hidden lg:block absolute bottom-0 text-[12px] lg:text-base p-2 w-full text-center uppercase font-source-code-pro font-medium bg-card text-brand-white rounded-b-[15px] transition-opacity duration-300"
						>
							Book Treatment
						</span>
					</Card.Content>

					<Card.Footer
						class="items-center lg:items-start justify-between text-foreground lg:bg-transparent lg:flex-col lg:pt-4"
					>
						<div class="lg:self-end order-2 lg:order-1">
							<span
								class="font-source-code-pro text-[12px] lg:text-[14px] font-bold uppercase"
							>
								{getDuration(service)} min{getPrice(service) !==
								"0"
									? ` — $${getPrice(service)}`
									: " – Member Rate"}
							</span>
						</div>
						<h2
							class="text-[14px] lg:text-[32px] order-1 lg:order-2 tracking-tight line-clamp-1 text-ellipsis"
						>
							{service.item_data.name}
						</h2>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div> -->
	{/if}

	<!-- Press Ticker -->
	<InfiniteMovingCards
		items={pressBrands}
		direction="right"
		speed="normal"
		class="mask-[linear-gradient(to_right,transparent_0%,white_20%,white_100%)] w-full max-w-[unset]"
	/>

	<!-- Social & Testimonials -->
	<UgcCarousel />

	<Reviews />

	<MembershipSection />

	<div class="relative w-full py-8 px-12 lg:p-20">
		<div class="absolute inset-0 m-0 -z-10">
			<enhanced:img
				src={"https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/studio/IMG%20Glass.webp"}
				alt="Background"
				class="object-cover xl:object-cover object-center size-full"
			/>
		</div>

		<div
			class="w-full flex flex-col lg:flex-row space-y-8 lg:space-y-0 lg:space-x-6"
		>
			<Item.Root
				variant="outline"
				class="px-8 flex-col lg:flex-row text-cold-ivory border-2 border-warm-ivory rounded-2xl"
			>
				<Item.Header
					class="uppercase text-primary-foreground/70 font-source-code-pro"
					>Custom Facial</Item.Header
				>
				<Item.Content>
					<Item.Title
						class="uppercase text-warm-ivory text-2xl w-full justify-center lg:justify-start"
						>Better Skin Health</Item.Title
					>
					<Item.Description
						class="uppercase text-primary-foreground/75 font-source-code-pro text-base lg:text-start text-center"
					>
						Experience our core Facial ST - $235
					</Item.Description>
				</Item.Content>
				<Item.Actions class="place-self-center">
					<Button
						variant="link"
						href="/booking/TWYSCIC46EIMS3SD2A6UMJ5H"
						class={buttonVariants({
							variant: "outline",
							class: "uppercase border-warm-ivory text-warm-ivory rounded-full px-8",
						})}>Book Facial ST</Button
					>
				</Item.Actions>
			</Item.Root>

			<Item.Root
				variant={"outline"}
				class="px-8 flex-col lg:flex-row text-cold-ivory border-2 border-warm-ivory rounded-2xl"
			>
				<Item.Header
					class={"uppercase text-primary-foreground/70 font-source-code-pro text-center"}
					>Skin Consult</Item.Header
				>
				<Item.Content>
					<Item.Title
						class="line-clamp-none uppercase text-warm-ivory text-2xl w-full"
						>Not sure where to start?</Item.Title
					>
					<Item.Description
						class={"uppercase text-primary-foreground/75 font-source-code-pro text-base text-center lg:text-start"}
						>Begin with a consultation - $120</Item.Description
					>
				</Item.Content>
				<Item.Actions>
					<Button
						variant="link"
						href="/booking/VPBLSPCLMIHKLXHLPE2BY3HL"
						class={buttonVariants({
							variant: "outline",
							class: "uppercase border-warm-ivory text-warm-ivory rounded-full px-8",
						})}>Book Consult</Button
					>
				</Item.Actions>
			</Item.Root>
		</div>
	</div>

	<SkinHealthGuide />
</section>
