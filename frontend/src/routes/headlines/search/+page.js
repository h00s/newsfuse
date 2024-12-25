import { fetchSources } from '$svc/sources';

export async function load({ fetch, params }) {
  let sources = await fetchSources(fetch);

  const sourcesMap = sources.reduce((map, obj) => {
    map[obj.id] = obj;
    return map;
  }, {});

  return {
    sources: sourcesMap,
  };
}