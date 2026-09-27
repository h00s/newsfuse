<script lang="ts">
  import { page } from "$app/state";
  import { routes } from "$lib/helpers/routes";
</script>

<svelte:head><title>Pogreška {page.status} · FusedNews</title></svelte:head>

<div class="flex flex-col items-center gap-4 py-24 text-center">
  <p class="text-sm text-base-content/60">Pogreška {page.status}</p>
  <h1 class="text-2xl font-semibold text-balance">{page.error?.message ?? "Nešto je pošlo po zlu."}</h1>
  <a class="btn btn-outline btn-sm" href={routes.home}>Natrag na početnu</a>
</div>
