<script>
  import { onMount } from 'svelte';
  import { fetchStoryByHeadline, fetchStorySummary } from '$svc/stories';

  onMount(() => {
    if (story === undefined) {
      fetchStoryByHeadline(headlineId).then((data) => {
        story = data;
        if (story.content.length > 650) {
          story.long = true;
          story.intro = story.content.substring(0, 600).trim() + '...';
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
  });

  function summary() {
    buttonSummarizeText = '<span class="loading loading-spinner loading-sm"></span>';
    if (!story.summarized) {
      fetchStorySummary(story.id).then((data) => {
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

  let story;
  let buttonSummarizeEnabled = false;
  let buttonContentEnabled = false;
  let buttonSummarizeText = "Sažmi članak";
  let displayedStory = "";
  export let headlineId;
</script>

<div class="story pb-2 text-gray-300">
  {#if story === undefined}
    <div class="pt-4">
      <span class="loading loading-spinner loading-sm"></span>
    </div>
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
