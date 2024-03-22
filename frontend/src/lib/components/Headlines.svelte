<script>
  import Headline from '$lib/components/Headline.svelte';

  export let headlines;
  export let sources;
  export let lastAccessedAt;
</script>

<div class="border">
  {#each headlines as headline (headline.id)}
    <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} />
    <hr>
  {/each}
</div>

<style>
  .border {
    margin: 0;
    padding: 0;
    border: 1px solid #172F47;
  }
  
  hr {
    border-color: #172F47;
  }
</style>
