<script lang="ts">
  import type { Component, Snippet } from "svelte";

  type Props = { icon: Component<{ class?: string }>; title: string; children?: Snippet };
  let { icon: Icon, title, children }: Props = $props();
</script>

<div class="flex flex-col items-center gap-3 rounded-box border border-dashed border-base-300 px-6 py-16 text-center">
  <Icon class="size-10 text-base-content/30" />
  <p class="font-medium">{title}</p>
  {#if children}
    <div class="max-w-sm text-sm text-base-content/60">{@render children()}</div>
  {/if}
</div>
