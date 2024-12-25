<script>
  import Navbar from '$comp/layout/Navbar.svelte';
  import Search from '$comp/layout/Search.svelte';
  import Footer from '$comp/layout/Footer.svelte';

  export let data;
</script>

<Navbar />
<Search sources={data.sources} />
<Footer />