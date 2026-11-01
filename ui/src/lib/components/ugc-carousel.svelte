<script lang="ts">
	import BookButton from "./book-button.svelte";
	import UgcCard from "./ugc-card.svelte";
	import { Button, buttonVariants } from "./ui/button";

	// Still carousel: nothing moves on its own. Desktop shows three posts with the next one fading out at the edge
	// (a 20% fade on each side that has more posts); phones show one post with the next peeking in.
	// Back / Next arrows (or a swipe) move one post at a time.
	let { page = "home" } = $props();

	const social = [
		{
			type: "video",
			autoplay: true,
			label: "TikTok",
			title: "Zach",
			caption: "“Favorite place to get a facial in New York.”",
			handle: "@zachboussatta",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-zachboussatta.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/zach-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/zach.mp4",
			link: "https://www.tiktok.com/@zachboussatta/video/7662920524040785165",
		},
		{
			type: "video",
			label: "Instagram",
			title: "Indya",
			caption: "“Wellness means pouring into myself.”",
			handle: "@indyabrown",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-indyabrown.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/indya-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/indya.mp4",
			link: "https://www.instagram.com/indyabrown/",
		},
		{
			type: "still",
			label: "Instagram",
			title: "Olga",
			caption: "Olga at the studio.",
			handle: "@olgayanul",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-olgayanul.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Photos/olga.jpg",
			link: "https://www.instagram.com/p/DFNtdpGx9OT/",
		},
		{
			type: "still",
			label: "Instagram",
			title: "Unconventional Beauty",
			caption: "“When looking in the mirror we question our own beauty.”",
			handle: "@bluenomadworld",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-bluenomadworld.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Photos/unconventional-beauty.jpg",
			link: "https://www.instagram.com/p/DP4Y1uGCEfN/",
		},
		{
			type: "still",
			label: "Instagram",
			title: "Back in Stock",
			caption: "A few of our favourite skin hydration essentials.",
			handle: "@bluenomadworld",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-bluenomadworld.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Photos/back-in-stock.jpg",
			link: "https://www.instagram.com/p/DcCAzN0o-ep/",
		},
		{
			type: "video",
			label: "Instagram",
			title: "Jeremy",
			caption: "Jeremy at the studio.",
			handle: "@jeremyasgari",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-jeremyasgari.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/jeremy-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/jeremy.mp4",
			link: "https://www.instagram.com/p/DO_x46iEXXp/",
		},
		{
			type: "video",
			label: "Instagram",
			title: "Kameron",
			caption: "Leaving facials educated and with good skin.",
			handle: "@kameronmonet",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-kameronmonet.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/kameron-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/kameron.mp4",
			link: "https://www.instagram.com/p/DPW-egwiUpD/",
		},
		{
			type: "video",
			label: "TikTok",
			title: "Reed",
			caption: "Black women owned spa, Blue Nomad Midtown.",
			handle: "@reed.g_",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-reed.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/reed-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/reed.mp4",
			link: "https://www.tiktok.com/@reed.g_/video/7666871387662847245",
		},
		{
			type: "video",
			label: "Instagram",
			title: "Kevin",
			caption: "A while back with Kevin ✨",
			handle: "@bluenomadworld",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-bluenomadworld.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/kevin-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/kevin.mp4",
			link: "https://www.instagram.com/p/DNGqq40q5QT/",
		},
		{
			type: "video",
			label: "Instagram",
			title: "Miles",
			caption: "Miles at the studio.",
			handle: "@makeupbymilessss",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-makeupbymilessss.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/miles-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/miles.mp4",
			link: "https://www.instagram.com/p/DRzxWPJDoSe/",
		},
		// Bonus slide: Carey.
		{
			type: "video",
			label: "Instagram",
			title: "Carey",
			caption: "Come with me to get a Blue Nomad facial.",
			handle: "@thedaintydetailsblog",
			avatar: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Avatar%20Assets/avatar-thedaintydetailsblog.jpg",
			image: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Video%20Stills/carey-still.jpg",
			video: "https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/UGC%20Ticker%20Assets/UGC%20Ticker%20Videos/carey.mp4",
			link: "https://www.instagram.com/p/DQIRV60ACXz/",
		},
	];

	let rail = $state<HTMLDivElement>();
	let active = $state(0);
	let atStart = $state(true);
	let atEnd = $state(false);
	const last = social.length - 1;

	function step(): number {
		const railElement = rail;
		if (!railElement) return 0;

		const first = railElement.children.item(0) as HTMLElement | null;
		const second = railElement.children.item(1) as HTMLElement | null;

		return first && second
			? second.offsetLeft - first.offsetLeft
			: railElement.clientWidth;
	}

	function go(i: number) {
		const railElement = rail;
		if (!railElement) return;

		active = Math.max(0, Math.min(last, i));
		railElement.scrollTo({ left: active * step(), behavior: "smooth" });
	}

	function onScroll() {
		const railElement = rail;
		if (!railElement) return;

		const distance = step();
		active =
			distance > 0
				? Math.max(
						0,
						Math.min(
							last,
							Math.round(railElement.scrollLeft / distance),
						),
					)
				: 0;

		atStart = railElement.scrollLeft < 4;
		atEnd =
			railElement.scrollLeft + railElement.clientWidth >
			railElement.scrollWidth - 4;
	}
</script>

<section class="bn-section social">
	<div class="bn-inner head">
		<div>
			<p class="bn-kicker">Skin Health and Beyond</p>
			<h2 class="bn-title">@bluenomadworld</h2>
		</div>
		<BookButton
			href="https://www.instagram.com/bluenomadworld"
			label="Follow on Instagram"
			event="social_follow"
			location="social"
			external
		/>
	</div>

	<div
		class="carousel"
		role="region"
		aria-roledescription="carousel"
		aria-label="Posts from @bluenomadworld"
	>
		<div
			class="rail"
			class:fade-l={!atStart}
			class:fade-r={!atEnd}
			bind:this={rail}
			onscroll={onScroll}
		>
			{#each social as item (item.title)}
				<div class="slide">
					<UgcCard {item} />
				</div>
			{/each}
		</div>
		<button
			class="arrow prev"
			onclick={() => go(active - 1)}
			disabled={atStart}
			aria-label="Previous post"
		>
			<img
				src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/arrows/arrow-left.svg"
				alt=""
			/>
		</button>
		<button
			class="arrow next"
			onclick={() => go(active + 1)}
			disabled={atEnd}
			aria-label="Next post"
		>
			<img
				src="https://blue-nomad.nyc3.cdn.digitaloceanspaces.com/assets/arrows/arrow-right.svg"
				alt=""
			/>
		</button>
	</div>
</section>

<style>
	.social {
		gap: 40px;
		padding-left: 0;
		padding-right: 0;
	}
	.head {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		flex-wrap: wrap;
		gap: 16px;
		padding: 0 var(--bn-gutter);
		box-sizing: content-box;
	}
	.head .bn-title {
		margin-top: 12px;
	}

	.carousel {
		--card-w: 300px;
		--fade-l: 0%;
		--fade-r: 100%;
		position: relative;
		width: 100%;
	}
	.rail {
		display: flex;
		gap: 24px;
		overflow-x: auto;
		scroll-snap-type: x mandatory;
		scrollbar-width: none;
		-webkit-mask-image: linear-gradient(
			90deg,
			transparent 0%,
			#000 var(--fade-l),
			#000 var(--fade-r),
			transparent 100%
		);
		mask-image: linear-gradient(
			90deg,
			transparent 0%,
			#000 var(--fade-l),
			#000 var(--fade-r),
			transparent 100%
		);
	}
	/* 20% fade on each edge that still has posts beyond it. */
	.rail.fade-l {
		--fade-l: 20%;
	}
	.rail.fade-r {
		--fade-r: 80%;
	}
	.rail::-webkit-scrollbar {
		display: none;
	}
	.slide {
		flex: 0 0 var(--card-w);
		scroll-snap-align: start;
	}

	/* Arrow graphic: assets/arrows/arrow-left.svg and arrow-right.svg (black circle outline with a chevron). */
	.arrow {
		position: absolute;
		top: 50%;
		width: 44px;
		height: 44px;
		margin-top: -22px;
		padding: 0;
		border: 0;
		background: none;
		cursor: pointer;
		transition: opacity 0.2s ease;
	}
	.arrow img {
		width: 100%;
		height: 100%;
		display: block;
	}
	.arrow.prev {
		left: 24px;
	}
	.arrow.next {
		right: 24px;
	}
	.arrow:disabled {
		opacity: 0;
		pointer-events: none;
	}

	@media (max-width: 1280px) {
		.social {
			gap: 28px;
		}
		.head {
			flex-direction: column;
			align-items: flex-start;
			padding: 0 var(--bn-gutter-phone);
			box-sizing: border-box;
		}
		.head .bn-title {
			font-size: 30px;
		}
		/* Phone: the Follow on Instagram button sits centred under the heading. */
		.head :global(.bn-btn) {
			align-self: center;
		}
		.rail {
			gap: 12px;
		}
		.arrow {
			width: 36px;
			height: 36px;
			margin-top: -18px;
		}
		.arrow.prev {
			left: 6px;
		}
		.arrow.next {
			right: 6px;
		}
	}
	@media (prefers-reduced-motion: reduce) {
		.rail {
			scroll-behavior: auto;
		}
	}
</style>
