/** A news site. Seeded reference data, like topics. */
export interface Source {
  id: number;
  name: string;
  topicId: number;
  /** Its stories can be fetched and summarized; otherwise the headline only links out. */
  isScrapable: boolean;
}
