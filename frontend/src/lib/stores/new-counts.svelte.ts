import { browser } from "$app/environment";
import { untrack } from "svelte";
import { SvelteMap } from "svelte/reactivity";
import { countHeadlinesSince } from "$lib/services/headlines";
import { reading } from "$lib/stores/reading.svelte";

const POLL_MS = 2 * 60 * 1000;

/** How many headlines each topic has gained since this browser last opened it: the tab badges
 *  and the open topic's "new headlines" banner. Polls only while a consumer is mounted and the
 *  tab is visible, with a chained timeout so a slow response never stacks requests. */
function createNewCountsStore() {
  const counts = new SvelteMap<number, number>();
  let topicIds: number[] = [];
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;
  let running = false;

  function stop() {
    if (timer !== null) clearTimeout(timer);
    timer = null;
  }

  // untrack: consumers subscribe from effects, which must not come to depend on this store's
  // internals.
  function schedule(delay: number) {
    untrack(() => {
      stop();
      if (!browser || consumers === 0 || topicIds.length === 0 || document.visibilityState === "hidden") return;
      timer = setTimeout(tick, delay);
    });
  }

  async function tick() {
    timer = null;
    if (running) return;
    running = true;
    try {
      await Promise.all(
        topicIds.map(async (id) => {
          const since = reading.lastSeen(id);
          if (since === undefined) return; // never opened: nothing to compare against
          try {
            const { count } = await countHeadlinesSince(id, new Date(since));
            // Opened meanwhile: this count is against the old visit, so drop it.
            if (reading.lastSeen(id) === since) counts.set(id, count);
          } catch {
            // no toast for a background poll: the next one retries
          }
        }),
      );
    } finally {
      running = false;
      schedule(POLL_MS);
    }
  }

  if (browser) document.addEventListener("visibilitychange", () => schedule(0));

  return {
    count(topicId: number): number {
      return counts.get(topicId) ?? 0;
    },
    /** The topic was just opened or refreshed: nothing in it is new any more. */
    seen(topicId: number) {
      counts.set(topicId, 0);
    },
    /** `$effect(() => newCounts.subscribe(ids))` polls these topics while the component lives. */
    subscribe(ids: number[]): () => void {
      topicIds = ids;
      consumers++;
      schedule(0);
      return () => {
        consumers = Math.max(0, consumers - 1);
        if (consumers === 0) stop();
      };
    },
  };
}

export const newCounts = createNewCountsStore();
