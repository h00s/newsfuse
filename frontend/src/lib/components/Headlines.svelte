<script>
  import { fetchHeadlinesByTopicAndLastId } from '$lib/repositories/headlines';
  import Headline from '$lib/components/Headline.svelte';

  async function handleHeadlineDisplayed(event) {
    let headlineId = event.detail;
    if (headlineId == headlines[headlines.length - 1].id) {
      let additionalHeadlines = await fetchHeadlinesByTopicAndLastId(selectedTopic, headlineId);
      headlines = [...headlines, ...additionalHeadlines];
    }
  }

  export let selectedTopic;
  export let headlines;
  export let sources;
  export let lastAccessedAt;
</script>

<div>
  {#each headlines as headline (headline.id)}
    <Headline {headline} source={sources[headline.sourceId]} {lastAccessedAt} on:headlineDisplayed={handleHeadlineDisplayed} />
    <hr class="border-news-separator">
  {/each}
</div>