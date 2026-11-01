import { browser } from "$app/environment";

type PlausibleEventOptions = {
	props?: Record<string, string | number | boolean>;
	callback?: () => void;
};

/**
 * Track a custom event with Plausible Analytics
 * @param eventName - The name of the event to track
 * @param options - Optional props and callback
 */
export function trackEvent(
	eventName: string,
	options?: PlausibleEventOptions,
): void {
	if (!browser) return;

	if (typeof window.plausible === "function") {
		window.plausible(eventName, options);
	}
	window.dataLayer?.push({ event: eventName, ...options?.props });
}

// Type declaration for the plausible function on window
declare global {
	interface Window {
		plausible?: (
			eventName: string,
			options?: {
				props?: Record<string, string | number | boolean>;
				callback?: () => void;
			},
		) => void;
	}
}
