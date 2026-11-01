<script lang="ts">
	import { slide } from "svelte/transition";
	import { Button } from "./ui/button";
	import BookButton from "./book-button.svelte";

	let active = $state(0); // selected topic pill
	let pinned = $state(-1); // question opened by a click or tap
	let hovered = $state(-1); // question under the mouse (desktop only)
	let canHover = $state(false);

	$effect(() => {
		const query = window.matchMedia("(hover: hover) and (pointer: fine)");
		canHover = query.matches;
		const update = (e) => (canHover = e.matches);
		query.addEventListener("change", update);
		return () => query.removeEventListener("change", update);
	});

	const topic = $derived(faq.topics[active]);

	function selectTopic(i) {
		active = i;
		pinned = -1;
		hovered = -1;
	}
	// Desktop: hovering a question shows its answer. Phone: tapping opens and closes it.
	const isOpen = (j) => pinned === j || hovered === j;

	const faq = {
		title: "Skin Health Guide",
		intro: "Personalized skin care starts with understanding your skin. Explore our approach to common skin concerns, treatments and membership.",
		ready: "Ready to begin?",
		cta: "Book Facial ST",
		topics: [
			{
				name: "Breakouts",
				items: [
					[
						"Can facials help with acne and recurring breakouts?",
						"Yes. Facial Skin Therapy supports clients working on breakouts over time through personalized treatments and home-care guidance. We adapt each session to your skin’s needs, tolerance and progress.",
					],
					[
						"Will my facial include extractions?",
						"Where appropriate, we perform extractions to address congestion. We never force a blemish that is not ready to release, and we adjust our approach to minimize unnecessary irritation.",
					],
					[
						"Will I receive facial massage?",
						"Facial massage may be incorporated into your treatment, depending on your skin’s needs and the approach selected for your session.",
					],
				],
			},
			{
				name: "Dark Spots",
				items: [
					[
						"Can facials help with dark spots and uneven pigmentation?",
						"Personalized treatments and consistent home care can help improve the appearance of dark spots over time. We take care to minimize unnecessary irritation, which can contribute to further pigmentation.",
					],
					[
						"How long does it take to see results?",
						"This depends on the cause and depth of the pigmentation, your skin’s response and the consistency of your home care. We focus on a considered, gradual approach rather than promising overnight results.",
					],
				],
			},
			{
				name: "Skin of Color",
				items: [
					[
						"Do you have particular expertise in skin of color?",
						"Yes. We welcome clients of all skin tones, with particular expertise in concerns such as breakouts and post-inflammatory hyperpigmentation. These concerns can be especially persistent in melanin-rich skin, so we take care when selecting products, exfoliation and treatment techniques.",
					],
					[
						"What do you mean by skin of color?",
						"Skin of color is a broad term generally used to describe people with more melanin-rich skin. This includes many people of African, South Asian, Southeast Asian, East Asian, Middle Eastern, Latino and Indigenous heritage. It is not limited to one ethnicity or skin tone, and people within the same community can have very different complexions.\n\nWe may also consider the Fitzpatrick scale, which classifies skin into six types based on its response to UV exposure. Types IV to VI are more commonly associated with darker skin that tans more readily and burns less easily. However, the scale does not fully capture skin tone, pigmentation tendencies or how your skin may respond to treatment. We tailor our approach to your individual skin rather than relying on ethnicity or a classification alone.",
					],
					[
						"Why does this matter when treating my skin?",
						"Inflammation and irritation can contribute to persistent dark marks, particularly in skin prone to post-inflammatory hyperpigmentation. We take this into account when selecting products and treatment techniques, with the aim of supporting your skin without causing unnecessary irritation.",
					],
				],
			},
			{
				name: "Sensitive Skin",
				items: [
					[
						"Are your treatments suitable for sensitive skin?",
						"Yes. We tailor each treatment to your skin’s current condition and tolerance, selecting products and techniques with care. Our priority is to support your skin while minimizing unnecessary irritation.",
					],
					[
						"What if my skin reacts easily to products or treatments?",
						"Let us know about previous reactions and any current sensitivity before your treatment. We will take these into account and adjust our approach accordingly. We do not believe in pushing your skin beyond what it can comfortably tolerate.",
					],
				],
			},
			{
				name: "Fine Lines + Skin Longevity",
				items: [
					[
						"Are facials suitable if I am focused on fine lines and skin longevity?",
						"Absolutely. Regular treatments can support hydration, skin maintenance and a consistent approach to skin health over time. We tailor each session to your skin’s needs and evolving goals.",
					],
					[
						"Can facials help with fine lines?",
						"Personalized treatments can help improve the appearance of fine lines by supporting hydration and overall skin condition. Results vary, and we focus on thoughtful, consistent care rather than quick fixes.",
					],
					[
						"Can you help me adapt my routine as my skin changes?",
						"Yes. Your skin’s needs can shift with age, seasons and environmental factors. We provide personalized guidance to help you make thoughtful adjustments as needed.",
					],
				],
			},
			{
				name: "Treatments + Memberships",
				items: [
					[
						"Do I need a specific skin concern to book a facial?",
						"Not at all. Some clients come to work on concerns such as breakouts or pigmentation over time. Others visit monthly for skin maintenance, hydration, wellbeing and skin longevity. Your treatment is tailored to your skin and your goals.",
					],
					[
						"What is included in the Blue Nomad membership?",
						"Membership includes one monthly 60-minute Facial Skin Therapy treatment, personalized skin-care guidance, preferred rates on additional treatments and products, and free shipping on product orders.",
					],
					[
						"Do I need to be a member to book?",
						"No. You are welcome to book a single Facial Skin Therapy session or become a member for ongoing monthly care.",
					],
					[
						"How do I book?",
						"Book your Facial Skin Therapy session online. If you are looking for consistent, personalized care, explore our membership for monthly treatments and ongoing skin-care guidance.",
					],
				],
			},
		],
	};
</script>

<section class="bn-section" id="faq">
	<div class="bn-inner faq">
		<header>
			<h2 class="bn-title">{faq.title}</h2>
			<p class="bn-body bn-brief intro">{faq.intro}</p>
		</header>

		<div class="pills" role="tablist" aria-label="FAQ topics">
			{#each faq.topics as t, i}
				<button
					class="pill"
					class:on={i === active}
					role="tab"
					aria-selected={i === active}
					onclick={() => selectTopic(i)}
				>
					{t.name}
				</button>
			{/each}
		</div>

		<div class="list" role="tabpanel">
			{#each topic.items as [question, answer], j (topic.name + j)}
				<div
					class="item"
					role="presentation"
					onmouseenter={() => canHover && (hovered = j)}
					onmouseleave={() => (hovered = -1)}
				>
					<button
						class="q"
						aria-expanded={isOpen(j)}
						onclick={() => (pinned = pinned === j ? -1 : j)}
					>
						<span>{question}</span>
						<span class="icon" aria-hidden="true"
							>{isOpen(j) ? "−" : "+"}</span
						>
					</button>
					{#if isOpen(j)}
						<div class="a" transition:slide={{ duration: 250 }}>
							{#each answer.split("\n\n") as paragraph}
								<p class="bn-body">{paragraph}</p>
							{/each}
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<div class="ready">
			<h3>{faq.ready}</h3>
			<BookButton
				href="https://bluenomadworld.com/booking/TWYSCIC46EIMS3SD2A6UMJ5H"
				label={faq.cta}
				event="book_facial_st_click"
				location="faq"
			/>
		</div>
	</div>
</section>

<style>
	.faq {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}
	header {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}
	.bn-title {
		font-size: 48px;
	}
	.intro {
		max-width: 640px;
	}
	.pills {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}
	.pill {
		padding: 9px 18px;
		border-radius: 999px;
		border: 1px solid black;
		background: transparent;
		color: black;
		font-family: var(--font-source-code-pro);
		font-weight: 500;
		font-size: 13px;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		cursor: pointer;
	}
	.pill.on {
		background: black;
		color: var(--warm-ivory);
	}
	.list {
		border-top: 1px solid black;
	}
	.item {
		border-bottom: 1px solid black;
	}
	.q {
		width: 100%;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 24px;
		padding: 20px 0;
		background: none;
		border: 0;
		text-align: left;
		cursor: pointer;
		color: black;
		font-family: var(--font-source-code-pro);
		font-weight: 500;
		font-size: 15px;
		letter-spacing: 0.04em;
		line-height: 1.4;
		text-transform: uppercase;
	}
	.icon {
		font-size: 20px;
		line-height: 1;
	}
	.a {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 0 48px 20px 0;
	}
	.a .bn-body {
		font-size: 15px;
	}
	.ready {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: center;
		gap: 20px 40px;
		padding-top: 8px;
	}
	.ready h3 {
		margin: 0;
		font-family: var(--font-harmony);
		font-weight: 500;
		font-size: 28px;
		line-height: 1;
		text-transform: uppercase;
	}
	@media (max-width: 1280px) {
		.bn-title {
			font-size: 40px;
		}
		.ready {
			flex-direction: column;
		}
		.a {
			padding-right: 0;
		}
	}
</style>
