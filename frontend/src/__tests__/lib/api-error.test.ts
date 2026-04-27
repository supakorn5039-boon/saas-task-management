import { AxiosError } from "axios";
import { describe, expect, it } from "vitest";
import { getApiError } from "@/lib/api-error";

describe("getApiError", () => {
  it("extracts the backend error envelope", () => {
    const err = new AxiosError(
      "Request failed",
      "ERR_BAD_REQUEST",
      undefined,
      undefined,
      // @ts-expect-error — partial response is fine for the test
      { status: 400, data: { error: "title cannot be empty" } },
    );
    expect(getApiError(err)).toBe("title cannot be empty");
  });

  it("falls back to axios error.message when no envelope", () => {
    const err = new AxiosError("Network Error");
    expect(getApiError(err)).toBe("Network Error");
  });

  it("uses fallback for non-axios errors", () => {
    expect(getApiError("oops", "default msg")).toBe("default msg");
  });

  it("respects fallback when error has no message", () => {
    expect(getApiError(undefined, "Failed to load")).toBe("Failed to load");
  });
});
