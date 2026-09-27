<script lang="ts">
  import SearchX from "@lucide/svelte/icons/search-x";
  import { goto } from "$app/navigation";
  import EmptyState from "$lib/components/empty-state.svelte";
  import HeadlineList from "$lib/components/headline-list.svelte";
  import SearchForm from "$lib/components/search-form.svelte";
  import { HeadlineFeed } from "$lib/feeds/headline-feed.svelte";
  import { searchPath } from "$lib/helpers/routes";
  import { searchHeadlines } from "$lib/services/headlines";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const feed = $derived(
    data.headlines && new HeadlineFeed(data.headlines, (beforeId) => searchHeadlines(data.query, beforeId)),
  );
</script>

<svelte:head><title>{data.headlines ? `${data.query} · Pretraga` : "Pretraga"} · FusedNews</title></svelte:head>

<h1 class="sr-only">Pretraga</h1>

<div class="space-y-5">
  {#key data.query}
    <SearchForm query={data.query} onSearch={(query) => goto(searchPath(query), { keepFocus: true })} />
  {/key}

  {#if feed}
    <p class="text-sm text-base-content/60">
      Naslovi koji sadrže <span class="font-semibold text-base-content">„{data.query}”</span>, najnoviji prvi
    </p>
    <HeadlineList {feed}>
      {#snippet empty()}
        <EmptyState icon={SearchX} title="Ništa nije pronađeno">
          Nijedan naslov ne sadrži „{data.query}”. Pokušajte s drugim ili kraćim pojmom.
        </EmptyState>
      {/snippet}
    </HeadlineList>
  {:else}
    <p class="text-center text-sm text-base-content/50">Pretražite naslove svih tema i izvora.</p>
  {/if}
</div>
