<script lang="ts">
	import { generateSrcSet, roundTo } from "$lib/utils";
	import Picture from "./picture.svelte";
	import { Button, buttonVariants } from "./ui/button";

	type LogoProps = {
		src: string;
		alt: string;
	};

	const centralLogo: LogoProps = {
		src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/blue-nomad.png",
		alt: "Blue Nomad Logo",
	};

	const surroundingLogos: LogoProps[] = [
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/lookout-and-wonderland.png",
			alt: "Lookout and Wonderland",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/woods-copenhagen.png",
			alt: "Woods Copenhagen",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/amolin.png",
			alt: "Amolin",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/27-rosiers.png",
			alt: "27 Rosiers",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/dermalogica.png",
			alt: "Dermalogica",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/earl-of-east.png",
			alt: "Earl of East",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/arami.png",
			alt: "Arami",
		},
		{
			src: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/logos/lesse.png",
			alt: "Lesse",
		},
	];

	let innerWidth = $state(1024);
	let innerHeight = $state(768);

	const totalImages = surroundingLogos.length;

	let ovalWidth = $derived(innerWidth * 0.7);
	let ovalHeight = $derived(innerHeight * 0.5);
	let radiusX = $derived(ovalWidth / 2);
	let radiusY = $derived(ovalHeight / 2);

	function getPosition(index: number) {
		const theta = (2 * Math.PI * index) / totalImages;
		const x = roundTo(Math.cos(theta) * radiusX, 3);
		const y = roundTo(Math.sin(theta) * radiusY, 3);
		return `translate(${x}px, ${y}px)`;
	}
</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class="mt-8 pt-12 md:px-4 lg:px-8 w-full lg:flex-col">
	<div class="lg:mb-32">
		<p class="lg:text-2xl uppercase text-end">
			Some Favorites From Near &amp; Far:
		</p>
		<h1 class="uppercase text-lg lg:text-[3rem] text-end my-2 lg:my-6">
			Our In-Studio Curation
		</h1>
	</div>

	<!-- Elliptical layout for medium screens and above -->
	<div class="relative w-full h-[650px] self-center hidden md:block">
		<!-- Central logo -->
		<div
			class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2"
		>
			<!-- <img
				src={centralLogo.src}
				srcset={generateSrcSet(
					centralLogo.src,
					[400, 800, 1200, 1600],
					"webp",
					85,
				)}
				alt={centralLogo.alt}
				width="300"
				height="300"
				sizes="(max-width: 640px) 100vw, (max-width: 1024px) 70vw, (max-width: 1280px) 50vw, 33vw"
			/> -->
			<Picture
				src={centralLogo.src}
				srcset={generateSrcSet(
					centralLogo.src,
					[400, 800, 1200, 1600],
					"webp",
					85,
				)}
				alt={centralLogo.alt}
				width="300"
				height="300"
				sizes="(max-width: 640px) 100vw, (max-width: 1024px) 70vw, (max-width: 1280px) 50vw, 33vw"
				loading="lazy"
			/>
		</div>

		<!-- Surrounding logos in ellipse -->
		{#each surroundingLogos as logo, index (logo.alt)}
			<div
				class="absolute top-1/2 left-1/2"
				style:transform={getPosition(index)}
			>
				<div
					class="flex items-center justify-center"
					style:transform="translate(-50%, -50%)"
				>
					<!-- <img
						src={logo.src}
						srcset={generateSrcSet(
							logo.src,
							[400, 800, 1200, 1600],
							"webp",
							85,
						)}
						alt={logo.alt}
						width="175"
						height="75"
						loading="lazy"
						sizes="(max-width: 640px) 100vw, (max-width: 1024px) 70vw, (max-width: 1280px) 50vw, 33vw"
					/> -->
					<Picture
						src={logo.src}
						alt={logo.alt}
						width="175"
						height="75"
						sizes="(max-width: 640px) 100vw, 33vw"
						loading="lazy"
						sources={[
							{
								type: "image/webp",
								srcset: generateSrcSet(
									logo.src,
									[400, 800, 1200, 1600],
									"webp",
									85,
								),
							},
						]}
					/>
				</div>
			</div>
		{/each}
	</div>

	<!-- Stacked layout for mobile -->
	<div class="flex flex-col items-center gap-4 mt-10 lg:hidden">
		<div>
			<Picture
				src={centralLogo.src}
				srcset={generateSrcSet(
					centralLogo.src,
					[400, 800, 1200, 1600],
					"webp",
					85,
				)}
				alt={centralLogo.alt}
				width="250"
				height="200"
				loading="lazy"
			/>
		</div>

		{#each surroundingLogos as logo (logo.alt)}
			<div>
				<Picture
					src={logo.src}
					srcset={generateSrcSet(
						logo.src,
						[400, 800, 1200, 1600],
						"webp",
						85,
					)}
					alt={logo.alt}
					width="125"
					height="75"
					loading="lazy"
				/>
			</div>
		{/each}
	</div>
</div>
