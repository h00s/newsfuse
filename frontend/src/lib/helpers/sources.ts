import type { Source } from "$lib/types/source";

/** Logos are static files named after the source. */
export const sourceLogo = (source: Pick<Source, "name">) => `/img/sources/${encodeURIComponent(source.name)}.webp`;

/** "HN" for Hacker News, "IN" for Index.hr: the fallback when a logo is missing. */
export function sourceInitials(name: string): string {
  const words = name.split(/[\s.]+/).filter(Boolean);
  const letters = words.length > 1 ? words[0][0] + words[1][0] : name.slice(0, 2);
  return letters.toUpperCase();
}
