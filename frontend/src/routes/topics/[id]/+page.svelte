<script lang="ts">
  import TopicView from "$lib/components/topic-view.svelte";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();
</script>

<svelte:head><title>{data.topic.name} · FusedNews</title></svelte:head>

<h1 class="sr-only">{data.topic.name}</h1>

<!-- Another topic is a fresh view: what counts as new is decided when it opens. -->
{#key data.topic.id}
  <TopicView topic={data.topic} sources={data.sources} sourceId={data.sourceId} headlines={data.headlines} />
{/key}
