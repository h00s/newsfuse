import type { Source } from "$lib/types/source";

/** Headlines are written by the scrapers and only read here. */
export interface Headline {
  id: number;
  title: string;
  url: string;
  /** ISO 8601. */
  publishedAt: string;
  source: Source;
}

export interface HeadlineCount {
  count: number;
}
