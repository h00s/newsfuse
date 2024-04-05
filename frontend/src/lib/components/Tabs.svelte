<script>
  import Tab from '$lib/components/Tab.svelte'

  export let topics;
  export let selectedTopic;
</script>

<div role="tablist" class="tabs tabs-lg tabs-bordered mt-2">
  {#each topics as topic}
    <Tab topic={topic} selected={topic.id == selectedTopic} />
  {/each}
</div>
