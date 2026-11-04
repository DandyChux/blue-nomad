<script lang="ts">
	import { onMount } from "svelte";
	import Picture from "./picture.svelte";

	export type MovingCardProps = {
		name: string;
		image: string;
		link: string;
		quote?: string;
	};

	let {
		items,
		direction = "left",
		speed = "fast",
		pauseOnHover = false,
		class: className = "",
	}: {
		items: MovingCardProps[];
		direction?: "left" | "right";
		speed?: "fast" | "normal" | "slow";
		pauseOnHover?: boolean;
		class?: string;
	} = $props();

	let containerEl: HTMLDivElement | undefined = $state();
	let scrollerEl: HTMLUListElement | undefined = $state();
	let started = $state(false);

	function getAnimationDuration(speed: string): string {
		switch (speed) {
			case "fast":
				return "20s";
			case "normal":
				return "40s";
			default:
				return "80s";
		}
	}

	function getAnimationDirection(direction: string): string {
		return direction === "left" ? "forwards" : "reverse";
	}

	onMount(() => {
		if (containerEl && scrollerEl) {
			// Duplicate items for infinite scroll effect
			const scrollerContent = Array.from(scrollerEl.children);
			scrollerContent.forEach((item) => {
				const duplicatedItem = item.cloneNode(true);
				scrollerEl!.appendChild(duplicatedItem);
			});

			containerEl.style.setProperty(
				"--animation-direction",
				getAnimationDirection(direction),
			);
			containerEl.style.setProperty(
				"--animation-duration",
				getAnimationDuration(speed),
			);
			started = true;
		}
	});
</script>

<div
	bind:this={containerEl}
	class={[
		"scroller relative z-20 max-w-7xl overflow-hidden mask-[linear-gradient(to_right,transparent,white_20%,white_80%,transparent)]",
		className,
	]}
>
	<ul
		bind:this={scrollerEl}
		class={[
			"flex w-max min-w-full shrink-0 flex-nowrap py-4",
			started && "animate-scroll",
			pauseOnHover && "hover:[animation-play-state:paused]",
		]}
	>
		{#each items as item (item.name)}
			<li
				class="relative md:w-[300px] w-[220px] max-w-full shrink-0 rounded-2xl md:px-8 px-3 md:py-6 py-4"
			>
				<blockquote>
					<div
						aria-hidden="true"
						class="user-select-none pointer-events-none absolute -top-0.5 -left-0.5 -z-1 h-[calc(100%+4px)] w-[calc(100%+4px)]"
					></div>

					<div
						class="relative w-full md:h-40 h-32 mb-4 overflow-hidden rounded-lg flex place-items-center aspect-4/3"
					>
						<Picture
							src={item.image}
							alt="Image for {item.name}"
							class="object-cover object-center mx-auto"
							width="125"
							height="100"
							loading="lazy"
						/>
					</div>
				</blockquote>
			</li>
		{/each}
	</ul>
</div>
