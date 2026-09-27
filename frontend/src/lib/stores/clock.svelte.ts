const TICK_MS = 30 * 1000;

/** A shared "now" for relative times ("prije 5 minuta"), so they stay true while the page is
 *  open. Ticks only while subscribed. */
function createClock() {
  let now = $state(Date.now());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  function tick() {
    now = Date.now();
    timer = setTimeout(tick, TICK_MS);
  }

  return {
    get now(): number {
      return now;
    },
    /** `$effect(() => clock.subscribe())` keeps it ticking while the component lives. */
    subscribe(): () => void {
      consumers++;
      if (consumers === 1) timer = setTimeout(tick, TICK_MS);
      return () => {
        consumers = Math.max(0, consumers - 1);
        if (consumers === 0 && timer !== null) {
          clearTimeout(timer);
          timer = null;
        }
      };
    },
  };
}

export const clock = createClock();
