import { api } from "$lib/api/client";
import type { Topic } from "$lib/types/topic";

export const fetchTopics = (fetch?: typeof globalThis.fetch) => api.get<Topic[]>("/topics", { fetch });
