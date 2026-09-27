/** Smooth scrolling, unless the reader asked the system for reduced motion. */
export const scrollBehavior = (): ScrollBehavior =>
  matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth";
