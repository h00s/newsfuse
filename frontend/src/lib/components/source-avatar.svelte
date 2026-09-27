<script lang="ts">
  import { sourceInitials, sourceLogo } from "$lib/helpers/sources";
  import type { Source } from "$lib/types/source";

  let { source, class: className = "size-8" }: { source: Pick<Source, "name">; class?: string } = $props();

  let failed = $state(false);
</script>

<!-- Decorative: the source's name is always written next to it. -->
{#if failed}
  <span
    class={["grid shrink-0 place-items-center rounded-lg bg-neutral text-[0.625rem] font-semibold text-neutral-content", className]}
    aria-hidden="true"
  >
    {sourceInitials(source.name)}
  </span>
{:else}
  <img
    src={sourceLogo(source)}
    alt=""
    width="32"
    height="32"
    loading="lazy"
    decoding="async"
    class={["shrink-0 rounded-lg bg-base-300 object-cover", className]}
    onerror={() => (failed = true)}
  />
{/if}
