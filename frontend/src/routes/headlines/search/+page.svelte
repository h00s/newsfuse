<script>
  import Navbar from '$lib/components/Navbar.svelte';
  import Search from '$lib/components/Search.svelte';

  export let data;
</script>

<Navbar />
<Search sources={data.sources} />