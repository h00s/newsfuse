<script>
  export let topics;
  export let selectedTopic;
</script>

<div role="tablist" class="tabs tabs-lg tabs-bordered mt-2 mb-3">
  {#each topics as topic}
    <a role="tab" class="tab" href="/topics/{topic.ID}" class:tab-active={topic.ID == selectedTopic}>{topic.Name}</a>
  {/each}
</div>
