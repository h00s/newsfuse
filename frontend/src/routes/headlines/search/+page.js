import { fetchSources } from '$lib/repositories/sources';

export async function load({ fetch, params }) {
  let [sources] = await Promise.all([
    fetchSources(fetch)
  ]);

  const sourcesMap = sources.reduce((map, obj) => {
    map[obj.id] = obj;
    return map;
  }, {});

  return {
    sources: sourcesMap,
  };
}