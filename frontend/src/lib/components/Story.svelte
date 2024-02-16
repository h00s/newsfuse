<script>
  import { PUBLIC_API_URL } from '$env/static/public'

  export let headlineId;
  let story;
  let buttonSummarizeText = "Sažmi članak";

  if (story === undefined) {
    fetch(`${PUBLIC_API_URL}/headlines/${headlineId}/story`)
      .then((response) => response.json())
      .then((data) => {
        story = data;
      });
  }

  function summarizeStory() {
    buttonSummarizeText = '<span class="loading loading-spinner loading-sm"></span>';
    fetch(`${PUBLIC_API_URL}/stories/${story.ID}/summarize`)
    .then((response) => response.json())
      .then((data) => {
        story = data;
        story.Content = story.Summary;
      });
  }
</script>

<div class="pt-4 text-gray-700 dark:text-gray-300">
  <p>
    {#if story === undefined}
    <span class="loading loading-spinner loading-sm"></span>
    {:else}
    {#if story.Content.length > 800}
        <div class="flex justify-center pb-2">
          <button class="btn btn-sm btn-primary" on:click={summarizeStory}>{@html buttonSummarizeText}</button>
        </div>
      {/if}
      {@html story.Content}
    {/if}
  </p>
</div>
