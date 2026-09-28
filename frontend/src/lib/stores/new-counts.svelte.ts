import { browser } from "$app/environment";
import { untrack } from "svelte";
import { SvelteMap } from "svelte/reactivity";
import { countHeadlinesSince } from "$lib/services/headlines";
import { reading } from "$lib/stores/reading.svelte";

const POLL_MS = 2 * 60 * 1000;

/** The topic being read, and the moment its news counts from: what was new when it opened. */
interface Visit {
  topicId: number;
  since: number;
}

/** New-headline counts for the tab badges and the open topic's banner. Polls only while a
 *  consumer is mounted and the tab is visible, with a chained timeout so a slow response never
 *  stacks requests.
 *  - Every topic: headlines since this browser last opened it. For the open topic that is what
 *    arrived during the visit, which the banner offers to load.
 *  - The open topic: headlines since the visit began, the ones marked Novo, which its own tab
 *    shows until the reader moves to another topic. */
function createNewCountsStore() {
  const counts = new SvelteMap<number, number>();
  let visit = $state.raw<Visit | null>(null);
  let visitCount = $state(0);
  let topicIds: number[] = [];
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;
  let running = false;
  let again = false;

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
    if (running) {
      again = true; // asked for fresh numbers mid-poll: run once more right after
      return;
    }
    running = true;
    const current = untrack(() => visit);
    try {
      await Promise.all([
        ...topicIds.map(async (id) => {
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
        (async () => {
          if (!current) return;
          try {
            const { count } = await countHeadlinesSince(current.topicId, new Date(current.since));
            if (untrack(() => visit) === current) visitCount = count;
          } catch {
            // as above
          }
        })(),
      ]);
    } finally {
      running = false;
      schedule(again ? 0 : POLL_MS);
      again = false;
    }
  }

  if (browser) document.addEventListener("visibilitychange", () => schedule(0));

  return {
    /** Headlines since the topic was last opened. */
    count(topicId: number): number {
      return counts.get(topicId) ?? 0;
    },
    /** The open topic's headlines since its visit began; 0 for any other topic. */
    visitCount(topicId: number): number {
      return visit?.topicId === topicId ? visitCount : 0;
    },
    /** The reader opened the topic, or refreshed it, with everything after since marked new. */
    open(topicId: number, since: number) {
      counts.set(topicId, 0);
      visit = { topicId, since };
      visitCount = 0;
      schedule(0);
    },
    /** The reader moved on: the topic's news has been seen. */
    close(topicId: number) {
      if (untrack(() => visit)?.topicId !== topicId) return;
      visit = null;
      visitCount = 0;
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
