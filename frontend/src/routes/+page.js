export async function load({ fetch, params }) {
	const res = await fetch('http://localhost:3000/api/v1/headlines');
	const headlines = await res.json();

	return {
    headlines: headlines
  };
}