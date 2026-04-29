import { createLazyFileRoute } from "@tanstack/react-router";
import { AuditFilters } from "@/components/audit/audit-filters";
import { AuditTable } from "@/components/audit/audit-table";
import { PageError } from "@/components/shared/page-state";
import { TableSkeleton } from "@/components/shared/skeletons";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toAuditParams, useAuditListState } from "@/hooks/use-audit-list-state";
import { useAdminAuditLogs } from "@/hooks/use-audit-logs";

export const Route = createLazyFileRoute("/audit-logs")({
  component: AuditLogsPage,
});

function AuditLogsPage() {
  const list = useAuditListState();
  const { data, isLoading, isFetching, error } = useAdminAuditLogs(
    toAuditParams(list),
  );

  if (error) return <PageError label="Error loading audit logs" />;

  return (
    <div className="mx-auto max-w-6xl space-y-6 p-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            Audit logs
            {isFetching && !isLoading && (
              <span className="text-muted-foreground text-xs font-normal">
                Updating…
              </span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <AuditFilters state={list} actions={list} />
          {isLoading ? (
            <TableSkeleton rows={5} columns={6} />
          ) : (
            <AuditTable
              rows={data?.data ?? []}
              meta={
                data?.meta ?? {
                  page: list.page,
                  perPage: list.perPage,
                  total: 0,
                }
              }
              state={list}
              actions={list}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
