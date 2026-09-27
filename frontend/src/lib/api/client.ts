/** Raptor's error envelope: `{"code":404,"message":"Topic not found","attrs":{…}}`. `message` is
 *  English and meant for developers; the UI shows errorMessage() instead. */
export interface ApiErrorBody {
  code?: number;
  message?: string;
  attrs?: Record<string, unknown>;
}

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public body?: ApiErrorBody,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export interface ApiOptions {
  /** SvelteKit's fetch, passed in from a load. */
  fetch?: typeof globalThis.fetch;
  /** Query params: null and undefined entries are dropped, the rest are encoded. */
  query?: Record<string, string | number | boolean | null | undefined>;
  /** Aborts the request. The rejection is an AbortError (see isAbortError), never toasted. */
  signal?: AbortSignal;
}

/** Same origin everywhere: the Go server serves the SPA in production, and Vite proxies /api in
 *  development. So there is nothing to configure, no CORS, and no $env. */
const BASE_URL = "/api/v1";

function buildUrl(endpoint: string, query?: ApiOptions["query"]): string {
  const url = `${BASE_URL}${endpoint}`;
  if (!query) return url;
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value !== undefined && value !== null) params.set(key, String(value));
  }
  const qs = params.toString();
  return qs ? `${url}?${qs}` : url;
}

/** Sends the request and turns any non-2xx into an ApiError. */
async function send(method: string, endpoint: string, body: unknown, opts: ApiOptions = {}): Promise<Response> {
  const { fetch: customFetch, query, signal } = opts;

  const response = await (customFetch ?? fetch)(buildUrl(endpoint, query), {
    method,
    headers: body === undefined ? undefined : { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
    credentials: "include",
    signal,
  });

  if (!response.ok) {
    // Raptor's errors are JSON; a proxy's error page is not, and the status still carries the meaning.
    const errorBody: ApiErrorBody = await response.json().catch(() => ({}));
    throw new ApiError(errorBody.message ?? `Request failed (${response.status})`, response.status, errorBody);
  }
  return response;
}

async function request<T>(method: string, endpoint: string, body: unknown, opts?: ApiOptions): Promise<T> {
  const response = await send(method, endpoint, body, opts);
  if (response.status === 204 || response.headers.get("content-length") === "0") {
    return undefined as T;
  }
  return response.json() as Promise<T>;
}

/** True when a rejection is the caller's own abort rather than a failure. */
export const isAbortError = (e: unknown) => e instanceof DOMException && e.name === "AbortError";

export const api = {
  get: <T>(endpoint: string, opts?: ApiOptions) => request<T>("GET", endpoint, undefined, opts),
  post: <T>(endpoint: string, body?: unknown, opts?: ApiOptions) => request<T>("POST", endpoint, body, opts),
};
