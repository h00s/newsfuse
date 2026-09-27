import { describe, expect, it } from "vitest";
import { sourceInitials, sourceLogo } from "$lib/helpers/sources";

describe("sources", () => {
  it("takes initials from the first two words, or the first two letters", () => {
    expect(sourceInitials("Hacker News")).toBe("HN");
    expect(sourceInitials("Index.hr")).toBe("IH");
    expect(sourceInitials("Bug")).toBe("BU");
  });

  it("encodes the logo file name", () => {
    expect(sourceLogo({ name: "Radio Daruvar" })).toBe("/img/sources/Radio%20Daruvar.webp");
  });
});
