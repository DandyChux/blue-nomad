import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function roundTo(num: number, decimals: number): number {
	const multiplier = Math.pow(10, decimals);
	return Math.round(num * multiplier) / multiplier;
}

/**
 * Convert a hex string to Uint8Array
 */
function hexToBytes(hex: string): Uint8Array<ArrayBuffer> {
	const bytes = new Uint8Array(hex.length / 2);
	for (let i = 0; i < hex.length; i += 2) {
		bytes[i] = parseInt(hex.substring(i, i + 2), 16);
	}
	return bytes;
}

/**
 * Base64URL encode a buffer (RFC 4648)
 */
function base64UrlEncode(buffer: ArrayBuffer): string {
	const bytes = new Uint8Array(buffer);
	let binary = "";
	for (let i = 0; i < bytes.byteLength; i++) {
		binary += String.fromCharCode(bytes[i]);
	}
	return btoa(binary)
		.replace(/=/g, "")
		.replace(/\+/g, "-")
		.replace(/\//g, "_");
}

/**
 * Generate optimized image URL via Imgproxy
 */
export async function getOptimizedImageUrl(
	src: string,
	options?: {
		width?: number | "auto";
		format?: "webp" | "avif";
		quality?: number;
	},
): Promise<string> {
	const IMGPROXY_URL = "https://img.ceokoroji.dev";

	// If no width is provided, default to 0 (Imgproxy keeps original width)
	const width = options?.width === "auto" ? 0 : options?.width || 0;
	const format = options?.format || "webp";
	const quality = options?.quality || 85;

	const keyHex = import.meta.env.VITE_IMGPROXY_KEY;
	const saltHex = import.meta.env.VITE_IMGPROXY_SALT;

	// Imgproxy processing options:
	// rs:fill:WIDTH:HEIGHT / q:QUALITY / f:FORMAT
	const processingOptions = `rs:fill:${width}:0/q:${quality}/f:${format}`;
	const path = `/${processingOptions}/plain/${src}`;

	if (!keyHex || !saltHex) {
		return `${IMGPROXY_URL}/insecure/${path}`;
	}

	// Convert hex strings to byte arrays
	const keyBytes = hexToBytes(keyHex);
	const saltBytes = hexToBytes(saltHex);
	const pathBytes = new TextEncoder().encode(path);

	// Combine salt + path
	const saltedPath = new Uint8Array(saltBytes.length + pathBytes.length);
	saltedPath.set(saltBytes, 0);
	saltedPath.set(pathBytes, saltBytes.length);

	// Import key into Web Crypto
	const cryptoKey = await crypto.subtle.importKey(
		"raw",
		keyBytes,
		{ name: "HMAC", hash: "SHA-256" },
		false,
		["sign"],
	);

	// Calculate HMAC-SHA256 signature
	const signatureBuffer = await crypto.subtle.sign(
		"HMAC",
		cryptoKey,
		saltedPath,
	);
	const signature = base64UrlEncode(signatureBuffer);

	return `${IMGPROXY_URL}/${signature}${path}`;
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
		.map(async (width) => {
			const url = await getOptimizedImageUrl(src, {
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
