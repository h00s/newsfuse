import { api } from "$lib/api/client";
import type { Story } from "$lib/types/story";

/** The article behind a headline; the server scrapes it on the first request. */
export const fetchStory = (headlineId: number, fetch?: typeof globalThis.fetch) =>
  api.get<Story>(`/headlines/${headlineId}/story`, { fetch });

/** Summarizes a story once (a paid LLM call, rate-limited); later calls return the stored summary. */
export const summarizeStory = (storyId: number, fetch?: typeof globalThis.fetch) =>
  api.post<Story>(`/stories/${storyId}/summarize`, undefined, { fetch });
