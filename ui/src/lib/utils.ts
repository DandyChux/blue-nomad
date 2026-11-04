import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function roundTo(num: number, decimals: number): number {
	const multiplier = Math.pow(10, decimals);
	return Math.round(num * multiplier) / multiplier;
}

function imageURL(
	src: string,
	width: number,
	format: "webp" | "avif",
	quality: number,
): string {
	if (!src.startsWith("https://")) return src;

	const query = new URLSearchParams({
		src,
		format,
		quality: String(quality),
	});
	return `/api/images/${width}?${query}`;
}

/**
 * Generate a srcset string for responsive images
 */
export function generateSrcSet(
	src: string,
	widths: number[] = [200, 400, 600, 800],
	format: "webp" | "avif" = "webp",
	quality: number = 85,
): string {
	return widths
		.map((width) => {
			return `${imageURL(src, width, format, quality)} ${width}w`;
		})
		.join(", ");
}

/**
 * Create LCP preload link for above-fold images
 */
export function createLCPPreload(src: string): HTMLLinkElement | null {
	const link = document.createElement("link");
	link.rel = "preload";
	link.as = "image";
	link.href = src;
	Object.assign(link, {
		crossorigin: "anonymous",
	});
	return link;
}

/**
 * Create a lightweight error handler for performance monitoring
 */
// export function createPerformanceObserver(callback: PerformanceObserverEntryList => void): PerformanceObserver | null {
// 	try {
// 		return new PerformanceObserver((list) => callback(list.getEntries()))
// 	} catch (e) {
// 		return null
// 	}
// }

/**
 * Memoized debounce for performance
 */
export function debounce<T extends (...args: any[]) => void>(
	func: T,
	wait: number,
): (...args: Parameters<T>) => void {
	let timeout: ReturnType<typeof setTimeout>;
	return ((...args: Parameters<T>) => {
		clearTimeout(timeout);
		timeout = setTimeout(() => func(...args), wait);
	}) as T;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any }
	? Omit<T, "children">
	: T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & {
	ref?: U | null;
};
