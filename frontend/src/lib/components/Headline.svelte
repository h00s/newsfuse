<script>
  import { onMount } from 'svelte';
  import { humanizeDuration } from '$lib/helpers/date'
  import { inview } from 'svelte-inview';
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

<div class="rounded overflow-hidden ml-2 p-2">
  <div class="flex justify-between items-center">
    <img src="/sources/{source.name}.webp" class="inline-block rounded-lg pr-3" alt="{source.name} logo" width="32" height="32" use:inview>
    <div class="flex-1">
      <h3 class="text-gray-900 dark:text-white font-bold inline">
        <a class="mb-4" href="{headline.url}">
          {headline.title}
        </a>
      </h3>

      {#if newHeadline}
        <span class="text-xs bg-green-500 text-white rounded-full px-2 py-1">NOVO</span>
      {/if}

      <span class="text-gray-400 pl-1">
        {humanizeDuration(headline.published_at)}
      </span>
    </div>

    {#if source.is_scrapable}
      <button on:click="{toggleStory}" class="text-news-light focus:outline-none pl-1">
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

  {#if showStory}
    <Story headlineId={headline.id} />
  {/if}
</div>