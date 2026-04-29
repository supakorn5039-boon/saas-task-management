import { createLazyFileRoute } from "@tanstack/react-router";
import { AuditFilters } from "@/components/audit/audit-filters";
import { AuditTable } from "@/components/audit/audit-table";
import { PageError } from "@/components/shared/page-state";
import { TableSkeleton } from "@/components/shared/skeletons";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { toAuditParams, useAuditListState } from "@/hooks/use-audit-list-state";
import { useMyActivity } from "@/hooks/use-audit-logs";

export const Route = createLazyFileRoute("/activity")({
  component: ActivityPage,
});

function ActivityPage() {
  const list = useAuditListState();
  const { data, isLoading, isFetching, error } = useMyActivity(
    toAuditParams(list),
  );

  if (error) return <PageError label="Error loading activity" />;

  return (
    <div className="mx-auto max-w-6xl space-y-6 p-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            My activity
            {isFetching && !isLoading && (
              <span className="text-muted-foreground text-xs font-normal">
                Updating…
              </span>
            )}
          </CardTitle>
          <CardDescription>
            A history of every action you've taken in the app.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* No actor-search field — every event is the calling user. */}
          <AuditFilters state={list} actions={list} showActorSearch={false} />
          {isLoading ? (
            <TableSkeleton rows={5} columns={5} />
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
              showActor={false}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
