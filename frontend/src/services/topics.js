import { PUBLIC_API_URL } from '$env/static/public'

export async function fetchTopics(svelteFetch = undefined) {
  svelteFetch = svelteFetch || fetch;
  return svelteFetch(`${PUBLIC_API_URL}/topics`).then(res => res.json());
}