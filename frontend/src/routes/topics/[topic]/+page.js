import { getLastAccessedAt } from '$lib/stores/topics'
import { fetchHeadlinesByTopic } from '$svc/headlines';
import { fetchSources } from '$svc/sources';
import { fetchTopics } from '$svc/topics';

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
