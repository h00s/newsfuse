import { describe, expect, it } from "vitest";
import { stepIndex } from "$lib/helpers/shortcuts";

describe("stepIndex", () => {
  it("starts at the first or last item when nothing is focused", () => {
    expect(stepIndex(-1, 5, 1)).toBe(0);
    expect(stepIndex(-1, 5, -1)).toBe(4);
  });

  it("moves and stops at the ends", () => {
    expect(stepIndex(2, 5, 1)).toBe(3);
    expect(stepIndex(4, 5, 1)).toBe(4);
    expect(stepIndex(0, 5, -1)).toBe(0);
  });

  it("has nowhere to go in an empty list", () => {
    expect(stepIndex(-1, 0, 1)).toBe(-1);
  });
});
