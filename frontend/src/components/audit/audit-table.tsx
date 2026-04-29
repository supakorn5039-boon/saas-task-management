import { ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";
import { AuditActionBadge } from "@/components/audit/audit-action-badge";
import { AuditStatusBadge } from "@/components/audit/audit-status-badge";
import { DataTablePagination } from "@/components/shared/data-table/data-table-pagination";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { formatRelative } from "@/lib/date";
import type {
  AuditListActions,
  AuditListState,
} from "@/hooks/use-audit-list-state";
import type { AuditLog, AuditListMeta, AuditStatus } from "@/types/audit";

interface Props {
  rows: AuditLog[];
  meta: AuditListMeta;
  state: AuditListState;
  actions: AuditListActions;
  // The "my activity" view collapses the actor column — every row is
  // the calling user, so showing it adds nothing.
  showActor?: boolean;
}

export function AuditTable({
  rows,
  meta,
  state,
  actions,
  showActor = true,
}: Props) {
  return (
    <div className="space-y-4">
      {/* Mobile: stacked cards. Desktop: table. */}
      <div className="md:hidden">
        <AuditCardList rows={rows} showActor={showActor} />
      </div>
      <div className="hidden md:block">
        <DesktopTable rows={rows} showActor={showActor} />
      </div>

      <DataTablePagination
        page={meta.page}
        perPage={meta.perPage}
        total={meta.total}
        onPageChange={actions.setPage}
        onPerPageChange={actions.setPerPage}
      />

      {/* Empty state — render after pagination so total/perPage selectors
          stay visible. */}
      {rows.length === 0 && (
        <p className="text-muted-foreground py-12 text-center text-sm">
          No audit events match your filters.
          {state.action || state.search || state.from || state.to ? (
            <>
              {" "}
              <button
                onClick={actions.clear}
                className="text-primary underline"
              >
                Clear filters
              </button>
            </>
          ) : null}
        </p>
      )}
    </div>
  );
}

// ----- Desktop table -----

function DesktopTable({
  rows,
  showActor,
}: {
  rows: AuditLog[];
  showActor: boolean;
}) {
  return (
    <div className="overflow-hidden rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-32">When</TableHead>
            <TableHead className="w-44">Action</TableHead>
            <TableHead className="w-28">Status</TableHead>
            {showActor && <TableHead>Actor</TableHead>}
            <TableHead>Target / details</TableHead>
            <TableHead className="w-32">IP</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row) => (
            <TableRow key={row.id} className="group">
              <TableCell
                className="text-muted-foreground text-sm"
                title={row.createdAt}
              >
                {formatRelative(row.createdAt)}
              </TableCell>
              <TableCell>
                <AuditActionBadge action={row.action} />
              </TableCell>
              <TableCell>
                <AuditStatusBadge status={row.status as AuditStatus} />
              </TableCell>
              {showActor && (
                <TableCell className="max-w-xs truncate text-sm">
                  {row.actorEmail || (
                    <span className="text-muted-foreground italic">
                      anonymous
                    </span>
                  )}
                </TableCell>
              )}
              <TableCell className="text-sm">
                <TargetCell log={row} />
              </TableCell>
              <TableCell className="text-muted-foreground font-mono text-xs">
                {row.ip || "—"}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

// ----- Mobile card list -----

function AuditCardList({
  rows,
  showActor,
}: {
  rows: AuditLog[];
  showActor: boolean;
}) {
  if (rows.length === 0) return null;
  return (
    <ul className="divide-y rounded-lg border">
      {rows.map((row) => (
        <li key={row.id} className="space-y-2 p-3">
          <div className="flex items-center justify-between gap-2">
            <AuditActionBadge action={row.action} />
            <AuditStatusBadge status={row.status as AuditStatus} />
          </div>
          {showActor && row.actorEmail && (
            <p className="truncate text-sm font-medium">{row.actorEmail}</p>
          )}
          <div className="text-muted-foreground text-xs">
            <TargetCell log={row} />
          </div>
          <div className="text-muted-foreground flex items-center justify-between text-xs">
            <span title={row.createdAt}>{formatRelative(row.createdAt)}</span>
            {row.ip && <span className="font-mono">{row.ip}</span>}
          </div>
        </li>
      ))}
    </ul>
  );
}

// ----- Target / metadata cell -----

// Renders a short summary plus an expandable JSON blob if there's metadata.
// Keeps the row scannable while still letting power users dig in.
function TargetCell({ log }: { log: AuditLog }) {
  const [open, setOpen] = useState(false);
  const summary = summarize(log);
  const hasMeta = log.metadata && Object.keys(log.metadata).length > 0;

  return (
    <div className="space-y-1">
      <div className="flex items-center gap-1">
        <span className="truncate">{summary || "—"}</span>
        {hasMeta && (
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setOpen((o) => !o)}
            aria-label={open ? "Hide details" : "Show details"}
            className="h-6 w-6 shrink-0"
          >
            {open ? (
              <ChevronDown className="h-3 w-3" />
            ) : (
              <ChevronRight className="h-3 w-3" />
            )}
          </Button>
        )}
      </div>
      {open && hasMeta && (
        <pre className="bg-muted overflow-x-auto rounded p-2 font-mono text-[11px] leading-tight">
          {JSON.stringify(log.metadata, null, 2)}
        </pre>
      )}
    </div>
  );
}

// summarize picks the most useful text per action so the row tells a story
// without expansion. Falls back to "<targetType> #<id>" when there's nothing
// nicer to say.
function summarize(log: AuditLog): string {
  const meta = log.metadata ?? {};
  switch (log.action) {
    case "auth.login":
    case "auth.register":
    case "auth.logout":
      return "—";
    case "auth.login_failed":
      return typeof meta.reason === "string" ? meta.reason : "bad credentials";
    case "user.password_changed":
      return "Password updated";
    case "admin.user_updated":
      return formatUserUpdate(meta);
    case "admin.user_deleted":
      return typeof meta.targetEmail === "string"
        ? `Deleted ${meta.targetEmail}`
        : `Deleted user #${log.targetId ?? "?"}`;
    case "task.created":
      return typeof meta.title === "string"
        ? meta.title
        : `Task #${log.targetId}`;
    case "task.updated":
      return typeof meta.title === "string"
        ? `${meta.title}${meta.status ? ` → ${meta.status}` : ""}`
        : `Task #${log.targetId}`;
    case "task.deleted":
      return `Task #${log.targetId}`;
    default:
      if (log.targetType && log.targetId)
        return `${log.targetType} #${log.targetId}`;
      return "";
  }
}

function formatUserUpdate(meta: Record<string, unknown>): string {
  const parts: string[] = [];
  if (typeof meta.targetEmail === "string") parts.push(meta.targetEmail);
  if (typeof meta.role === "string") parts.push(`role=${meta.role}`);
  if (typeof meta.status === "number")
    parts.push(`status=${meta.status === 1 ? "active" : "inactive"}`);
  return parts.join(" · ");
}
