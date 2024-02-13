<script>
  import { PUBLIC_API_URL } from '$env/static/public'

  export let headlineId;
  let story;

  if (story === undefined) {
    fetch(`${PUBLIC_API_URL}/headlines/${headlineId}/story`)
      .then((response) => response.json())
      .then((data) => {
        story = data;
      });
  }
</script>

<div class="pt-4 text-gray-700 dark:text-gray-300">
  <p>
    {#if story === undefined}
      <span class="loading loading-spinner loading-sm"></span>
    {:else}
      {@html story.Content}
    {/if}
  </p>
</div>
