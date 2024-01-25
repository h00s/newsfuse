<script>
  export let data;
  console.log(data.headlines);
</script>

<div>
  {#each data.headlines as headline}
    <div class="pb-8">
      <h3 class="text-lg font-medium">{headline.Title}</h3>
      <a href="{headline.URL}">{headline.URL}</a>
      <hr>
    </div>
  {/each}
</div>