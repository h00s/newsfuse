<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import ExternalLink from "@lucide/svelte/icons/external-link";
  import SourceAvatar from "$lib/components/source-avatar.svelte";
  import StoryPanel from "$lib/components/story-panel.svelte";
  import { absoluteTime, relativeTime } from "$lib/helpers/format";
  import { clock } from "$lib/stores/clock.svelte";
  import type { Headline } from "$lib/types/headline";

  let { headline, isNew = false }: { headline: Headline; isNew?: boolean } = $props();

  const id = $props.id();
  let open = $state(false);

  /** A plain click on a headline whose story can be fetched opens it here. A modified click (new
   *  tab or window) and a headline that can't be fetched go to the news site. */
  function onTitleClick(e: MouseEvent) {
    if (!headline.source.isScrapable || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;
    e.preventDefault();
    open = !open;
  }
</script>

<li
  class="list-row scroll-mt-26 gap-3 rounded-none px-3 py-3 [--list-grid-cols:auto_minmax(0,1fr)_auto] focus-within:bg-base-300/40 sm:px-4"
  data-headline
>
  <SourceAvatar source={headline.source} class="mt-0.5 size-8" />

  <div class="min-w-0">
    <a
      href={headline.url}
      target="_blank"
      rel="noopener noreferrer"
      class={[
        "text-[0.9375rem] leading-snug text-pretty transition-colors hover:text-primary focus-visible:text-primary",
        isNew ? "font-semibold text-base-content" : "text-base-content/85",
      ]}
      onclick={onTitleClick}
      data-headline-link
    >
      {headline.title}
    </a>
    <div class="mt-1 flex flex-wrap items-center gap-x-1.5 gap-y-1 text-xs text-base-content/60">
      {#if isNew}
        <span class="badge badge-xs font-semibold badge-primary">Novo</span>
      {/if}
      <span>{headline.source.name}</span>
      <span aria-hidden="true">·</span>
      <time datetime={headline.publishedAt} title={absoluteTime(headline.publishedAt)}>
        {relativeTime(headline.publishedAt, clock.now)}
      </time>
    </div>
  </div>

  {#if headline.source.isScrapable}
    <button
      type="button"
      class="btn btn-square btn-ghost btn-sm text-base-content/60 hover:text-primary"
      aria-label={open ? "Sakrij članak" : "Prikaži članak"}
      aria-expanded={open}
      aria-controls="{id}-story"
      onclick={() => (open = !open)}
      data-story-toggle
    >
      <ChevronDown class={["size-5 transition-transform motion-reduce:transition-none", open && "rotate-180"]} />
    </button>
  {:else}
    <a
      href={headline.url}
      target="_blank"
      rel="noopener noreferrer"
      class="btn btn-square btn-ghost btn-sm text-base-content/60 hover:text-primary"
      aria-label="Otvori članak na izvornoj stranici"
    >
      <ExternalLink class="size-4" />
    </a>
  {/if}

  {#if open}
    <div id="{id}-story" class="list-col-wrap col-span-full">
      <StoryPanel {headline} />
    </div>
  {/if}
</li>
