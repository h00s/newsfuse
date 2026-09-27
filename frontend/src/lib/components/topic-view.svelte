<script lang="ts">
  import Inbox from "@lucide/svelte/icons/inbox";
  import { refreshAll } from "$app/navigation";
  import { onMount, untrack } from "svelte";
  import EmptyState from "$lib/components/empty-state.svelte";
  import HeadlineList from "$lib/components/headline-list.svelte";
  import NewHeadlinesBanner from "$lib/components/new-headlines-banner.svelte";
  import SourceFilter from "$lib/components/source-filter.svelte";
  import { HeadlineFeed } from "$lib/feeds/headline-feed.svelte";
  import { fetchHeadlines } from "$lib/services/headlines";
  import { newCounts } from "$lib/stores/new-counts.svelte";
  import { reading } from "$lib/stores/reading.svelte";
  import type { Headline } from "$lib/types/headline";
  import type { Source } from "$lib/types/source";
  import type { Topic } from "$lib/types/topic";

  type Props = { topic: Topic; sources: Source[]; sourceId: number | null; headlines: Headline[] };
  let { topic, sources, sourceId, headlines }: Props = $props();

  // What was new when the topic opened stays marked while it is open. A first visit marks nothing.
  let seenBefore = $state(untrack(() => reading.lastSeen(topic.id) ?? Date.now()));

  // Every load (a source picked, the list refreshed) starts a fresh feed.
  const feed = $derived(
    new HeadlineFeed(headlines, (beforeId) => fetchHeadlines({ topicId: topic.id, sourceId, beforeId })),
  );

  onMount(() => {
    reading.markSeen(topic.id);
    newCounts.seen(topic.id);
  });

  async function showNew() {
    seenBefore = reading.lastSeen(topic.id) ?? Date.now();
    reading.markSeen(topic.id);
    newCounts.seen(topic.id);
    await refreshAll();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }
</script>

<div class="space-y-4">
  <SourceFilter topicId={topic.id} {sources} activeSourceId={sourceId} />
  <NewHeadlinesBanner count={newCounts.count(topic.id)} onShow={showNew} />
  <HeadlineList {feed} {seenBefore}>
    {#snippet empty()}
      <EmptyState icon={Inbox} title="Još nema vijesti">
        Novi naslovi stižu svakih nekoliko minuta.
      </EmptyState>
    {/snippet}
  </HeadlineList>
</div>
