<script lang="ts">
	import { Star } from "@lucide/svelte";

	interface Props {
		value?: number;
		max?: number;
		class?: string;
	}

	let { value = 5, max = 5, class: className = "" }: Props = $props();

	const normalizedMax = $derived(Math.max(0, max));
	const normalizedValue = $derived(
		Math.max(0, Math.min(normalizedMax, Math.round(value))),
	);
</script>

<div
	class={["inline-flex items-center gap-1 text-black my-2", className]}
	aria-label={`${normalizedValue} out of ${normalizedMax} stars`}
>
	{#each Array.from({ length: normalizedMax }) as _, index (index)}
		<span
			aria-hidden="true"
			class={index < normalizedValue
				? "opacity-100 text-2xl rounded-lg"
				: "opacity-25"}
		>
			<!-- ★ -->
			<Star size={24} class="fill-inherit" />
		</span>
	{/each}
</div>
