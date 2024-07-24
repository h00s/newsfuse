<script>
  import Headline from "$lib/components/Headline.svelte";

  let lastAccessedAt = Date.now();
  export let searchedTerm;
  export let searchedHeadlines;
  export let sources;

  $: resultsNumber = searchedHeadlines ? searchedHeadlines.length : 0;
</script>

{#if searchedHeadlines !== undefined}
  <div class="mt-6 mb-2 text-center text-news-light">
    <p>Pronađeno je <span class="font-bold">{resultsNumber}</span> članaka za <span class="font-bold">"{searchedTerm}"</span></p>
  </div>
  <div class="border-t-2 border-news-light">
    {#each searchedHeadlines as headline (headline.id)}
      <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} />
      <hr class="border-news-separator">
    {/each}
  </div>
{/if}