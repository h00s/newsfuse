<script>
  import Navbar from '$lib/components/Navbar.svelte';
  import Search from '$lib/components/Search.svelte';
  import Footer from '$lib/components/Footer.svelte';

  export let data;
</script>

<Navbar />
<Search sources={data.sources} />
<Footer />