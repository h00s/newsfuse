<script>
  import { onMount } from 'svelte';
  import { searchHeadlines } from '$lib/repositories/headlines';
  import Headline from '$lib/components/Headline.svelte';
  
  onMount(() => {
    inputSearch.focus();
  });

  async function handleInputSearch(event) {
    if (event.key === 'Enter') {
      foundHeadlines = await searchHeadlines(inputSearch.value);
      inputSearch.blur();
    }
  }

  let inputSearch;
  let foundHeadlines;
  let lastAccessedAt = new Date();
  export let sources;
</script>

<div>
  <label class="input input-bordered flex items-center gap-2 bg-news-gray ml-1 mr-1 mb-2">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 opacity-70"><path fill-rule="evenodd" d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z" clip-rule="evenodd" /></svg>
    <input bind:this={inputSearch} on:keydown={handleInputSearch} type="text" class="grow" placeholder="Upišite tekst za pretragu..." />
  </label>

  {#if foundHeadlines !== undefined && foundHeadlines.length > 0}
    <div class="border mt-6">
      {#each foundHeadlines as headline (headline.id)}
        <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} />
        <hr>
      {/each}
    </div>
  {/if}
</div>

<style>
  .border {
    padding: 0;
    border: 1px solid #172F47;
  }
  
  hr {
    border-color: #172F47;
  }
</style>