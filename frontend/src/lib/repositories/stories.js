import { PUBLIC_API_URL } from '$env/static/public'

export async function fetchStoryByHeadline(headlineId, svelteFetch = undefined) {
  svelteFetch = svelteFetch || fetch;
  return svelteFetch(`${PUBLIC_API_URL}/headlines/${headlineId}/story`).then(res => res.json());
}