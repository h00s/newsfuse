<script>
  import { onMount } from 'svelte';
  import { humanizeDuration } from '$lib/helpers/date'
  import Story from '$lib/components/Story.svelte'

  let showStory = false;

  function toggleStory() {
    showStory = !showStory;
  }

  onMount(() => {
    newHeadline = lastAccessedAt < new Date(headline.published_at).getTime() ? true : false;
  });

  export let headline;
  export let source;
  export let lastAccessedAt;
  export let newHeadline = false;
</script>

<div class="rounded overflow-hidden m-2 p-2">
  <div class="flex justify-between items-center">
    <h3 class="text-gray-900 dark:text-white font-bold">
      <a class="mb-4" href="{headline.url}">{headline.title}</a>
      {#if newHeadline}
        <span class="text-xs bg-green-500 text-white rounded-full px-2 py-1">NOVO</span>
      {/if}
    </h3>

    {#if source.is_scrapable}
      <button on:click="{toggleStory}" class="text-news-light focus:outline-none">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          {#if showStory}
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
          {:else}
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
          {/if}
        </svg>
      </button>
    {/if}
  </div>

  <span class="text-gray-400 mb-4">
    {humanizeDuration(headline.published_at)} | {source.name}
  </span>

  {#if showStory}
    <Story headlineId={headline.id} />
  {/if}
</div>
