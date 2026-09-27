/** What this browser has read: when each topic was last opened, which marks newer headlines as
 *  new, and which topic was open last, which "/" returns to. Kept in localStorage, so it is
 *  per-browser convenience and survives a missing or blocked storage by starting empty. */

export interface ReadingState {
  /** Topic id → epoch milliseconds. */
  lastSeen: Record<string, number>;
  lastTopicId: number | null;
}

const KEY = "fusednews:reading";
/** Where svelte-persisted-store kept the per-topic timestamps before. */
const LEGACY_KEY = "topicsLastAccessedAt";

export function parseReading(raw: string | null, legacyRaw: string | null): ReadingState {
  const state: ReadingState = { lastSeen: {}, lastTopicId: null };
  try {
    if (raw !== null) {
      const parsed = JSON.parse(raw);
      if (parsed && typeof parsed.lastSeen === "object" && !Array.isArray(parsed.lastSeen)) {
        for (const [id, at] of Object.entries(parsed.lastSeen)) {
          if (typeof at === "number") state.lastSeen[id] = at;
        }
      }
      if (typeof parsed?.lastTopicId === "number") state.lastTopicId = parsed.lastTopicId;
      return state;
    }
    if (legacyRaw !== null) {
      const legacy = JSON.parse(legacyRaw);
      if (legacy && typeof legacy === "object") {
        for (const [id, at] of Object.entries(legacy)) {
          if (typeof at === "number") state.lastSeen[id] = at;
        }
      }
    }
  } catch {
    // corrupt storage: start empty
  }
  return state;
}

function read(key: string): string | null {
  try {
    return globalThis.localStorage?.getItem(key) ?? null;
  } catch {
    return null;
  }
}

function createReadingStore() {
  const state = $state<ReadingState>(parseReading(read(KEY), read(LEGACY_KEY)));

  function save() {
    try {
      globalThis.localStorage?.setItem(KEY, JSON.stringify(state));
    } catch {
      // private mode or blocked storage: the markers last until the tab closes
    }
  }

  return {
    /** When the topic was last opened, or undefined if never. */
    lastSeen(topicId: number): number | undefined {
      return state.lastSeen[topicId];
    },
    get lastTopicId(): number | null {
      return state.lastTopicId;
    },
    /** The topic is open now: its headlines up to this moment count as seen. */
    markSeen(topicId: number, at = Date.now()) {
      state.lastSeen[topicId] = at;
      state.lastTopicId = topicId;
      save();
    },
  };
}

export const reading = createReadingStore();
