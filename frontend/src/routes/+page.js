import { PUBLIC_API_URL } from '$env/static/public'

export async function load({ fetch, params }) {
	const res = await fetch(PUBLIC_API_URL + '/headlines');
	const headlines = await res.json();

	return {
    headlines: headlines
  };
}