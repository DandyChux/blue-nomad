<script lang="ts">
	import { Button, buttonVariants } from "./ui/button";
	import { generateSrcSet, roundTo } from "$lib/utils";

	type TestimonialProps = {
		name: string;
		profession: string;
		image: string;
		description: string;
	};

	const testimonials: TestimonialProps[] = [
		{
			name: "Andie",
			profession: "fashion producer",
			description: "Rediscover what feels good",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Andie%201.webp",
		},
		{
			name: "Gohar",
			profession: "artist",
			description: "Feel inspired",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Gohar.webp",
		},
		{
			name: "Kevin",
			profession: "model",
			description: "Refresh from within",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Kevin.webp",
		},
		{
			name: "Silvana",
			profession: "architect + interior designer",
			description: "Relax and unwind",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Tobin-162.webp",
		},
		{
			name: "Mawatle",
			profession: "thespian",
			description: "Embrace my authentic self",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Mawatle%201.webp",
		},
		{
			name: "Ron",
			profession: "Tech founder",
			description: "Feel energized!",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/testimonial/Ron%201.webp",
		},
	];

	const totalTestimonials = testimonials.length;

	let innerWidth = $state(1024);
	let innerHeight = $state(768);

	// const ovalWidth = $derived(Math.min(innerWidth * 0.75, 600));
	// const ovalHeight = $derived(Math.min(innerHeight * 0.75, 550));
	const ovalWidth = $derived(innerWidth * 0.7);
	const ovalHeight = 600;
	const radiusX = $derived(ovalWidth / 2);
	const radiusY = $derived(ovalHeight / 2);
</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class="relative w-full flex justify-center">
	<!-- Desktop Layout with Circular Arrangement -->
	{#if innerWidth >= 768}
		<div class="flex justify-center items-center h-[750px] mt-32 w-full">
			<!-- Main Title Overlay -->
			<div
				class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 flex flex-col items-center"
			>
				<h2
					class="uppercase text-xl md:text-2xl lg:text-4xl mb-2 font-medium"
				>
					Earned Love
				</h2>
				<p class="text-lg md:text-xl lg:text-2xl mb-6">
					Blue Nomad makes me...
				</p>

				<Button
					variant="link"
					href="/booking"
					class={buttonVariants({
						class: "h-auto uppercase rounded-full py-2 px-10 mt-4 pointer-events-auto font-source-code-pro",
						variant: "outline",
						size: "xl",
					})}
				>
					Your <br />
					Turn
				</Button>
			</div>

			<!-- Circular Testimonial Cards -->
			{#each testimonials as testimonial, index (testimonial.name)}
				{@const θ = (2 * Math.PI * index) / totalTestimonials}
				{@const x = Math.cos(θ) * radiusX}
				{@const y = Math.sin(θ) * radiusY}

				<div
					class="absolute top-1/2 left-1/2"
					style={"transform: translate(" + x + "px, " + y + "px);"}
				>
					<figure
						class="flex flex-col items-center z-10"
						style="transform: translate(-50%, -50%);"
					>
						<div
							class="relative w-[125px] h-[200px] overflow-hidden"
						>
							<img
								src={testimonial.image}
								srcset={generateSrcSet(
									testimonial.image,
									[400, 800, 1200, 1600],
									"webp",
									85,
								)}
								alt={`Picture of ${testimonial.name}`}
								class="object-cover size-full"
								loading="lazy"
								placeholder="blur"
								sizes="(max-width: 640px) 100vw, (max-width: 1024px) 70vw, (max-width: 1280px) 50vw, 33vw"
							/>
						</div>

						<figcaption
							class="font-bold mt-5 font-source-code-pro text-center"
						>
							<span class="uppercase">
								{testimonial.name}, {testimonial.profession}
							</span>
							<p class="font-normal">
								{testimonial.description}
							</p>
						</figcaption>
					</figure>
				</div>
			{/each}
		</div>
	{:else}
		<!-- Mobile Layout -->
		<div class="flex flex-col items-center gap-4 mt-16 lg:hidden relative">
			<div class="flex flex-col items-center mb-4">
				<h2 class="uppercase text-xl lg:text-2xl font-medium">
					Earned Love
				</h2>
				<p class="text-lg lg:text-xl">Blue Nomad makes me...</p>
			</div>

			<div class="grid grid-cols-2 gap-4 w-full px-4">
				{#each testimonials as testimonial, index (testimonial.name)}
					<div class="relative w-full h-auto p-1">
						<figure class="flex flex-col items-center mx-auto">
							<div
								class="relative w-[125px] h-[200px] my-2 overflow-hidden"
							>
								<img
									src={testimonial.image}
									srcset={generateSrcSet(
										testimonial.image,
										[100, 200, 300],
										"webp",
										85,
									)}
									alt={`Picture of ${testimonial.name}`}
									class="object-cover size-full"
									loading="lazy"
									placeholder="blur"
									sizes="(max-width: 640px) 75px, (max-width: 1280px) 100px, 125px"
								/>
							</div>

							<figcaption
								class="font-bold text-center font-source-code-pro text-sm"
							>
								<span class="uppercase">
									{testimonial.name}, {testimonial.profession}
								</span>
								<p
									class="font-normal text-center font-source-code-pro text-xs opacity-80"
								>
									{testimonial.description}
								</p>
							</figcaption>
						</figure>
					</div>
				{/each}
			</div>

			<Button
				variant="outline"
				size="xl"
				class="h-auto uppercase rounded-full py-2 px-10 mt-6 font-source-code-pro"
				href="#treatments"
				target="_blank"
				rel="noopener noreferrer"
			>
				Your <br />
				Turn
			</Button>
		</div>
	{/if}
</div>
