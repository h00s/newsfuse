<script>
  import Headlines from '$lib/components/Headlines.svelte';
  import Tabs from '$lib/components/Tabs.svelte';

  export let data;
</script>

<Tabs topics={data.topics} selectedTopic={data.selectedTopic} />
<Headlines headlines={data.headlines} sources={data.sources} />
