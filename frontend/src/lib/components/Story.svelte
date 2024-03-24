<script>
  import { PUBLIC_API_URL } from '$env/static/public'

  export let headlineId;
  let story;
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
    fetch(`${PUBLIC_API_URL}/stories/${story.id}/summarize`)
      .then((response) => response.json())
      .then((data) => {
        story = data;
        storySummarized = true;
      });
  }
</script>

<div class="story pt-2 pb-4 text-gray-700 dark:text-gray-300">
  {#if story === undefined}
    <span class="loading loading-spinner loading-sm"></span>
  {:else}
    {#if !storySummarized}
      {#if story.content.length > 800}
        <div class="flex justify-center">
          <button class="btn btn-sm btn-primary text-white" on:click={summarizeStory}>
            Sažmi članak
          </button>
        </div>
      {/if}
      {@html story.content}
    {:else}
      {@html story.summary}
    {/if}
  {/if}
</div>
