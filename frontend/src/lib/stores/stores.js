import { persisted } from 'svelte-persisted-store'

export const topicsLastAccessedAt = persisted('topicsLastAccessedAt', [])
