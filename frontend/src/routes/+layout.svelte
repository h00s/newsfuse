<script lang="ts">
  import "./layout.css";
  import { goto } from "$app/navigation";
  import { navigating, page } from "$app/state";
  import { Toaster } from "svelte-sonner";
  import AppFooter from "$lib/components/app-footer.svelte";
  import AppHeader from "$lib/components/app-header.svelte";
  import ShortcutsDialog from "$lib/components/shortcuts-dialog.svelte";
  import { routes, topicPath } from "$lib/helpers/routes";
  import { isTypingTarget, moveHeadlineFocus, toggleFocusedStory } from "$lib/helpers/shortcuts";
  import { clock } from "$lib/stores/clock.svelte";
  import { newCounts } from "$lib/stores/new-counts.svelte";
  import type { LayoutProps } from "./$types";

  let { data, children }: LayoutProps = $props();

  let shortcuts = $state<ShortcutsDialog>();

  const activeTopicId = $derived(page.route.id === "/topics/[id]" ? Number(page.params.id) : null);

  $effect(() => clock.subscribe());
  $effect(() => newCounts.subscribe(data.topics.map((t) => t.id)));

  function focusSearch() {
    const input = document.querySelector<HTMLInputElement>("[data-search-input]");
    if (input) input.focus();
    else goto(routes.search);
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.defaultPrevented || e.metaKey || e.ctrlKey || e.altKey || isTypingTarget(e.target)) return;
    if (document.querySelector("dialog[open]")) return;

    if (e.key === "/") {
      e.preventDefault();
      focusSearch();
    } else if (e.key === "?") {
      e.preventDefault();
      shortcuts?.show();
    } else if (e.key === "j" || e.key === "k") {
      e.preventDefault();
      moveHeadlineFocus(e.key === "j" ? 1 : -1);
    } else if (e.key === "o") {
      toggleFocusedStory();
    } else if (e.key === "Escape") {
      toggleFocusedStory(true);
    } else if (/^[1-9]$/.test(e.key)) {
      const topic = data.topics[Number(e.key) - 1];
      if (topic) goto(topicPath(topic.id));
    }
  }
</script>

<svelte:window onkeydown={onKeydown} />

{#if navigating.to}
  <div class="fixed inset-x-0 top-0 z-50 h-0.5 animate-pulse bg-primary" role="progressbar" aria-label="Učitavanje"></div>
{/if}

<AppHeader topics={data.topics} {activeTopicId} onShortcuts={() => shortcuts?.show()} />

<main class="mx-auto max-w-3xl px-3 py-4 sm:px-4 sm:py-6">
  {@render children()}
</main>

<AppFooter />

<ShortcutsDialog bind:this={shortcuts} topicCount={data.topics.length} />

<Toaster
  position="bottom-center"
  toastOptions={{
    unstyled: true,
    classes: {
      toast: "alert w-full shadow-lg text-sm",
      error: "alert-error",
      success: "alert-success",
    },
  }}
/>
