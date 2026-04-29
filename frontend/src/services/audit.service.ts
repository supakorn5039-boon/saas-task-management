import api from "@/lib/axios";
import type { AuditListResponse, ListAuditParams } from "@/types/audit";

// Mirror taskKeys / authKeys / adminKeys — single source of truth for
// invalidation and query identity.
export const auditKeys = {
  all: ["audit"] as const,
  // Admin: every actor's events.
  adminLogs: () => [...auditKeys.all, "admin"] as const,
  adminLogList: (params: ListAuditParams) =>
    [...auditKeys.adminLogs(), "list", params] as const,
  // User: only the calling user's events.
  myActivity: () => [...auditKeys.all, "mine"] as const,
  myActivityList: (params: ListAuditParams) =>
    [...auditKeys.myActivity(), "list", params] as const,
};

const toQuery = (params: ListAuditParams) => ({
  page: params.page,
  per_page: params.perPage,
  action: params.action || undefined,
  search: params.search || undefined,
  from: params.from || undefined,
  to: params.to || undefined,
  sort: params.sort,
  order: params.order,
});

export const auditService = {
  listAdminLogs: async (
    params: ListAuditParams,
  ): Promise<AuditListResponse> => {
    const response = await api.get<AuditListResponse>("/admin/audit-logs", {
      params: toQuery(params),
    });
    return response.data;
  },

  listMyActivity: async (
    params: ListAuditParams,
  ): Promise<AuditListResponse> => {
    const response = await api.get<AuditListResponse>("/user/activity", {
      params: toQuery(params),
    });
    return response.data;
  },

  // Best-effort logout — records the event in the audit log. Failures are
  // swallowed by the caller since the user is logging out anyway.
  logout: async (): Promise<void> => {
    await api.post("/auth/logout");
  },
};
