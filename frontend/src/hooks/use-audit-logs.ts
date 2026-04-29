import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { auditKeys, auditService } from "@/services/audit.service";
import type { ListAuditParams } from "@/types/audit";

// Admin: every actor's events.
export function useAdminAuditLogs(params: ListAuditParams) {
  return useQuery({
    queryKey: auditKeys.adminLogList(params),
    queryFn: () => auditService.listAdminLogs(params),
    staleTime: 5_000,
    placeholderData: keepPreviousData,
  });
}

// "My activity" — the calling user's own events.
export function useMyActivity(params: ListAuditParams) {
  return useQuery({
    queryKey: auditKeys.myActivityList(params),
    queryFn: () => auditService.listMyActivity(params),
    staleTime: 5_000,
    placeholderData: keepPreviousData,
  });
}
