import { Search, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type {
  AuditListActions,
  AuditListState,
} from "@/hooks/use-audit-list-state";
import { AUDIT_ACTIONS, AUDIT_ACTION_LABEL } from "@/types/audit";

interface Props {
  state: AuditListState;
  actions: AuditListActions;
  // The "my activity" page hides the actor-search field — the user is
  // already filtered to themselves at the API level.
  showActorSearch?: boolean;
}

const ALL_ACTIONS_VALUE = "__all__";

export function AuditFilters({
  state,
  actions,
  showActorSearch = true,
}: Props) {
  const hasFilters =
    Boolean(state.action) ||
    Boolean(state.search) ||
    Boolean(state.from) ||
    Boolean(state.to);

  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
      {showActorSearch && (
        <div className="space-y-1.5">
          <Label htmlFor="audit-search" className="text-xs">
            Actor email
          </Label>
          <div className="relative">
            <Search className="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
            <Input
              id="audit-search"
              placeholder="someone@example.com"
              value={state.search}
              onChange={(e) => actions.setSearch(e.target.value)}
              className="pl-9"
            />
          </div>
        </div>
      )}

      <div className="space-y-1.5">
        <Label className="text-xs">Action</Label>
        <Select
          value={state.action || ALL_ACTIONS_VALUE}
          onValueChange={(v) =>
            actions.setAction(v && v !== ALL_ACTIONS_VALUE ? v : "")
          }
        >
          <SelectTrigger className="w-full">
            <SelectValue
              labels={{
                [ALL_ACTIONS_VALUE]: "All actions",
                ...AUDIT_ACTION_LABEL,
              }}
            />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={ALL_ACTIONS_VALUE}>All actions</SelectItem>
            {AUDIT_ACTIONS.map((a) => (
              <SelectItem key={a} value={a}>
                {AUDIT_ACTION_LABEL[a]}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="audit-from" className="text-xs">
          From
        </Label>
        <Input
          id="audit-from"
          type="date"
          value={state.from ? state.from.slice(0, 10) : ""}
          onChange={(e) => {
            const iso = e.target.value
              ? new Date(e.target.value + "T00:00:00").toISOString()
              : "";
            actions.setDateRange(iso, state.to);
          }}
        />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="audit-to" className="text-xs">
          To
        </Label>
        <div className="flex gap-2">
          <Input
            id="audit-to"
            type="date"
            value={state.to ? state.to.slice(0, 10) : ""}
            onChange={(e) => {
              const iso = e.target.value
                ? new Date(e.target.value + "T23:59:59").toISOString()
                : "";
              actions.setDateRange(state.from, iso);
            }}
          />
          {hasFilters && (
            <Button
              variant="outline"
              size="icon"
              onClick={actions.clear}
              aria-label="Clear filters"
              title="Clear filters"
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
