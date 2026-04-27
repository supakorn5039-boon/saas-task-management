import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TASK_STATUS_LABEL, TASK_STATUSES } from "@/constants/task";
import type { TaskListCounts, TaskStatus } from "@/types/task";

interface Props {
  counts: TaskListCounts;
  statusFilter: TaskStatus | undefined;
  search: string;
  onStatusFilterChange: (status: TaskStatus | undefined) => void;
  onSearchChange: (search: string) => void;
}

export function TaskToolbar({
  counts,
  statusFilter,
  search,
  onStatusFilterChange,
  onSearchChange,
}: Props) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <Tabs
        value={statusFilter ?? "all"}
        onValueChange={(v) =>
          onStatusFilterChange(v === "all" ? undefined : (v as TaskStatus))
        }
      >
        <TabsList>
          <TabsTrigger value="all">
            All
            <span className="text-muted-foreground ml-1.5">{counts.all}</span>
          </TabsTrigger>
          {TASK_STATUSES.map((s) => (
            <TabsTrigger key={s} value={s}>
              {TASK_STATUS_LABEL[s]}
              <span className="text-muted-foreground ml-1.5">{counts[s]}</span>
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="relative w-full sm:w-64">
        <Search className="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
        <Input
          placeholder="Search tasks..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          className="pl-9"
        />
      </div>
    </div>
  );
}
