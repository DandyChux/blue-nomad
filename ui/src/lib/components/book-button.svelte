<script>
	import { trackEvent } from "$lib/analytics.svelte";
	import { Button, buttonVariants } from "./ui/button";

	let {
		href,
		label,
		variant = "ink",
		event = "book_click",
		location = "",
		external = false,
	} = $props();
</script>

<Button
	variant={"link"}
	{href}
	rel={external ? "noopener noreferrer" : undefined}
	data-track={event}
	onclick={() => trackEvent(event, { props: { label, location } })}
	class={buttonVariants({
		variant: "outline",
		size: "xl",
		class: "rounded-full px-8 h-10 uppercase",
	})}
>
	{label}
</Button>
