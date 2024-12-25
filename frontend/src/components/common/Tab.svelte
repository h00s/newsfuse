<script>
  import { afterUpdate } from 'svelte';
  import { setLastAccessedAt } from '$lib/stores/topics'
  import { countNewSince } from '$svc/headlines'

  afterUpdate(() => {
    countNewSince(topic.id, lastAccessedAt).then((data) => {
      count = data.count;
      if (selected) {
        setLastAccessedAt(topic.id);
      }
    });
  });

  let count = 0;
  export let topic;
  export let selected;
  export let lastAccessedAt;
</script>

<a role="tab" class="tab" href="/topics/{topic.id}" class:tab-active={selected}>
  {topic.name}
  {#if count > 0}
    <span class="text-xs bg-green-500 text-white rounded-full px-2 py-1 ml-2">
      {count > 99 ? '99+' : count}
    </span>
  {/if}
</a>