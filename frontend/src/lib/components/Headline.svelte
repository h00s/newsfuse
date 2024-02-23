<script>
  import { humanizeDuration } from '$lib/helpers/date'
  import Story from '$lib/components/Story.svelte'

  let showStory = false;

  function toggleStory() {
    showStory = !showStory;
  }

  export let headline;
  export let source;
</script>

<div class="rounded overflow-hidden shadow-lg m-2 p-3 dark:bg-gray-800">
  <div class="flex justify-between items-center">
    <h3 class="text-gray-900 dark:text-white text-lg font-bold">
      <a class="mb-4" href="{headline.url}">{headline.title}</a>
    </h3>

    <button on:click="{toggleStory}" class="text-gray-300 focus:outline-none">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        {#if showStory}
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
        {:else}
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
        {/if}
      </svg>
    </button>
  </div>

  <span class="text-gray-400 mb-4">
    {humanizeDuration(headline.published_at)} | {source.name}
  </span>

  {#if showStory}
    <Story headlineId={headline.id} />
  {/if}
</div>
