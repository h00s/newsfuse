import { error, redirect } from "@sveltejs/kit";
import { topicPath } from "$lib/helpers/routes";
import { reading } from "$lib/stores/reading.svelte";
import type { PageLoad } from "./$types";

/** Home is the topic read last, or the first one on a first visit. */
export const load: PageLoad = async ({ parent }) => {
  const { topics, unavailable } = await parent();
  if (unavailable) error(503, unavailable);

  const topic = topics.find((t) => t.id === reading.lastTopicId) ?? topics[0];
  if (!topic) error(404, "Još nema nijedne teme.");
  redirect(307, topicPath(topic.id));
};
