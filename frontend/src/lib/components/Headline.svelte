<script>
  import { humanizeDuration } from '$lib/helpers/date';

  let showStory = false;

  function toggleStory() {
    showStory = !showStory;
  }

  export let headline;
</script>

<div class="rounded overflow-hidden shadow-lg m-2 p-3 dark:bg-gray-800">
  <div class="flex justify-between items-center">
    <h3 class="text-gray-900 dark:text-white text-lg font-bold">
      <a class="mb-4" href="{ headline.URL }">{ headline.Title }</a>
    </h3>
    <button on:click="{toggleStory}" class="text-gray-300 focus:outline-none">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
    </button>
  </div>

  <span class="text-gray-400 mb-4">{ humanizeDuration(headline.PublishedAt) } | { headline.Source.Name }</span>

  {#if showStory}
    <div class="story pt-4 text-gray-700 dark:text-gray-300">
      <p> uskoro... </p>
    </div>
  {/if}
</div>
