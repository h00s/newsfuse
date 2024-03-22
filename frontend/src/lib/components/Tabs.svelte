<script>
  export let topics;
  export let selectedTopic;
</script>

<div role="tablist" class="tabs tabs-lg tabs-bordered mt-2">
  {#each topics as topic}
    <a role="tab" class="tab" href="/topics/{topic.id}" class:tab-active={topic.id == selectedTopic}>{topic.name}</a>
  {/each}
</div>
