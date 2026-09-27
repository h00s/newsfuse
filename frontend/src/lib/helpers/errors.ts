import { error } from "@sveltejs/kit";
import { toast } from "svelte-sonner";
import { ApiError, isAbortError } from "$lib/api/client";

type Overrides = Partial<Record<number, string>>;

const byStatus: Record<number, string> = {
  400: "Zahtjev nije ispravan.",
  403: "Zahtjev je odbijen.",
  404: "Traženi sadržaj ne postoji.",
  429: "Previše zahtjeva. Pokušajte ponovno za minutu.",
  502: "Izvor trenutno nije dostupan. Pokušajte ponovno kasnije.",
  503: "Poslužitelj trenutno nije dostupan. Pokušajte ponovno kasnije.",
  504: "Izvor se predugo ne javlja. Pokušajte ponovno.",
};

/** Croatian text for any thrown value. The API's own message is English and meant for
 *  developers, so it never reaches the UI. A call site can override any status with a message
 *  that names the actual problem: errorMessage(e, { 404: "Tema ne postoji." }). */
export function errorMessage(e: unknown, overrides: Overrides = {}): string {
  if (e instanceof ApiError) {
    return (
      overrides[e.status] ??
      byStatus[e.status] ??
      (e.status >= 500 ? "Greška na poslužitelju. Pokušajte ponovno." : "Došlo je do pogreške.")
    );
  }
  if (e instanceof TypeError) return "Poslužitelj nije dostupan. Provjerite vezu.";
  return "Došlo je do pogreške.";
}

/** The catch of a user action. An abort gets no toast: the caller asked for it. */
export function toastError(e: unknown, overrides?: Overrides): void {
  if (isAbortError(e)) return;
  toast.error(errorMessage(e, overrides));
}

/** The catch of a load: the API's status becomes the error page's, with Croatian text. A network
 *  failure has no status, so it is a 503. */
export function loadFailed(e: unknown, overrides?: Overrides): never {
  if (e instanceof ApiError) error(e.status, errorMessage(e, overrides));
  if (e instanceof TypeError) error(503, errorMessage(e, overrides));
  throw e;
}
