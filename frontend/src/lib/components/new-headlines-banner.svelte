<script lang="ts">
  import ArrowUp from "@lucide/svelte/icons/arrow-up";
  import { newHeadlinesLabel } from "$lib/helpers/format";

  let { count, onShow }: { count: number; onShow: () => void } = $props();
</script>

<!-- Zero height, so it takes no room in the page; the pill floats under the header. -->
<div class="pointer-events-none sticky top-30 z-20 flex h-0 justify-center" role="status">
  {#if count > 0}
    <button
      type="button"
      class="btn pointer-events-auto rounded-full shadow-lg shadow-primary/20 btn-sm btn-primary"
      onclick={onShow}
    >
      <ArrowUp class="size-4" aria-hidden="true" />
      {newHeadlinesLabel(count)}
    </button>
  {/if}
</div>
