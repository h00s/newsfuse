import { PUBLIC_API_URL } from '$env/static/public'

export async function load({ fetch, params }) {
	let [topics, headlines] = await Promise.all([
    fetch(`${PUBLIC_API_URL}/topics`).then(res => res.json()),
    fetch(`${PUBLIC_API_URL}/topics/${params.topic}/headlines`).then(res => res.json())
	]);

	return {
		topics: topics,
		selectedTopic: params.topic,
    headlines: headlines,
  };
}