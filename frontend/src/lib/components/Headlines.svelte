<script>
  import Headline from '$lib/components/Headline.svelte';

  export let headlines;
  export let sources;
  export let lastAccessedAt;
</script>

<div>
  {#each headlines as headline (headline.id)}
    <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} />
  {/each}
</div>
