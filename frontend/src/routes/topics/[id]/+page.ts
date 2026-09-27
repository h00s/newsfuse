import { error } from "@sveltejs/kit";
import { loadFailed } from "$lib/helpers/errors";
import { fetchHeadlines } from "$lib/services/headlines";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, url, parent, fetch }) => {
  const topicId = Number(params.id);
  if (!Number.isInteger(topicId) || topicId < 1) error(400, "Neispravna tema.");

  const rawSource = url.searchParams.get("source");
  const sourceId = rawSource === null ? null : Number(rawSource);
  if (sourceId !== null && (!Number.isInteger(sourceId) || sourceId < 1)) error(400, "Neispravan izvor.");

  const { topics, sources, unavailable } = await parent();
  if (unavailable) error(503, unavailable);
  const topic = topics.find((t) => t.id === topicId);
  if (!topic) error(404, "Ta tema ne postoji.");

  try {
    const headlines = await fetchHeadlines({ topicId, sourceId }, fetch);
    return { topic, sources: sources.filter((s) => s.topicId === topicId), sourceId, headlines };
  } catch (e) {
    loadFailed(e, { 404: "Taj izvor ne pripada ovoj temi." });
  }
};
