import { redirect } from "@sveltejs/kit";
import { routes } from "$lib/helpers/routes";
import type { PageLoad } from "./$types";

/** The search page's old address, kept for bookmarks. */
export const load: PageLoad = () => {
  redirect(307, routes.search);
};
