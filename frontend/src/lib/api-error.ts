import { isAxiosError } from "axios";

// Backend always responds with `{ error: string }` for non-2xx (see
// backend/src/apiwebserver/controller/response.go). This helper unwraps that
// shape and falls back to a generic message — every route should call this
// instead of poking at error.response?.data?.error directly.
export function getApiError(
  error: unknown,
  fallback = "Something went wrong",
): string {
  if (isAxiosError(error)) {
    const data = error.response?.data as { error?: unknown } | undefined;
    if (data && typeof data.error === "string" && data.error.length > 0) {
      return data.error;
    }
    if (error.message) return error.message;
  }
  if (error instanceof Error && error.message) return error.message;
  return fallback;
}
