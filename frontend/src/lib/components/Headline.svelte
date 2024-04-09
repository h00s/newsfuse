<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { humanizeDuration } from '$lib/helpers/date'
  import { inview } from 'svelte-inview';
  import Story from '$lib/components/Story.svelte'
  
  onMount(() => {
    newHeadline = lastAccessedAt < new Date(headline.published_at).getTime() ? true : false;
  });

  function sendDataToParent() {
    dispatch('headlineDisplayed', headline.id);
  }

  function toggleStory() {
    showStory = !showStory;
  }

  function handleHeadlineView({ detail }) {
    sendDataToParent();
    sourceLogoInView = true;
  }

  const headlineInViewOptions = {
    rootMargin: '50px',
    unobserveOnEnter: true,
  };
  const dispatch = createEventDispatcher();
  let showStory = false;
  let sourceLogoInView = false;
  export let headline;
  export let source;
  export let lastAccessedAt;
  export let newHeadline = false;
</script>

<div class="rounded overflow-hidden ml-2 p-2" use:inview={headlineInViewOptions} on:inview_enter="{handleHeadlineView}">
  <div class="flex justify-between items-center">
    {#if sourceLogoInView}
      <img src="/img/sources/{source.name}.webp" class="inline-block rounded-lg pr-3" alt="{source.name} logo" width="32" height="32">
    {:else}
      <div class="inline-block source-logo-placeholder"></div>
    {/if}
    <div class="flex-1">
      <h3 class="text-gray-900 inline"
        class:dark:text-white={newHeadline}
        class:font-bold={newHeadline}
        class:dark:text-slate-300={!newHeadline}
        class:font-medium={!newHeadline}
      >
        <a class="mb-4" href="{headline.url}">
          {headline.title}
        </a>
      </h3>

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

<style>
  .source-logo-placeholder {
    width: 32px;
    height: 36px;
  }
</style>