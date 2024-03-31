import { PUBLIC_API_URL } from '$env/static/public'

export async function fetchSources(svelteFetch = undefined) {
  svelteFetch = svelteFetch || fetch;
  return svelteFetch(`${PUBLIC_API_URL}/sources`).then(res => res.json());
}