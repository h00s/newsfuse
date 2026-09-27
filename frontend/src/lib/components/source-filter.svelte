<script lang="ts">
  import SourceAvatar from "$lib/components/source-avatar.svelte";
  import { topicPath } from "$lib/helpers/routes";
  import type { Source } from "$lib/types/source";

  type Props = { topicId: number; sources: Source[]; activeSourceId: number | null };
  let { topicId, sources, activeSourceId }: Props = $props();

  const chip = (active: boolean) => [
    "btn btn-sm shrink-0 gap-2 rounded-full font-medium whitespace-nowrap",
    active ? "btn-primary" : "border-base-300 bg-base-200 text-base-content/80 btn-ghost hover:text-base-content",
  ];
</script>

<!-- Links rather than toggles: the filter lives in the URL, so it survives a reload and can be
     shared, and switching it replaces the history entry instead of piling them up. -->
{#if sources.length > 1}
  <nav aria-label="Izvori" class="-mx-3 overflow-x-auto px-3 sm:mx-0 sm:px-0">
    <ul class="flex gap-2 py-1">
      <li>
        <a
          href={topicPath(topicId)}
          class={chip(activeSourceId === null)}
          aria-current={activeSourceId === null ? "page" : undefined}
          data-sveltekit-replacestate
          data-sveltekit-noscroll
          data-sveltekit-keepfocus
        >
          Svi izvori
        </a>
      </li>
      {#each sources as source (source.id)}
        <li>
          <a
            href={topicPath(topicId, source.id)}
            class={chip(activeSourceId === source.id)}
            aria-current={activeSourceId === source.id ? "page" : undefined}
            data-sveltekit-replacestate
            data-sveltekit-noscroll
            data-sveltekit-keepfocus
          >
            <SourceAvatar {source} class="size-5 rounded-md" />
            {source.name}
          </a>
        </li>
      {/each}
    </ul>
  </nav>
{/if}
