<script>
  export let data;
</script>

<div>
  {#each data.headlines as headline}
    <div class="rounded overflow-hidden shadow-lg m-2 p-3 dark:bg-gray-800">
      <h3 class="text-gray-900 dark:text-white text-lg font-bold">{headline.Title}</h3>
      <a class="text-gray-400 mb-4" href="{headline.URL}">{headline.URL}</a>
    </div>
  {/each}
</div>