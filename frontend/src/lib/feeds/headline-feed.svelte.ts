import type { Headline } from "$lib/types/headline";

/** The API's page size: a shorter page means there is nothing more. */
export const PAGE_SIZE = 30;

/** An infinite list of headlines for one page instance: seeded from the load, extended page by
 *  page. A new load (another topic, a refresh) makes a new feed rather than re-seeding this one. */
export class HeadlineFeed {
  items = $state.raw<Headline[]>([]);
  loading = $state(false);
  done = $state(false);
  failed = $state(false);

  readonly #fetchPage: (beforeId: number) => Promise<Headline[]>;
  readonly #pageSize: number;

  constructor(initial: Headline[], fetchPage: (beforeId: number) => Promise<Headline[]>, pageSize = PAGE_SIZE) {
    this.items = initial;
    this.done = initial.length < pageSize;
    this.#fetchPage = fetchPage;
    this.#pageSize = pageSize;
  }

  /** Appends the next page. Rethrows a failure for the caller to report; failed then stays set
   *  until a retry, so an in-view trigger doesn't retry on its own. */
  async loadMore(): Promise<void> {
    const last = this.items.at(-1);
    if (this.loading || this.done || !last) return;

    this.loading = true;
    this.failed = false;
    try {
      const page = await this.#fetchPage(last.id);
      const known = new Set(this.items.map((h) => h.id));
      this.items = [...this.items, ...page.filter((h) => !known.has(h.id))];
      this.done = page.length < this.#pageSize;
    } catch (e) {
      this.failed = true;
      throw e;
    } finally {
      this.loading = false;
    }
  }
}
