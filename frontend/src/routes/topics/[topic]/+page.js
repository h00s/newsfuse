import { PUBLIC_API_URL } from '$env/static/public'
import { getLastAccessedAt } from '$lib/stores/topics'

export async function load({ fetch, params }) {
  let [topics, headlines, sources] = await Promise.all([
    fetch(`${PUBLIC_API_URL}/topics`).then(res => res.json()),
    fetch(`${PUBLIC_API_URL}/topics/${params.topic}/headlines`).then(res => res.json()),
    fetch(`${PUBLIC_API_URL}/sources`).then(res => res.json())
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
