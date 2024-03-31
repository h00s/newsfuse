<script>
  import { setLastAccessedAt } from '$lib/stores/topics'
  import { afterUpdate } from 'svelte';

  import Navbar from '$lib/components/Navbar.svelte';
  import Tabs from '$lib/components/Tabs.svelte';
  import Headlines from '$lib/components/Headlines.svelte';
  import Footer from '$lib/components/Footer.svelte';

  afterUpdate(async () => {
    setLastAccessedAt(data.selectedTopic);
  });

  export let data;
</script>

<Navbar />
<Tabs topics={data.topics} selectedTopic={data.selectedTopic} />
<Headlines topic={data.selectedTopic} headlines={data.headlines} sources={data.sources} lastAccessedAt={data.lastAccessedAt} />
<Footer />

