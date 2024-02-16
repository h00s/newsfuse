<script>
  import { PUBLIC_API_URL } from '$env/static/public'

  export let headlineId;
  let story;
  let buttonSummarizeText = "Sažmi članak";
  let storySummarized = false;

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
        storySummarized = true;
      });
  }
</script>

<div class="pt-4 text-gray-700 dark:text-gray-300">
  <p>
    {#if story === undefined}
      <span class="loading loading-spinner loading-sm"></span>
    {:else}
      {#if !storySummarized}
        {#if story.Content.length > 800}
          <div class="flex justify-center pb-2">
            <button class="btn btn-sm btn-primary" on:click={summarizeStory}>{@html buttonSummarizeText}</button>
          </div>
        {/if}
        {@html story.Content}
      {:else}
        {@html story.Summary}
      {/if}
    {/if}
  </p>
</div>
