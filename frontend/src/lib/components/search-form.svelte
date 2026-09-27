<script lang="ts">
  import SearchIcon from "@lucide/svelte/icons/search";
  import { untrack } from "svelte";
  import { z } from "zod";
  import { searchInput } from "$lib/types/search";

  let { query = "", onSearch }: { query?: string; onSearch: (query: string) => void } = $props();

  const id = $props.id();
  // A one-time working copy: the page remounts this form for every new query.
  let value = $state(untrack(() => query));
  let error = $state<string>();

  function submit(e: SubmitEvent) {
    e.preventDefault();
    const result = searchInput.safeParse({ query: value });
    if (!result.success) {
      error = z.flattenError(result.error).fieldErrors.query?.[0];
      return;
    }
    error = undefined;
    onSearch(result.data.query);
  }
</script>

<form role="search" novalidate onsubmit={submit}>
  <label class={["input input-lg w-full bg-base-200", error && "input-error"]}>
    <SearchIcon class="size-5 opacity-60" aria-hidden="true" />
    <input
      id="{id}-query"
      type="search"
      class="grow"
      placeholder="Pretražite naslove…"
      autocomplete="off"
      enterkeyhint="search"
      aria-label="Pojam za pretragu"
      aria-invalid={error ? "true" : undefined}
      aria-describedby={error ? `${id}-error` : undefined}
      bind:value
      data-search-input
      {@attach (input) => input.focus()}
    />
    <kbd class="kbd hidden kbd-sm sm:inline-flex">Enter</kbd>
  </label>
  {#if error}
    <p id="{id}-error" class="mt-2 text-sm text-error">{error}</p>
  {/if}
</form>
