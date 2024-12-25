<script>
  import Navbar from '$comp/layout/Navbar.svelte';
  import Tabs from '$comp/layout/Tabs.svelte';
  import Headlines from '$comp/layout/Headlines.svelte';
  import Footer from '$comp/layout/Footer.svelte';

  export let data;
</script>

<Navbar />
<Tabs topics={data.topics} selectedTopic={data.selectedTopic} />
<Headlines selectedTopic={data.selectedTopic} headlines={data.headlines} sources={data.sources} lastAccessedAt={data.lastAccessedAt} />
<Footer />
