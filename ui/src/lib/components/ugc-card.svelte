<script lang="ts">
	import { trackEvent } from "$lib/analytics.svelte";
	import { generateSrcSet } from "$lib/utils";
	import { Volume2, VolumeOff, Play } from "@lucide/svelte";

	// item.type: 'video' | 'still' | 'embed' (live TikTok player)
	// Video cards: tap the video to pause or play (a play button shows while paused), and a speaker
	// button turns sound on or off. The one with `autoplay: true` (Zach) starts playing muted on its own;
	// the others start paused. Only one card has sound at a time: turning one on mutes the others.
	let { item } = $props();

	let video = $state<HTMLVideoElement>();
	let muted = $state(true);
	let playing = $state(false);

	$effect(() => {
		const videoElement = video;
		if (item.type !== "video" || !videoElement) return;

		const reduce = window.matchMedia(
			"(prefers-reduced-motion: reduce)",
		).matches;

		if (item.autoplay && !reduce) {
			videoElement.muted = true;
			videoElement
				.play()
				.then(() => (playing = true))
				.catch(() => {});
		}

		const onSound = (event: Event) => {
			if ((event as CustomEvent<string>).detail === item.title) return;
			muted = true;
			videoElement.muted = true;
		};

		window.addEventListener("bn-sound", onSound);
		return () => window.removeEventListener("bn-sound", onSound);
	});

	function togglePlay() {
		const videoElement = video;
		if (!videoElement) return;

		if (playing) {
			videoElement.pause();
			playing = false;
		} else {
			videoElement
				.play()
				.then(() => (playing = true))
				.catch(() => {});
			// Keep the existing trackEvent call here.
		}
	}

	function toggleSound() {
		const videoElement = video;
		if (!videoElement) return;

		if (muted) {
			window.dispatchEvent(
				new CustomEvent("bn-sound", { detail: item.title }),
			);
			muted = false;
			videoElement.muted = false;
			// Keep the existing trackEvent call here.
		} else {
			muted = true;
			videoElement.muted = true;
		}
	}
</script>

<article class="card" class:embed={item.type === "embed"}>
	{#if item.type === "embed"}
		<iframe
			src={item.embed}
			title="{item.title} on {item.label}"
			loading="lazy"
			allow="autoplay; fullscreen; encrypted-media"
		></iframe>
	{:else if item.type === "video"}
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={video}
			src={item.video}
			poster={item.image}
			muted
			loop
			playsinline
			preload={item.autoplay ? "auto" : "none"}
		></video>
		<button
			class="tap"
			onclick={togglePlay}
			aria-pressed={playing}
			aria-label={playing
				? `Pause ${item.title}’s video`
				: `Play ${item.title}’s video`}
		>
			{#if !playing}
				<span class="play"
					><svg viewBox="0 0 24 24" aria-hidden="true"
						><path d="M8 5v14l11-7z" /></svg
					></span
				>
			{/if}
		</button>
		<button
			class="sound"
			onclick={toggleSound}
			aria-pressed={!muted}
			aria-label={muted
				? `Turn on sound for ${item.title}’s video`
				: `Mute ${item.title}’s video`}
		>
			{#if muted}
				<svg viewBox="0 0 24 24" aria-hidden="true"
					><path d="M11 5 6 9H2v6h4l5 4z" /><line
						x1="22"
						y1="9"
						x2="16"
						y2="15"
					/><line x1="16" y1="9" x2="22" y2="15" /></svg
				>
			{:else}
				<svg viewBox="0 0 24 24" aria-hidden="true"
					><path d="M11 5 6 9H2v6h4l5 4z" /><path
						d="M15.54 8.46a5 5 0 0 1 0 7.07"
					/><path d="M19.07 4.93a10 10 0 0 1 0 14.14" /></svg
				>
			{/if}
		</button>
	{:else}
		<img
			src={item.image}
			alt="{item.title}, from {item.handle}"
			loading="lazy"
			sizes="(max-width: 400px) 200px 300px"
			srcset={generateSrcSet(item.image, [200, 300], "webp", 100)}
		/>
	{/if}

	<!-- The badge at the top links to the original post (or the account), so people can still follow. -->
	<a
		class="badge"
		href={item.link}
		target="_blank"
		rel="noopener noreferrer"
		aria-label="{item.title} on {item.label} (opens the post)"
		onclick={() =>
			trackEvent("ugc_follow_click", {
				props: {
					title: item.title,
					handle: item.handle,
				},
			})}
	>
		<span class="avatar"><img src={item.avatar} alt="" /></span>
		<span class="label">{item.label}</span>
	</a>
</article>

<style>
	.card {
		position: relative;
		width: 300px;
		height: 500px;
		flex-shrink: 0;
		border-radius: 16px;
		overflow: hidden;
		border: 3px solid var(--border);
		box-sizing: border-box;
		background: black;
		color: var(--warm-ivory);
	}
	.card > img,
	video,
	iframe {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		border: 0;
		display: block;
	}
	.badge {
		position: absolute;
		top: 14px;
		left: 14px;
		display: flex;
		align-items: center;
		gap: 10px;
		z-index: 2;
		color: var(--warm-ivory);
		text-decoration: none;
	}
	.avatar {
		width: 34px;
		height: 34px;
		border-radius: 999px;
		border: 1px solid var(--warm-ivory);
		overflow: hidden;
		background: rgba(0, 0, 0, 0.35);
	}
	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		display: block;
	}
	.label {
		font-family: var(--bn-font-mono);
		font-weight: 500;
		font-size: 11px;
		letter-spacing: 1.5px;
		text-transform: uppercase;
		text-shadow: 0 1px 6px rgba(0, 0, 0, 0.45);
	}
	.sound {
		position: absolute;
		top: 14px;
		right: 14px;
		width: 36px;
		height: 36px;
		padding: 0;
		border-radius: 999px;
		border: 1px solid var(--warm-ivory);
		background: rgba(0, 0, 0, 0.35);
		backdrop-filter: blur(6px);
		cursor: pointer;
		display: grid;
		place-items: center;
		z-index: 3;
	}
	.sound svg {
		width: 18px;
		height: 18px;
		fill: none;
		stroke: var(--warm-ivory);
		stroke-width: 2;
		stroke-linecap: round;
		stroke-linejoin: round;
	}
	/* Covers the video below the badge row: tap to pause or play. */
	.tap {
		position: absolute;
		top: 60px;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 0;
		border: 0;
		background: transparent;
		cursor: pointer;
		z-index: 1;
		display: grid;
		place-items: center;
	}
	.play {
		width: 56px;
		height: 56px;
		margin-top: -60px;
		border-radius: 999px;
		border: 1px solid var(--warm-ivory);
		background: rgba(17, 39, 26, 0.45);
		backdrop-filter: blur(6px);
		display: grid;
		place-items: center;
	}
	.play svg {
		width: 22px;
		height: 22px;
		fill: var(--warm-ivory);
		margin-left: 3px;
	}
	.tap:focus-visible,
	.sound:focus-visible,
	.badge:focus-visible {
		outline: 2px solid var(--warm-ivory);
		outline-offset: 3px;
	}
</style>
