<script>
  import { PUBLIC_API_URL } from '$env/static/public'

  export let headlineId;
  let story;
  let buttonSummarizeEnabled = false;
  let buttonContentEnabled = true;
  let buttonSummarizeText = "Sažmi članak";
  let displayedStory = "";

  if (story === undefined) {
    fetch(`${PUBLIC_API_URL}/headlines/${headlineId}/story`)
      .then((response) => response.json())
      .then((data) => {
        story = data;
        if (story.content.length > 800) {
          story.long = true;
          story.intro = story.content.substring(0, 800).trim() + '...';
          story.summary === "" ? story.summarized = false : story.summarized = true;
          buttonContentEnabled = true;
          buttonSummarizeEnabled = true;
          displayedStory = story.intro;
        } else {
          story.long = false;
          displayedStory = story.content;
        }
      });
  }

  function summary() {
    buttonSummarizeText = '<span class="loading loading-spinner loading-sm"></span>';
    if (!story.summarized) {
      fetch(`${PUBLIC_API_URL}/stories/${story.id}/summarize`)
        .then((response) => response.json())
        .then((data) => {
          story.summary = data.summary;
          story.summarized = true;
          displayedStory = story.summary;
          buttonSummarizeEnabled = false;
        });
    } else {
      displayedStory = story.summary;
      buttonSummarizeEnabled = false;
    }
  }

  function content() {
    displayedStory = story.content;
    buttonContentEnabled = false;
  }
</script>

<div class="story pb-2 text-gray-700 dark:text-gray-300">
  {#if story === undefined}
    <span class="loading loading-spinner loading-sm"></span>
  {:else}
    {@html displayedStory}
    {#if story.long}
      <div class="flex justify-center pt-4">
        {#if buttonContentEnabled}
          <button class="btn btn-sm btn-primary text-white mr-1" on:click={content}>
            Cijeli članak
          </button>
        {/if}
        {#if buttonSummarizeEnabled}
          <button class="btn btn-sm btn-primary text-white" on:click={summary}>
            {@html buttonSummarizeText}
          </button>
        {/if}
      </div>
    {/if}
  {/if}
</div>
