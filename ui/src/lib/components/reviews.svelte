<script lang="ts">
	import { onDestroy, type SvelteComponent } from "svelte";
	import Rating from "./rating.svelte";
	import * as Carousel from "./ui/carousel";
	import { fly } from "svelte/transition";
	import type { CarouselAPI } from "./ui/carousel/context";

	let atStart = $state(false);
	let atEnd = $state(false);
	let carouselApi: CarouselAPI | undefined;
	let updateBoundaries: (() => void) | undefined;

	const testimonials = [
		{
			quote: "Calming, relaxing, informative, honest. Highly recommend! My skin feels and looks great.",
			author: "Anna",
			location: "Brooklyn",
			rating: 5,
		},
		{
			quote: "Come here for a bespoke and private treatment in a space that feels like…a new home.",
			author: "Siba",
			location: "Bed-Stuy",
			rating: 5,
		},
		{
			quote: "What an incredible experience! 10/10 recommend!!! I'm hooked.",
			author: "Ron",
			location: "West Village",
			rating: 5,
		},
	];

	function handleCarouselApi(api: CarouselAPI | undefined) {
		if (!api) return;

		carouselApi = api;
		updateBoundaries = () => {
			atStart = !api.canScrollPrev();
			atEnd = !api.canScrollNext();
		};

		updateBoundaries();
		api.on("select", updateBoundaries);
		api.on("reInit", updateBoundaries);
	}

	onDestroy(() => {
		if (!carouselApi || !updateBoundaries) return;
		carouselApi.off("select", updateBoundaries);
		carouselApi.off("reInit", updateBoundaries);
	});
</script>

<Carousel.Root
	class="max-w-[750px] my-4 px-6 md:px-12 lg:px-16 lg:mx-auto"
	opts={{
		loop: false,
	}}
	setApi={handleCarouselApi}
>
	<Carousel.Content>
		{#each testimonials as testimonial, i}
			<Carousel.Item>
				<div
					class="flex flex-col lg:flex-row font-source-code-pro"
					in:fly={{ y: 40, duration: 600, delay: i * 250 }}
				>
					<div
						class="flex justify-between items-center lg:items-start lg:justify-normal lg:flex-col basis-[40%]"
					>
						<p class=" text-sm font-bold tracking-[0.2em]">
							<span class="uppercase">{testimonial.author}</span>, {testimonial.location}
						</p>
						<span
							class="font-source-code-pro text-foreground/70 text-xs uppercase tracking-wide"
							>Verified Nomad</span
						>
					</div>
					<div class="flex flex-col gap-4 basis-[60%]">
						<Rating value={testimonial.rating} />
						<p
							class={[
								"text-sm md:text-base lg:text-lg font-medium tracking-wide w-[250px] lg:w-auto md:text-left",
								{
									"text-left": atStart,
									"text-center": !atStart && !atEnd,
									"text-right": atEnd,
								},
							]}
						>
							{testimonial.quote}
						</p>
					</div>
				</div>
			</Carousel.Item>
		{/each}
	</Carousel.Content>
	<Carousel.Previous
		class="opacity-50 left-4 disabled:hidden"
		disabled={atStart}
	/>
	<Carousel.Next
		class="opacity-50 right-4 disabled:hidden"
		disabled={atEnd}
	/>
</Carousel.Root>
