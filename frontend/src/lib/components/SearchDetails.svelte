<script>
  import Headline from "$lib/components/Headline.svelte";

  let lastAccessedAt = Date.now();
  export let searchTerm;
  export let searchedHeadlines;
  export let sources;

  $: resultsNumber = searchedHeadlines ? searchedHeadlines.length : 0;
</script>

{#if searchedHeadlines !== undefined}
  <div class="mt-6 mb-2 text-center text-news-light">
    <p>Pronađeno je <span class="font-bold">{resultsNumber}</span> članaka za <span class="font-bold">"{searchTerm}"</span></p>
  </div>
  <div class="border">
    {#each searchedHeadlines as headline (headline.id)}
      <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} />
      <hr>
    {/each}
  </div>
{/if}

<style>
  .border {
    padding: 0;
    border: 1px solid #172F47;
    border-top: 2px solid #38BDF8;
  }
  
  hr {
    border-color: #172F47;
  }
</style>