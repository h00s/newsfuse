<script>
  import Navbar from '$lib/components/Navbar.svelte';
  import Tabs from '$lib/components/Tabs.svelte';
  import Headlines from '$lib/components/Headlines.svelte';
  import Footer from '$lib/components/Footer.svelte';

  export let data;
</script>

<Navbar />
<Tabs topics={data.topics} selectedTopic={data.selectedTopic} />
<Headlines selectedTopic={data.selectedTopic} headlines={data.headlines} sources={data.sources} lastAccessedAt={data.lastAccessedAt} />
<Footer />
