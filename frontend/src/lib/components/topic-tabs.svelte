<script lang="ts">
  import { topicPath } from "$lib/helpers/routes";
  import { newCounts } from "$lib/stores/new-counts.svelte";
  import type { Topic } from "$lib/types/topic";

  let { topics, activeTopicId }: { topics: Topic[]; activeTopicId: number | null } = $props();
</script>

<!-- Phones share the width between the tabs, and a count is a bubble on the name's corner, so a
     badge never widens a tab. Wider screens show it inline. -->
<nav aria-label="Teme" class="-mx-3 overflow-x-auto overflow-y-hidden px-1 sm:mx-0 sm:px-0">
  <div class="tabs flex-nowrap tabs-border">
    {#each topics as topic (topic.id)}
      {@const active = topic.id === activeTopicId}
      <!-- The open topic shows its news in the banner instead. -->
      {@const count = active ? 0 : newCounts.count(topic.id)}
      <a
        href={topicPath(topic.id)}
        class={[
          "tab flex-1 px-2 font-medium whitespace-nowrap sm:flex-none sm:px-4",
          active ? "tab-active text-primary" : "text-base-content/70",
        ]}
        aria-current={active ? "page" : undefined}
      >
        <span class="relative inline-flex items-center gap-2">
          {topic.name}
          {#if count > 0}
            <span
              class="absolute -top-2 left-full ml-0.5 badge h-4 min-w-4 px-1 text-[0.625rem] leading-none font-semibold badge-success sm:static sm:ml-0 sm:h-5 sm:px-1.5 sm:text-xs"
              aria-label="{count} novih"
            >
              {count > 99 ? "99+" : count}
            </span>
          {/if}
        </span>
      </a>
    {/each}
  </div>
</nav>
