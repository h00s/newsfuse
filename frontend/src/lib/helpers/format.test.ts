import { describe, expect, it } from "vitest";
import { absoluteTime, newHeadlinesLabel, relativeTime } from "$lib/helpers/format";

const now = Date.parse("2026-09-27T12:00:00Z");
const ago = (seconds: number) => new Date(now - seconds * 1000).toISOString();

describe("relativeTime", () => {
  it("says just now for the last few seconds", () => {
    expect(relativeTime(ago(10), now)).toBe("upravo sada");
  });

  it("counts minutes, hours and days in Croatian", () => {
    expect(relativeTime(ago(5 * 60), now)).toBe("prije 5 minuta");
    expect(relativeTime(ago(3 * 3600), now)).toBe("prije 3 sata");
    expect(relativeTime(ago(2 * 86400), now)).toBe("prekjučer");
    expect(relativeTime(ago(4 * 86400), now)).toBe("prije 4 dana");
  });

  it("falls back to the date after a week", () => {
    expect(relativeTime(ago(10 * 86400), now)).toBe(absoluteTime(ago(10 * 86400)));
  });
});

describe("newHeadlinesLabel", () => {
  it("follows Croatian plural forms", () => {
    expect(newHeadlinesLabel(1)).toBe("1 nova vijest");
    expect(newHeadlinesLabel(3)).toBe("3 nove vijesti");
    expect(newHeadlinesLabel(5)).toBe("5 novih vijesti");
    expect(newHeadlinesLabel(11)).toBe("11 novih vijesti");
    expect(newHeadlinesLabel(21)).toBe("21 nova vijest");
    expect(newHeadlinesLabel(24)).toBe("24 nove vijesti");
  });
});
