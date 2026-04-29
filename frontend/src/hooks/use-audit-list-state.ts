import { useCallback, useState } from "react";
import type { ListAuditParams } from "@/types/audit";

// State container for the audit log filters/pagination. Mirrors the shape
// of useTaskListState / useUserListState — anything that affects the result
// set resets the page to 1.
export interface AuditFilters {
  action: string;
  search: string;
  from: string;
  to: string;
}

export interface AuditListState extends AuditFilters {
  page: number;
  perPage: number;
}

export interface AuditListActions {
  setPage: (page: number) => void;
  setPerPage: (perPage: number) => void;
  setAction: (action: string) => void;
  setSearch: (search: string) => void;
  setDateRange: (from: string, to: string) => void;
  clear: () => void;
}

const EMPTY: AuditFilters = { action: "", search: "", from: "", to: "" };

export function useAuditListState(
  initialPerPage = 20,
): AuditListState & AuditListActions {
  const [page, setPage] = useState(1);
  const [perPage, setPerPageRaw] = useState(initialPerPage);
  const [filters, setFilters] = useState<AuditFilters>(EMPTY);

  const setAction = useCallback((v: string) => {
    setPage(1);
    setFilters((f) => ({ ...f, action: v }));
  }, []);
  const setSearch = useCallback((v: string) => {
    setPage(1);
    setFilters((f) => ({ ...f, search: v }));
  }, []);
  const setDateRange = useCallback((from: string, to: string) => {
    setPage(1);
    setFilters((f) => ({ ...f, from, to }));
  }, []);
  const setPerPage = useCallback((v: number) => {
    setPage(1);
    setPerPageRaw(v);
  }, []);
  const clear = useCallback(() => {
    setPage(1);
    setFilters(EMPTY);
  }, []);

  return {
    page,
    perPage,
    ...filters,
    setPage,
    setPerPage,
    setAction,
    setSearch,
    setDateRange,
    clear,
  };
}

// Build the query params accepted by the backend from the state.
export function toAuditParams(state: AuditListState): ListAuditParams {
  return {
    page: state.page,
    perPage: state.perPage,
    action: state.action || undefined,
    search: state.search || undefined,
    from: state.from || undefined,
    to: state.to || undefined,
  };
}
