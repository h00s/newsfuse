<script>
  export let data;
</script>

<div>
  {#each data.headlines as headline}
    <div class="pb-4">
      <h3 class="text-lg font-bold">{headline.Title}</h3>
      <a class="text-gray-400 mb-4" href="{headline.URL}">{headline.URL}</a>
      <hr class="mt-4">
    </div>
  {/each}
</div>