import { getLastAccessedAt } from '$lib/stores/topics'
import { fetchHeadlinesByTopic } from '$lib/repositories/headlines';
import { fetchSources } from '$lib/repositories/sources';
import { fetchTopics } from '$lib/repositories/topics';

export async function load({ fetch, params }) {
  let [topics, headlines, sources] = await Promise.all([
    fetchTopics(fetch),
    fetchHeadlinesByTopic(params.topic, fetch),
    fetchSources(fetch)
  ]);

  const sourcesMap = sources.reduce((map, obj) => {
    map[obj.id] = obj;
    return map;
  }, {});

  return {
    topics: topics,
    selectedTopic: params.topic,
    headlines: headlines,
    sources: sourcesMap,
    lastAccessedAt: getLastAccessedAt(params.topic)
  };
}
