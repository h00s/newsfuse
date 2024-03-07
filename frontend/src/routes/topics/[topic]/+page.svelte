<script>
  import { setLastAccessedAt } from '$lib/stores/topics'
  import { afterUpdate } from 'svelte';

  import Headlines from '$lib/components/Headlines.svelte';
  import Tabs from '$lib/components/Tabs.svelte';

  afterUpdate(async () => {
    setLastAccessedAt(data.selectedTopic);
  });

  export let data;
</script>

<Tabs topics={data.topics} selectedTopic={data.selectedTopic} />
<Headlines headlines={data.headlines} sources={data.sources} lastAccessedAt={data.lastAccessedAt} />
