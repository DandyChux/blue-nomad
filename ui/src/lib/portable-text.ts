import { toHTML, type PortableTextComponents } from "@portabletext/to-html";

/**
 * Custom Portable Text components that mirror the old Next.js renderers.
 * Each function receives context and returns an HTML string.
 */
const components: PortableTextComponents = {
	types: {
		image: ({ value }) => {
			const url = value?.asset?.url ?? "";
			const alt = value?.alt ?? "";
			const width = value?.width ?? 800;
			const height = value?.height;
			const loading = value?.loading ?? "lazy";
			const className = value?.className ?? "rounded-lg my-8 mx-auto";

			return `<img src="${url}" alt="${alt}" width="${width}" height="${height}" loading="${loading}" class="${className}" />`;
		},
	},

	block: {
		normal: ({ children }) => `<p class="mb-4 font-medium">${children}</p>`,
		blockquote: ({ children }) =>
			`<blockquote class="mb-4">${children}</blockquote>`,
		h1: ({ children }) =>
			`<h1 class="mb-4 font-extrabold">${children}</h1>`,
		h2: ({ children }) =>
			`<h2 class="mb-4 font-extrabold">${children}</h2>`,
		h3: ({ children }) =>
			`<h3 class="mb-4 font-extrabold">${children}</h3>`,
		h4: ({ children }) =>
			`<h4 class="mb-4 font-extrabold">${children}</h4>`,
		h5: ({ children }) =>
			`<h5 class="mb-4 font-extrabold">${children}</h5>`,
		h6: ({ children }) =>
			`<h6 class="mb-4 font-extrabold">${children}</h6>`,
	},

	list: {
		bullet: ({ children }) =>
			`<ul class="my-5 ml-10 list-disc tracking-[0.005em] text-base">${children}</ul>`,
		number: ({ children }) =>
			`<ol class="my-5 ml-10 list-decimal tracking-[0.005em] text-base">${children}</ol>`,
	},

	listItem: {
		bullet: ({ children }) =>
			`<li class="mb-2 list-disc list-inside">${children}</li>`,
		number: ({ children }) =>
			`<li class="mb-2 list-decimal list-inside">${children}</li>`,
	},

	marks: {
		strong: ({ children }) => `<strong>${children}</strong>`,
		em: ({ children }) => `<em>${children}</em>`,
		link: ({ children, value }) => {
			const href = value?.href ?? "#";
			const isExternal =
				href.startsWith("http://") || href.startsWith("https://");
			const rel = isExternal ? "noreferrer noopener" : undefined;
			const target = isExternal ? "_blank" : undefined;
			return `<a href="${href}"${target ? ` target="${target}"` : ""}${rel ? ` rel="${rel}"` : ""}>${children}</a>`;
		},
		footnote: ({ children, value }) => {
			const key = value?._key ?? "";
			return `<span id="src-${key}">${children}<sup><a href="#note-${key}">?</a></sup></span>`;
		},
	},

	hardBreak: () => `<br class="mb-4" />`,
};

/**
 * Renders a Portable Text array to an HTML string.
 * Use with {@html renderPortableText(body)} in Svelte templates.
 */
export function renderPortableText(value: any[]): string {
	if (!value) return "";
	return toHTML(value, { components });
}
