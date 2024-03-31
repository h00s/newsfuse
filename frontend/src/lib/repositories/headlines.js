import { PUBLIC_API_URL } from '$env/static/public'

export async function fetchHeadlinesByTopic(topic, svelteFetch = undefined) {
  svelteFetch = svelteFetch || fetch;
  return svelteFetch(`${PUBLIC_API_URL}/topics/${topic}/headlines`).then(res => res.json());
}

export async function fetchHeadlinesByTopicAndLastId(topic, lastId, svelteFetch = undefined) {
  svelteFetch = svelteFetch || fetch;
  return svelteFetch(`${PUBLIC_API_URL}/topics/${topic}/headlines?last_id=${lastId}`).then(res => res.json());
}