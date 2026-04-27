import { act, renderHook } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { useTaskListState } from "@/hooks/use-task-list-state";

describe("useTaskListState", () => {
  it("starts on page 1 with sensible defaults", () => {
    const { result } = renderHook(() => useTaskListState());
    expect(result.current.page).toBe(1);
    expect(result.current.perPage).toBe(10);
    expect(result.current.search).toBe("");
    expect(result.current.status).toBeUndefined();
    expect(result.current.sort).toEqual({
      field: "created_at",
      order: "desc",
    });
  });

  it("changing the search resets page back to 1", () => {
    const { result } = renderHook(() => useTaskListState());
    act(() => result.current.setPage(5));
    expect(result.current.page).toBe(5);
    act(() => result.current.setSearch("milk"));
    expect(result.current.page).toBe(1);
    expect(result.current.search).toBe("milk");
  });

  it("changing status, sort, perPage all reset page", () => {
    const { result } = renderHook(() => useTaskListState());

    act(() => result.current.setPage(3));
    act(() => result.current.setStatus("done"));
    expect(result.current.page).toBe(1);

    act(() => result.current.setPage(4));
    act(() => result.current.setSort({ field: "title", order: "asc" }));
    expect(result.current.page).toBe(1);

    act(() => result.current.setPage(2));
    act(() => result.current.setPerPage(50));
    expect(result.current.page).toBe(1);
    expect(result.current.perPage).toBe(50);
  });

  it("setPage by itself does not snap to 1", () => {
    const { result } = renderHook(() => useTaskListState());
    act(() => result.current.setPage(7));
    expect(result.current.page).toBe(7);
  });
});
