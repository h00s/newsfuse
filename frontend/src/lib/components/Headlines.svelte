<script>
  import Headline from '$lib/components/Headline.svelte';

  export let headlines;
  export let sources;
</script>

<div>
  {#each headlines as headline (headline.id)}
    <Headline headline={headline} source={sources[headline.source_id]} />
  {/each}
</div>
