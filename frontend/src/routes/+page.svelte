<script>
  export let data;

  function formatDateTime(date) {
    date = new Date(date);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    return `${day}.${month}.${year}. ${hours}:${minutes}`;
  }
</script>

<div>
  {#each data.headlines as headline}
    <div class="rounded overflow-hidden shadow-lg m-2 p-3 dark:bg-gray-800">
      <h3 class="text-gray-900 dark:text-white text-lg font-bold">
        <a class="mb-4" href="{headline.URL}">{headline.Title}</a>
      </h3>
      <span class="text-gray-400 mb-4">{formatDateTime(headline.PublishedAt)} | {headline.Source}</span>
    </div>
  {/each}
</div>