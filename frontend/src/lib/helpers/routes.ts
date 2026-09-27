/** One root per page. Derive paths from these instead of pasting literals, so a rename propagates. */
export const routes = {
  home: "/",
  topics: "/topics",
  search: "/search",
} as const;

/** A topic, optionally narrowed to one of its sources. The source lives in the query string so the
 *  filtered view survives a reload and can be linked. */
export const topicPath = (topicId: number, sourceId?: number | null) =>
  sourceId ? `${routes.topics}/${topicId}?source=${sourceId}` : `${routes.topics}/${topicId}`;

export const searchPath = (query?: string) =>
  query ? `${routes.search}?q=${encodeURIComponent(query)}` : routes.search;
