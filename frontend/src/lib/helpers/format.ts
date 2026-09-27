const relative = new Intl.RelativeTimeFormat("hr", { numeric: "auto" });
const absolute = new Intl.DateTimeFormat("hr-HR", { dateStyle: "medium", timeStyle: "short" });
const plural = new Intl.PluralRules("hr");

const MINUTE = 60;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;

/** "prije 5 minuta", "prekjučer", or the date once it is more than a week old. */
export function relativeTime(iso: string, now: number): string {
  const seconds = (Date.parse(iso) - now) / 1000; // negative: in the past
  const elapsed = Math.abs(seconds);
  if (elapsed < 45) return "upravo sada";
  if (elapsed < HOUR) return relative.format(Math.round(seconds / MINUTE), "minute");
  if (elapsed < DAY) return relative.format(Math.round(seconds / HOUR), "hour");
  if (elapsed < 7 * DAY) return relative.format(Math.round(seconds / DAY), "day");
  return absoluteTime(iso);
}

/** "27. ruj 2026. 14:05": the exact time, for tooltips. */
export const absoluteTime = (iso: string) => absolute.format(new Date(iso));

const newHeadlines: Record<string, string> = {
  one: "nova vijest",
  few: "nove vijesti",
  other: "novih vijesti",
};

/** "1 nova vijest", "3 nove vijesti", "5 novih vijesti". */
export function newHeadlinesLabel(count: number): string {
  return `${count} ${newHeadlines[plural.select(count)] ?? newHeadlines.other}`;
}
