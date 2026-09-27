/** The index one step from current in a list of length items, stopping at the ends. With nothing
 *  focused (-1), a step forward starts at the first item and a step back at the last. */
export function stepIndex(current: number, length: number, direction: 1 | -1): number {
  if (length === 0) return -1;
  if (current === -1) return direction === 1 ? 0 : length - 1;
  return Math.min(length - 1, Math.max(0, current + direction));
}

/** Keys typed into a field are text, not shortcuts. */
export function isTypingTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false;
  return target.isContentEditable || ["INPUT", "TEXTAREA", "SELECT"].includes(target.tagName);
}

/** Moves focus to the next or previous headline and scrolls it into view. */
export function moveHeadlineFocus(direction: 1 | -1): void {
  const links = [...document.querySelectorAll<HTMLElement>("[data-headline-link]")];
  const row = document.activeElement?.closest("[data-headline]");
  const current = row ? links.findIndex((link) => row.contains(link)) : -1;
  const next = links[stepIndex(current, links.length, direction)];
  if (!next) return;
  next.focus({ preventScroll: true });
  next.closest("[data-headline]")?.scrollIntoView({ block: "nearest" });
}

/** Clicks the story toggle of the focused headline; with onlyIfOpen, only to close it. */
export function toggleFocusedStory(onlyIfOpen = false): void {
  const row = document.activeElement?.closest("[data-headline]");
  const toggle = row?.querySelector<HTMLButtonElement>("[data-story-toggle]");
  if (!toggle || (onlyIfOpen && toggle.getAttribute("aria-expanded") !== "true")) return;
  toggle.click();
}
