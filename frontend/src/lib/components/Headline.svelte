<script>
  import { formatDate } from '$lib/helpers/date';

  export let headline;
</script>

<div class="rounded overflow-hidden shadow-lg m-2 p-3 dark:bg-gray-800">
  <h3 class="text-gray-900 dark:text-white text-lg font-bold">
    <a class="mb-4" href="{ headline.URL }">{ headline.Title }</a>
  </h3>

  <span class="text-gray-400 mb-4">{ formatDate(headline.PublishedAt) } | { headline.Source.Name }</span>
</div>