import { describe, expect, it } from "vitest";
import { parseReading } from "$lib/stores/reading.svelte";

describe("parseReading", () => {
  it("reads the stored state", () => {
    const raw = JSON.stringify({ lastSeen: { "2": 1000 }, lastTopicId: 2 });
    expect(parseReading(raw, null)).toEqual({ lastSeen: { "2": 1000 }, lastTopicId: 2 });
  });

  it("imports the old per-topic timestamps once", () => {
    // svelte-persisted-store kept an array indexed by topic id.
    expect(parseReading(null, JSON.stringify([null, 1000, null, 3000]))).toEqual({
      lastSeen: { "1": 1000, "3": 3000 },
      lastTopicId: null,
    });
  });

  it("starts empty on missing or corrupt data", () => {
    const empty = { lastSeen: {}, lastTopicId: null };
    expect(parseReading(null, null)).toEqual(empty);
    expect(parseReading("{not json", "also not")).toEqual(empty);
    expect(parseReading(JSON.stringify({ lastSeen: "x", lastTopicId: "y" }), null)).toEqual(empty);
  });
});
