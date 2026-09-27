export interface Story {
  id: number;
  headlineId: number;
  /** Sanitized HTML: paragraphs, emphasis, lists and safe links. */
  content: string;
  /** Plain text; empty until the story is summarized. */
  summary: string;
}
