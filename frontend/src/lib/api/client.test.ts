import { describe, expect, it } from "vitest";
import { api, ApiError, isAbortError } from "$lib/api/client";

function fakeFetch(respond: () => Response) {
  const calls: { url: string; init?: RequestInit }[] = [];
  const fetch = (async (url: string | URL | Request, init?: RequestInit) => {
    calls.push({ url: String(url), init });
    return respond();
  }) as typeof globalThis.fetch;
  return { fetch, calls };
}

describe("api client", () => {
  it("prefixes /api/v1, encodes the query and drops empty values", async () => {
    const { fetch, calls } = fakeFetch(() => Response.json([]));
    await api.get("/headlines", { fetch, query: { topicId: 4, sourceId: null, beforeId: undefined, q: "č & d" } });
    expect(calls[0].url).toBe("/api/v1/headlines?topicId=4&q=%C4%8D+%26+d");
    expect(calls[0].init?.credentials).toBe("include");
  });

  it("parses a JSON body", async () => {
    const { fetch } = fakeFetch(() => Response.json([{ id: 1, name: "BBŽ" }]));
    await expect(api.get("/topics", { fetch })).resolves.toEqual([{ id: 1, name: "BBŽ" }]);
  });

  it("sends a JSON body with its content type only when there is one", async () => {
    const { fetch, calls } = fakeFetch(() => Response.json({}));
    await api.post("/stories/1/summarize", undefined, { fetch });
    await api.post("/things", { a: 1 }, { fetch });
    expect(calls[0].init?.method).toBe("POST");
    expect(calls[0].init?.body).toBeUndefined();
    expect(calls[0].init?.headers).toBeUndefined();
    expect(calls[1].init?.body).toBe('{"a":1}');
    expect(calls[1].init?.headers).toEqual({ "Content-Type": "application/json" });
  });

  it("resolves a 204 to undefined", async () => {
    const { fetch } = fakeFetch(() => new Response(null, { status: 204 }));
    await expect(api.post("/things", undefined, { fetch })).resolves.toBeUndefined();
  });

  it("throws ApiError with the status and Raptor's envelope", async () => {
    const body = { code: 404, message: "Topic not found" };
    const { fetch } = fakeFetch(() => Response.json(body, { status: 404 }));
    const error = (await api.get("/headlines", { fetch }).catch((e: unknown) => e)) as ApiError;
    expect(error).toBeInstanceOf(ApiError);
    expect(error.status).toBe(404);
    expect(error.body).toEqual(body);
  });

  it("throws ApiError when the error body is not JSON", async () => {
    const { fetch } = fakeFetch(() => new Response("<html>Bad Gateway</html>", { status: 502 }));
    const error = (await api.get("/topics", { fetch }).catch((e: unknown) => e)) as ApiError;
    expect(error).toBeInstanceOf(ApiError);
    expect(error.status).toBe(502);
  });

  it("recognizes an abort", () => {
    expect(isAbortError(new DOMException("Aborted", "AbortError"))).toBe(true);
    expect(isAbortError(new TypeError("Failed to fetch"))).toBe(false);
  });
});
