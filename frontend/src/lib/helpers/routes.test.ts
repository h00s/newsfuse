import { describe, expect, it } from "vitest";
import { searchPath, topicPath } from "$lib/helpers/routes";

describe("routes", () => {
  it("builds topic paths with an optional source", () => {
    expect(topicPath(4)).toBe("/topics/4");
    expect(topicPath(4, null)).toBe("/topics/4");
    expect(topicPath(4, 9)).toBe("/topics/4?source=9");
  });

  it("encodes the search query", () => {
    expect(searchPath()).toBe("/search");
    expect(searchPath("Đakovo & 100%")).toBe("/search?q=%C4%90akovo%20%26%20100%25");
  });
});
