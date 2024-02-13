<script>
  import Headline from '$lib/components/Headline.svelte';

  export let headlines;
</script>

<div>
  {#each headlines as headline (headline.ID)}
    <Headline headline={ headline } />
  {/each}
</div>