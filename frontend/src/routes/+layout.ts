import { errorMessage } from "$lib/helpers/errors";
import { fetchSources } from "$lib/services/sources";
import { fetchTopics } from "$lib/services/topics";
import type { Source } from "$lib/types/source";
import type { Topic } from "$lib/types/topic";
import type { LayoutLoad } from "./$types";

/** A client SPA (adapter-static): the API exists only in the browser. */
export const ssr = false;

/** Topics and sources are reference data every page needs: loaded once, here. A failure doesn't
 *  fail the layout (that would leave only the static error page); the page's own load reports it
 *  inside the app shell instead. */
export const load: LayoutLoad = async ({ fetch }) => {
  try {
    const [topics, sources] = await Promise.all([fetchTopics(fetch), fetchSources(fetch)]);
    return { topics, sources, unavailable: null };
  } catch (e) {
    return { topics: [] as Topic[], sources: [] as Source[], unavailable: errorMessage(e) };
  }
};
