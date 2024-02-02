import { PUBLIC_API_URL } from '$env/static/public'

export async function load({ fetch, params }) {
	let res = await fetch(PUBLIC_API_URL + '/headlines');
	const headlines = await res.json();
	
	res = await fetch(PUBLIC_API_URL + '/sources');
	const sources = await res.json();

	headlines.forEach(headline => {
		const source = sources.find(source => source.id === headline.sourceId);
		headline.Source = source;
	});

	return {
    headlines: headlines,
  };
}