import { PUBLIC_API_URL } from '$env/static/public'

export async function load({ fetch, params }) {
	let res = await fetch(`${PUBLIC_API_URL}/topics/${params.topic}/headlines`);
	const headlines = await res.json();

	return {
    headlines: headlines,
  };
}