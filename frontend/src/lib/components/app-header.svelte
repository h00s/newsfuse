<script lang="ts">
  import Code from "@lucide/svelte/icons/code";
  import Keyboard from "@lucide/svelte/icons/keyboard";
  import Search from "@lucide/svelte/icons/search";
  import TopicTabs from "$lib/components/topic-tabs.svelte";
  import { routes } from "$lib/helpers/routes";
  import type { Topic } from "$lib/types/topic";

  type Props = { topics: Topic[]; activeTopicId: number | null; onShortcuts: () => void };
  let { topics, activeTopicId, onShortcuts }: Props = $props();
</script>

<header class="sticky top-0 z-30 border-b border-base-300 bg-base-100/85 backdrop-blur-md">
  <div class="mx-auto max-w-3xl px-3 sm:px-4">
    <div class="flex h-14 items-center gap-2">
      <a href={routes.home} class="mr-auto text-xl font-bold tracking-tight">
        Fused<span class="text-primary">News</span>
      </a>
      <a href={routes.search} class="btn gap-2 btn-ghost btn-sm" aria-label="Pretraga">
        <Search class="size-4" aria-hidden="true" />
        <span class="hidden text-base-content/70 sm:inline">Pretraga</span>
        <kbd class="kbd hidden kbd-sm sm:inline-flex">/</kbd>
      </a>
      <button
        type="button"
        class="btn hidden btn-square btn-ghost btn-sm sm:inline-flex"
        aria-label="Prečaci tipkovnice"
        onclick={onShortcuts}
      >
        <Keyboard class="size-4" />
      </button>
      <a
        href="https://github.com/h00s/newsfuse"
        target="_blank"
        rel="noopener noreferrer"
        class="btn btn-square btn-ghost btn-sm"
        aria-label="Izvorni kod na GitHubu"
      >
        <Code class="size-4" />
      </a>
    </div>
    <TopicTabs {topics} {activeTopicId} />
  </div>
</header>
