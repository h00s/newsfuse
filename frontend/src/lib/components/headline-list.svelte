<script lang="ts">
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import type { Snippet } from "svelte";
  import HeadlineItem from "$lib/components/headline-item.svelte";
  import { inView } from "$lib/attachments/in-view";
  import type { HeadlineFeed } from "$lib/feeds/headline-feed.svelte";
  import { toastError } from "$lib/helpers/errors";

  type Props = {
    feed: HeadlineFeed;
    /** Headlines published after this (epoch ms) are marked new. */
    seenBefore?: number;
    /** Shown instead of the list when there are no headlines. */
    empty: Snippet;
  };
  let { feed, seenBefore = Number.POSITIVE_INFINITY, empty }: Props = $props();

  async function loadMore() {
    try {
      await feed.loadMore();
    } catch (e) {
      toastError(e);
    }
  }
</script>

{#if feed.items.length === 0}
  {@render empty()}
{:else}
  <ul class="list rounded-box border border-base-300 bg-base-200 shadow-sm" aria-label="Naslovi">
    {#each feed.items as headline (headline.id)}
      <HeadlineItem {headline} isNew={Date.parse(headline.publishedAt) > seenBefore} />
    {/each}
  </ul>

  <div class="flex min-h-20 items-center justify-center py-4">
    {#if feed.loading}
      <span class="loading loading-md loading-dots text-primary" role="status" aria-label="Učitavanje"></span>
    {:else if feed.failed}
      <button type="button" class="btn btn-ghost btn-sm" onclick={loadMore}>
        <RotateCw class="size-4" aria-hidden="true" /> Pokušaj ponovno
      </button>
    {:else if !feed.done}
      <!-- Remounts after every page, so a button still in view loads the next one too. -->
      <button type="button" class="btn btn-ghost btn-sm" onclick={loadMore} {@attach inView(loadMore)}>Učitaj još</button>
    {:else}
      <p class="text-sm text-base-content/50">To je sve.</p>
    {/if}
  </div>
{/if}
