<script lang="ts">
	import { trackEvent } from "$lib/analytics.svelte";
	import { cn } from "$lib/utils";
	import { Button, buttonVariants } from "./ui/button";

	type Props = {
		variant?:
			| "default"
			| "outline"
			| "secondary"
			| "ghost"
			| "destructive"
			| "link";
		size?:
			| "default"
			| "xs"
			| "sm"
			| "lg"
			| "xl"
			| "icon"
			| "icon-xs"
			| "icon-sm"
			| "icon-lg";
		href: string;
		label: string;
		event: string;
		location: string;
		external?: boolean;
		class?: string;
	};

	let {
		href,
		label,
		variant = "outline",
		event = "book_click",
		location = "",
		external = false,
		class: className,
	}: Props = $props();
</script>

<Button
	variant={"link"}
	{href}
	rel={external ? "noopener noreferrer" : undefined}
	data-track={event}
	onclick={() => trackEvent(event, { props: { label, location } })}
	class={buttonVariants({
		variant: variant,
		size: "xl",
		class: cn("rounded-full px-8 h-10 uppercase", className),
	})}
>
	{label}
</Button>
