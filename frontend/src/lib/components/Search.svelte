<script>
  import { onMount } from 'svelte';

  let searchInput;

  onMount(() => {
    searchInput.focus();
  });
</script>

<div>
  <label class="input input-bordered flex items-center gap-2 bg-news-gray ml-1 mr-1">
    <input bind:this={searchInput} type="text" class="grow" placeholder="Upišite tekst za pretragu..." />
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 opacity-70"><path fill-rule="evenodd" d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z" clip-rule="evenodd" /></svg>
  </label>
</div>