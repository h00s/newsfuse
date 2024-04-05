<script>
  import { setLastAccessedAt } from '$lib/stores/topics'
  import { afterUpdate } from 'svelte';

  afterUpdate(async () => {
    if (selected) {
      setLastAccessedAt(topic.id);
    }
  });

  export let topic;
  export let selected;
</script>

<a role="tab" class="tab" href="/topics/{topic.id}" class:tab-active={selected}>
  {topic.name}
</a>