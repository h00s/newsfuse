import { get } from 'svelte/store'
import { topicsLastAccessedAt } from './stores'

export function getLastAccessedAt(topic) {
  const topics = get(topicsLastAccessedAt)
  return topics[topic]
}

export function setLastAccessedAt(topic) {
  const topics = get(topicsLastAccessedAt)
  topics[topic] = Date.now()
  topicsLastAccessedAt.set(topics)
}
