<script lang="ts">
  import { topicPath } from "$lib/helpers/routes";
  import { newCounts } from "$lib/stores/new-counts.svelte";
  import type { Topic } from "$lib/types/topic";

  let { topics, activeTopicId }: { topics: Topic[]; activeTopicId: number | null } = $props();
</script>

<nav aria-label="Teme" class="-mx-3 overflow-x-auto px-3 sm:mx-0 sm:px-0">
  <div class="tabs flex-nowrap tabs-border">
    {#each topics as topic (topic.id)}
      {@const active = topic.id === activeTopicId}
      <!-- The open topic shows its news in the banner instead. -->
      {@const count = active ? 0 : newCounts.count(topic.id)}
      <a
        href={topicPath(topic.id)}
        class={["tab gap-2 font-medium whitespace-nowrap", active ? "tab-active text-primary" : "text-base-content/70"]}
        aria-current={active ? "page" : undefined}
      >
        {topic.name}
        {#if count > 0}
          <span class="badge badge-sm font-semibold badge-success" aria-label="{count} novih">
            {count > 99 ? "99+" : count}
          </span>
        {/if}
      </a>
    {/each}
  </div>
</nav>
