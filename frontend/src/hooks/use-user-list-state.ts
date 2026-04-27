import { useCallback, useState } from "react";
import type { UserSortField } from "@/types/admin";
import type { SortOrder } from "@/types/task";

export interface UserSort {
  field: UserSortField;
  order: SortOrder;
}

export interface UserListState {
  page: number;
  perPage: number;
  search: string;
  sort: UserSort;
}

export interface UserListActions {
  setPage: (page: number) => void;
  setPerPage: (perPage: number) => void;
  setSearch: (search: string) => void;
  setSort: (sort: UserSort) => void;
}

const DEFAULT_SORT: UserSort = { field: "created_at", order: "desc" };

// Same shape as useTaskListState — anything that affects the result set
// resets the page to 1 so we never land on an empty page.
export function useUserListState(
  initialPerPage = 10,
): UserListState & UserListActions {
  const [page, setPage] = useState(1);
  const [perPage, setPerPageRaw] = useState(initialPerPage);
  const [search, setSearchRaw] = useState("");
  const [sort, setSortRaw] = useState<UserSort>(DEFAULT_SORT);

  const setSearch = useCallback((v: string) => {
    setPage(1);
    setSearchRaw(v);
  }, []);
  const setSort = useCallback((v: UserSort) => {
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
    search,
    sort,
    setPage,
    setPerPage,
    setSearch,
    setSort,
  };
}
