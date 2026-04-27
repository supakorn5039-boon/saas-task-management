import { describe, expect, it } from "vitest";
import { taskKeys } from "@/services/task.service";

describe("taskKeys", () => {
  it("hierarchy lets us invalidate everything via taskKeys.all", () => {
    const list = taskKeys.list({ page: 1, perPage: 10 });
    const detail = taskKeys.detail(42);

    // Both list and detail keys must start with the all-key prefix.
    expect(list[0]).toBe("tasks");
    expect(detail[0]).toBe("tasks");
  });

  it("list keys are stable for the same params", () => {
    const a = taskKeys.list({ page: 1, perPage: 10 });
    const b = taskKeys.list({ page: 1, perPage: 10 });
    expect(a).toEqual(b);
  });

  it("list keys differ when params differ", () => {
    const a = taskKeys.list({ page: 1, perPage: 10 });
    const b = taskKeys.list({ page: 2, perPage: 10 });
    expect(a).not.toEqual(b);
  });

  it("detail key includes the id", () => {
    expect(taskKeys.detail(42)).toEqual(["tasks", "detail", 42]);
  });
});
