import { api } from "$lib/api/client";
import type { Source } from "$lib/types/source";

export const fetchSources = (fetch?: typeof globalThis.fetch) => api.get<Source[]>("/sources", { fetch });
