import { api } from "$lib/api/client";
import type { Headline, HeadlineCount } from "$lib/types/headline";

export interface HeadlineFilter {
  topicId: number;
  /** One of the topic's sources. */
  sourceId?: number | null;
  /** Continue after this headline; omit for the first page. */
  beforeId?: number;
}

/** A page of a topic's headlines, newest first. */
export const fetchHeadlines = ({ topicId, sourceId, beforeId }: HeadlineFilter, fetch?: typeof globalThis.fetch) =>
  api.get<Headline[]>("/headlines", { query: { topicId, sourceId, beforeId }, fetch });

/** A page of headlines whose title contains query, newest first. */
export const searchHeadlines = (query: string, beforeId?: number, fetch?: typeof globalThis.fetch) =>
  api.get<Headline[]>("/headlines/search", { query: { query, beforeId }, fetch });

/** How many of a topic's headlines were published after since. */
export const countHeadlinesSince = (topicId: number, since: Date, fetch?: typeof globalThis.fetch) =>
  api.get<HeadlineCount>("/headlines/count", { query: { topicId, since: since.toISOString() }, fetch });
