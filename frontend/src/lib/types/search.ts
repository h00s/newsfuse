import { z } from "zod";

/** Mirrors the backend's rule for ?query=: 3 to 100 characters after trimming. */
export const searchInput = z.object({
  query: z.string().trim().min(3, "Upišite barem 3 znaka").max(100, "Najviše 100 znakova"),
});
export type SearchInput = z.infer<typeof searchInput>;
