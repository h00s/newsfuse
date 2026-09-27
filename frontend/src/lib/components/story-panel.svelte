<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import ExternalLink from "@lucide/svelte/icons/external-link";
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import Sparkles from "@lucide/svelte/icons/sparkles";
  import { tick } from "svelte";
  import { toastError } from "$lib/helpers/errors";
  import { scrollBehavior } from "$lib/helpers/motion";
  import { fetchStory, summarizeStory } from "$lib/services/stories";
  import type { Headline } from "$lib/types/headline";
  import type { Story } from "$lib/types/story";

  let { headline }: { headline: Headline } = $props();

  /** A story with more text than this starts collapsed and can be summarized. */
  const LONG_STORY = 650;
  /** A story that isn't back by then is being fetched from the news site. */
  const SLOW_LOAD_MS = 1000;

  let panel = $state<HTMLElement>();
  let story = $state<Story | null>(null);
  let failed = $state(false);
  let slow = $state(false);
  let expanded = $state(false);
  let showSummary = $state(false);
  let summarizing = $state(false);

  const long = $derived(story !== null && story.content.replace(/<[^>]*>/g, "").length > LONG_STORY);

  const action = (active: boolean) => [
    "btn btn-sm btn-soft gap-1.5 px-2 font-medium sm:px-3",
    active ? "btn-primary" : "text-base-content/80",
  ];

  async function load() {
    failed = false;
    slow = false;
    const timer = setTimeout(() => (slow = true), SLOW_LOAD_MS);
    try {
      story = await fetchStory(headline.id);
    } catch (e) {
      failed = true;
      toastError(e, {
        404: "Ovaj članak nije moguće prikazati.",
        429: "Previše otvorenih članaka zaredom. Pokušajte ponovno za minutu.",
        502: "Izvorna stranica trenutno nije dostupna.",
      });
    } finally {
      clearTimeout(timer);
    }
  }

  /** Shows the summary, generating it first if the story has none yet; hides it when shown. */
  async function toggleSummary() {
    if (!story || summarizing) return;
    if (showSummary || story.summary) {
      showSummary = !showSummary;
      return;
    }
    summarizing = true;
    try {
      story = await summarizeStory(story.id);
      showSummary = true;
    } catch (e) {
      toastError(e, {
        429: "Previše sažetaka zaredom. Pokušajte ponovno za minutu.",
        502: "Sažimanje trenutno nije dostupno.",
        504: "Sažimanje je predugo trajalo. Pokušajte ponovno.",
      });
    } finally {
      summarizing = false;
    }
  }

  /** Collapsing a long story that was read to its end would leave the reader far below it, so
   *  bring its headline back into view. */
  async function toggleExpanded() {
    expanded = !expanded;
    if (expanded) return;
    await tick();
    panel?.closest("[data-headline]")?.scrollIntoView({ block: "nearest", behavior: scrollBehavior() });
  }

  // Mounted when the reader opens the headline: that is the action that fetches.
  load();
</script>

<div bind:this={panel} class="rounded-box border border-base-300 bg-base-100/70 p-3 sm:p-4">
  {#if story}
    {#if showSummary && story.summary}
      <section aria-label="Sažetak" class="mb-4 rounded-box border border-primary/30 bg-primary/10 p-3">
        <p class="flex items-center gap-1.5 text-xs font-semibold tracking-wide text-primary uppercase">
          <Sparkles class="size-3.5" aria-hidden="true" /> Sažetak
        </p>
        <p class="mt-1.5 text-sm leading-relaxed whitespace-pre-line">{story.summary}</p>
      </section>
    {/if}

    <!-- Sanitized by the server: paragraphs, emphasis, lists, images and safe links only. -->
    <div
      class={[
        "prose prose-sm max-w-none prose-invert prose-p:my-2 prose-a:text-primary prose-img:my-3 prose-img:rounded-box",
        long && !expanded && "max-h-48 overflow-hidden mask-b-from-50%",
      ]}
    >
      {@html story.content}
    </div>

    <div class={["mt-3 gap-2", long ? "grid grid-cols-3 sm:flex" : "flex"]} data-story-actions>
      {#if long}
        <button type="button" class={action(expanded)} aria-expanded={expanded} onclick={toggleExpanded}>
          <ChevronDown class={["size-4 shrink-0 transition-transform motion-reduce:transition-none", expanded && "rotate-180"]} aria-hidden="true" />
          <span class="sm:hidden">{expanded ? "Manje" : "Više"}</span>
          <span class="hidden sm:inline">{expanded ? "Prikaži manje" : "Prikaži više"}</span>
        </button>
        <button
          type="button"
          class={action(showSummary)}
          aria-expanded={showSummary}
          disabled={summarizing}
          onclick={toggleSummary}
        >
          {#if summarizing}
            <span class="loading loading-xs loading-spinner"></span> Sažimam…
          {:else}
            <Sparkles class="size-4 shrink-0" aria-hidden="true" /> Sažetak
          {/if}
        </button>
      {/if}
      <a class={action(false)} href={headline.url} target="_blank" rel="noopener noreferrer">
        <ExternalLink class="size-4 shrink-0" aria-hidden="true" /> Izvornik
      </a>
    </div>
  {:else if failed}
    <div class="flex flex-wrap items-center gap-2 text-sm text-base-content/70">
      <span>Članak se nije učitao.</span>
      <button type="button" class={action(false)} onclick={load}>
        <RotateCw class="size-4" aria-hidden="true" /> Pokušaj ponovno
      </button>
      <a class={action(false)} href={headline.url} target="_blank" rel="noopener noreferrer">
        <ExternalLink class="size-4" aria-hidden="true" /> Izvornik
      </a>
    </div>
  {:else}
    <div class="flex items-center gap-3 py-1 text-sm text-base-content/70" role="status">
      <span class="loading loading-sm loading-spinner text-primary"></span>
      {slow ? `Dohvaćam članak s ${headline.source.name}…` : "Učitavanje članka…"}
    </div>
  {/if}
</div>
