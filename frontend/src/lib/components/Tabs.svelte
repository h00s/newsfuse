<script>
  import { getLastAccessedAt } from '$lib/stores/topics'
  import Tab from '$lib/components/Tab.svelte'

  export let topics;
  export let selectedTopic;
</script>

<div role="tablist" class="tabs tabs-lg tabs-bordered mt-2">
  {#each topics as topic (topic.id)}
    <Tab topic={topic} selected={topic.id == selectedTopic} lastAccessedAt={getLastAccessedAt(topic.id)} />
  {/each}
</div>
