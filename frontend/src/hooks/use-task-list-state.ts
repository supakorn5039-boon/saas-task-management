import { useCallback, useState } from "react";
import type { SortOrder, TaskSortField, TaskStatus } from "@/types/task";

export interface TaskSort {
  field: TaskSortField;
  order: SortOrder;
}

export interface TaskListState {
  page: number;
  perPage: number;
  status: TaskStatus | undefined;
  search: string;
  sort: TaskSort;
}

export interface TaskListActions {
  setPage: (page: number) => void;
  setPerPage: (perPage: number) => void;
  setStatus: (status: TaskStatus | undefined) => void;
  setSearch: (search: string) => void;
  setSort: (sort: TaskSort) => void;
}

const DEFAULT_SORT: TaskSort = { field: "created_at", order: "desc" };

// Owns every piece of server-driven list state for the tasks table: pagination,
// status tab, search, sort. Anything that changes the result set (status / search
// / sort / perPage) snaps page back to 1, so the user never lands on an empty page.
export function useTaskListState(
  initialPerPage = 10,
): TaskListState & TaskListActions {
  const [page, setPage] = useState(1);
  const [perPage, setPerPageRaw] = useState(initialPerPage);
  const [status, setStatusRaw] = useState<TaskStatus | undefined>(undefined);
  const [search, setSearchRaw] = useState("");
  const [sort, setSortRaw] = useState<TaskSort>(DEFAULT_SORT);

  const setStatus = useCallback((v: TaskStatus | undefined) => {
    setPage(1);
    setStatusRaw(v);
  }, []);
  const setSearch = useCallback((v: string) => {
    setPage(1);
    setSearchRaw(v);
  }, []);
  const setSort = useCallback((v: TaskSort) => {
    setPage(1);
    setSortRaw(v);
  }, []);
  const setPerPage = useCallback((v: number) => {
    setPage(1);
    setPerPageRaw(v);
  }, []);

  return {
    page,
    perPage,
    status,
    search,
    sort,
    setPage,
    setPerPage,
    setStatus,
    setSearch,
    setSort,
  };
}
