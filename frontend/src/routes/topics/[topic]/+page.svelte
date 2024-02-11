<script>
  import Headlines from '$lib/components/Headlines.svelte';

  export let data;
</script>

<Headlines headlines={ data.headlines } />
