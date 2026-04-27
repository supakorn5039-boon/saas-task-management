import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { formatRelative } from "@/lib/date";

describe("formatRelative", () => {
  beforeEach(() => {
    // Pin "now" so the relative buckets are deterministic.
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-04-27T12:00:00Z"));
  });
  afterEach(() => vi.useRealTimers());

  it("under 1 minute → just now", () => {
    expect(formatRelative("2026-04-27T11:59:30Z")).toBe("just now");
  });

  it("minutes bucket", () => {
    expect(formatRelative("2026-04-27T11:55:00Z")).toBe("5m ago");
  });

  it("hours bucket", () => {
    expect(formatRelative("2026-04-27T09:00:00Z")).toBe("3h ago");
  });

  it("days bucket", () => {
    expect(formatRelative("2026-04-25T12:00:00Z")).toBe("2d ago");
  });

  it("over a week → date string", () => {
    const result = formatRelative("2026-04-01T12:00:00Z");
    // toLocaleDateString output varies by locale — just assert it's not a relative string.
    expect(result).not.toMatch(/ago|just now/);
  });
});
