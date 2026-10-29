import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { generateImageUrl } from "@imgproxy/imgproxy-node";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function roundTo(num: number, decimals: number): number {
	const multiplier = Math.pow(10, decimals);
	return Math.round(num * multiplier) / multiplier;
}

/**
 * Generate optimized image URL via Imgproxy
 */
export function getOptimizedImageUrl(
	src: string,
	options?: {
		width?: number | "auto";
		format?: "webp" | "avif";
		quality?: number;
	},
): string {
	const width = options?.width === "auto" ? 0 : options?.width || 0;
	const format = options?.format || "webp";
	const quality = options?.quality || 85;

	return generateImageUrl({
		endpoint: "https://img.ceokoroji.dev",
		url: {
			value: src,
			displayAs: "plain",
		},
		key: import.meta.env.IMGPROXY_KEY,
		salt: import.meta.env.IMGPROXY_SALT,
		options: {
			resizing_type: "fill",
			width: width,
			height: 0,
			quality: quality,
			format: format,
		},
	});
}

/**
 * Generate a srcset string for responsive images
 */
export function generateSrcSet(
	src: string,
	widths: number[] = [400, 800, 1200, 1600],
	format: "webp" | "avif" = "webp",
	quality: number = 85,
): string {
	return widths
		.map((width) => {
			const url = getOptimizedImageUrl(src, {
				width,
				format,
				quality,
			});
			return `${url} ${width}w`; // The 'w' tells the browser the true pixel width
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
