import { describe, expect, it } from "vitest";
import { ApiError } from "$lib/api/client";
import { errorMessage } from "$lib/helpers/errors";

describe("errorMessage", () => {
  it("maps statuses to Croatian text and never shows the API's message", () => {
    const message = errorMessage(new ApiError("Topic not found", 404));
    expect(message).toBe("Traženi sadržaj ne postoji.");
    expect(errorMessage(new ApiError("Rate limit exceeded", 429))).toContain("Previše zahtjeva");
    expect(errorMessage(new ApiError("Summarization failed", 502))).toContain("nije dostupan");
    expect(errorMessage(new ApiError("Database error", 500))).toContain("poslužitelju");
  });

  it("prefers an override for the status", () => {
    expect(errorMessage(new ApiError("x", 404), { 404: "Tema ne postoji." })).toBe("Tema ne postoji.");
  });

  it("explains a network failure", () => {
    expect(errorMessage(new TypeError("Failed to fetch"))).toContain("Poslužitelj nije dostupan");
  });
});
