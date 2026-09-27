import { describe, expect, it } from "vitest";
import { HeadlineFeed } from "$lib/feeds/headline-feed.svelte";
import type { Headline } from "$lib/types/headline";

const source = { id: 9, name: "Bug", topicId: 4, isScrapable: true };
const headline = (id: number): Headline => ({
  id,
  title: `Naslov ${id}`,
  url: `https://example.test/${id}`,
  publishedAt: "2026-09-27T12:00:00Z",
  source,
});
const range = (from: number, to: number) => Array.from({ length: from - to + 1 }, (_, i) => headline(from - i));

describe("HeadlineFeed", () => {
  it("is done when the first page is short", () => {
    const feed = new HeadlineFeed(range(10, 1), async () => [], 30);
    expect(feed.done).toBe(true);
  });

  it("loads the next page after the last id and appends it", async () => {
    const asked: number[] = [];
    const feed = new HeadlineFeed(range(60, 31), async (beforeId) => {
      asked.push(beforeId);
      return range(30, 1);
    }, 30);

    await feed.loadMore();

    expect(asked).toEqual([31]);
    expect(feed.items.map((h) => h.id)).toEqual(range(60, 1).map((h) => h.id));
    expect(feed.done).toBe(false);
  });

  it("is done after a short page and skips headlines it already has", async () => {
    const feed = new HeadlineFeed(range(40, 11), async () => [headline(11), headline(10), headline(9)], 30);

    await feed.loadMore();

    expect(feed.items.map((h) => h.id).slice(-3)).toEqual([11, 10, 9]);
    expect(feed.items).toHaveLength(32);
    expect(feed.done).toBe(true);
  });

  it("loads one page at a time", async () => {
    let calls = 0;
    const feed = new HeadlineFeed(range(40, 11), async () => {
      calls++;
      return [];
    }, 30);

    await Promise.all([feed.loadMore(), feed.loadMore()]);

    expect(calls).toBe(1);
  });

  it("marks a failure and rethrows it for the caller to report", async () => {
    const feed = new HeadlineFeed(range(40, 11), async () => {
      throw new TypeError("Failed to fetch");
    }, 30);

    await expect(feed.loadMore()).rejects.toThrow("Failed to fetch");
    expect(feed.failed).toBe(true);
    expect(feed.loading).toBe(false);
    expect(feed.done).toBe(false);
  });
});
