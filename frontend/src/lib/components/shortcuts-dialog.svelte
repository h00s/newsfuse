<script lang="ts">
  let { topicCount }: { topicCount: number } = $props();

  let dialog = $state<HTMLDialogElement>();

  export function show() {
    dialog?.showModal();
  }

  const shortcuts = $derived([
    { keys: ["/"], action: "Pretraga" },
    { keys: ["j", "k"], action: "Sljedeći ili prethodni naslov" },
    { keys: ["Enter"], action: "Otvori ili zatvori članak" },
    { keys: ["Esc"], action: "Zatvori članak" },
    { keys: topicCount > 1 ? ["1", String(Math.min(topicCount, 9))] : ["1"], action: "Prebaci na temu", range: true },
    { keys: ["?"], action: "Ovaj popis" },
  ]);
</script>

<dialog bind:this={dialog} class="modal modal-bottom sm:modal-middle" aria-labelledby="shortcuts-title">
  <div class="modal-box">
    <h2 id="shortcuts-title" class="text-lg font-semibold">Prečaci tipkovnice</h2>
    <dl class="mt-4 grid grid-cols-[auto_1fr] items-center gap-x-6 gap-y-3 text-sm">
      {#each shortcuts as shortcut (shortcut.action)}
        <dt class="flex items-center gap-1">
          {#each shortcut.keys as key, i (key)}
            {#if i > 0}<span class="text-base-content/50">{shortcut.range ? "–" : "/"}</span>{/if}
            <kbd class="kbd kbd-sm">{key}</kbd>
          {/each}
        </dt>
        <dd class="text-base-content/80">{shortcut.action}</dd>
      {/each}
    </dl>
    <div class="modal-action">
      <form method="dialog"><button class="btn btn-sm">Zatvori</button></form>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button tabindex="-1">Zatvori</button></form>
</dialog>
