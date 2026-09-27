import type { Attachment } from "svelte/attachments";

/** Calls onEnter whenever the element scrolls into view, a little before it does by default. */
export function inView(onEnter: () => void, options: IntersectionObserverInit = { rootMargin: "600px 0px" }): Attachment<HTMLElement> {
  return (node) => {
    const observer = new IntersectionObserver((entries) => {
      if (entries.some((entry) => entry.isIntersecting)) onEnter();
    }, options);
    observer.observe(node);
    return () => observer.disconnect();
  };
}
