<script>
  import { PUBLIC_API_URL } from '$env/static/public'
  import Headline from '$lib/components/Headline.svelte';

  function handleHeadlineDisplayed(event) {
    let headlineId = event.detail;
    if (headlineId == headlines[headlines.length - 1].id) {
      fetch(`${PUBLIC_API_URL}/topics/${selectedTopic}/headlines?last_id=${headlineId}`)
      .then((response) => response.json())
      .then((data) => {
        headlines = [...headlines, ...data];
      });
    }
  }

  export let selectedTopic;
  export let headlines;
  export let sources;
  export let lastAccessedAt;
</script>

<div class="border">
  {#each headlines as headline (headline.id)}
    <Headline {headline} source={sources[headline.source_id]} {lastAccessedAt} on:headlineDisplayed={handleHeadlineDisplayed} />
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
