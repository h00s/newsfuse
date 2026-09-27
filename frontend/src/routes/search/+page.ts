import { loadFailed } from "$lib/helpers/errors";
import { searchHeadlines } from "$lib/services/headlines";
import { searchInput } from "$lib/types/search";
import type { PageLoad } from "./$types";

/** The query lives in ?q=, so a search survives a reload and can be linked. */
export const load: PageLoad = async ({ url, fetch }) => {
  const parsed = searchInput.safeParse({ query: url.searchParams.get("q") ?? "" });
  if (!parsed.success) return { query: url.searchParams.get("q") ?? "", headlines: null };

  try {
    return { query: parsed.data.query, headlines: await searchHeadlines(parsed.data.query, undefined, fetch) };
  } catch (e) {
    loadFailed(e);
  }
};
