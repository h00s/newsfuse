<script lang="ts">
  import ExternalLink from "@lucide/svelte/icons/external-link";
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import Sparkles from "@lucide/svelte/icons/sparkles";
  import { toastError } from "$lib/helpers/errors";
  import { fetchStory, summarizeStory } from "$lib/services/stories";
  import type { Headline } from "$lib/types/headline";
  import type { Story } from "$lib/types/story";

  let { headline }: { headline: Headline } = $props();

  /** A story with more text than this starts collapsed and can be summarized. */
  const LONG_STORY = 650;

  let story = $state<Story | null>(null);
  let failed = $state(false);
  let expanded = $state(false);
  let summarizing = $state(false);

  const long = $derived(story !== null && story.content.replace(/<[^>]*>/g, "").length > LONG_STORY);

  async function load() {
    failed = false;
    try {
      story = await fetchStory(headline.id);
    } catch (e) {
      failed = true;
      toastError(e, {
        404: "Ovaj članak nije moguće prikazati.",
        429: "Previše otvorenih članaka zaredom. Pokušajte ponovno za minutu.",
        502: "Izvorna stranica trenutno nije dostupna.",
      });
    }
  }

  async function summarize() {
    if (!story) return;
    summarizing = true;
    try {
      story = await summarizeStory(story.id);
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

  // Mounted when the reader opens the headline: that is the action that fetches.
  load();
</script>

<div class="rounded-box border border-base-300 bg-base-100/70 p-3 sm:p-4">
  {#if story}
    {#if story.summary}
      <section aria-label="Sažetak" class="mb-4 rounded-box border border-accent/30 bg-accent/10 p-3">
        <p class="flex items-center gap-1.5 text-xs font-semibold tracking-wide text-accent uppercase">
          <Sparkles class="size-3.5" aria-hidden="true" /> Sažetak
        </p>
        <p class="mt-1.5 text-sm leading-relaxed whitespace-pre-line">{story.summary}</p>
      </section>
    {/if}

    <!-- Sanitized by the server: paragraphs, emphasis, lists and safe links only. -->
    <div
      class={[
        "prose prose-sm max-w-none prose-invert prose-p:my-2 prose-a:text-primary",
        long && !expanded && "max-h-48 overflow-hidden mask-b-from-50%",
      ]}
    >
      {@html story.content}
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-2">
      {#if long}
        <button type="button" class="btn btn-ghost btn-xs" aria-expanded={expanded} onclick={() => (expanded = !expanded)}>
          {expanded ? "Skrati" : "Prikaži cijeli članak"}
        </button>
      {/if}
      {#if long && !story.summary}
        <button type="button" class="btn btn-soft btn-xs btn-accent" disabled={summarizing} onclick={summarize}>
          {#if summarizing}
            <span class="loading loading-xs loading-spinner"></span> Sažimam…
          {:else}
            <Sparkles class="size-3.5" aria-hidden="true" /> Sažmi članak
          {/if}
        </button>
      {/if}
      <a class="btn btn-soft btn-xs btn-primary" href={headline.url} target="_blank" rel="noopener noreferrer">
        <ExternalLink class="size-3.5" aria-hidden="true" /> Izvorni članak
      </a>
    </div>
  {:else if failed}
    <div class="flex flex-wrap items-center gap-2 text-sm text-base-content/70">
      <span>Članak se nije učitao.</span>
      <button type="button" class="btn btn-ghost btn-xs" onclick={load}>
        <RotateCw class="size-3.5" aria-hidden="true" /> Pokušaj ponovno
      </button>
      <a class="btn btn-soft btn-xs btn-primary" href={headline.url} target="_blank" rel="noopener noreferrer">
        <ExternalLink class="size-3.5" aria-hidden="true" /> Izvorni članak
      </a>
    </div>
  {:else}
    <div class="space-y-2.5" role="status" aria-label="Učitavanje članka">
      <div class="skeleton h-3.5 w-full"></div>
      <div class="skeleton h-3.5 w-11/12"></div>
      <div class="skeleton h-3.5 w-4/5"></div>
    </div>
  {/if}
</div>
