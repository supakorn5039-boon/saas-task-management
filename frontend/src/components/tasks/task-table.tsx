import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnMeta,
  type SortingState,
} from "@tanstack/react-table";
import { Search } from "lucide-react";
import { useMemo } from "react";
import { DataTablePagination } from "@/components/data-table/data-table-pagination";
import { buildTaskColumns } from "@/components/tasks/task-columns";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TASK_STATUS_LABEL, TASK_STATUSES } from "@/constants/task";
import type {
  SortOrder,
  Task,
  TaskListCounts,
  TaskListMeta,
  TaskSortField,
  TaskStatus,
} from "@/types/task";

interface SortState {
  field: TaskSortField;
  order: SortOrder;
}

interface Props {
  tasks: Task[];
  meta: TaskListMeta;
  counts: TaskListCounts;
  statusFilter: TaskStatus | undefined;
  search: string;
  sort: SortState;
  onStatusFilterChange: (status: TaskStatus | undefined) => void;
  onSearchChange: (search: string) => void;
  onSortChange: (sort: SortState) => void;
  onPageChange: (page: number) => void;
  onPerPageChange: (perPage: number) => void;
  onTaskStatusChange: (id: number, status: TaskStatus) => void;
  onTaskDelete: (id: number) => void;
}

export function TaskTable({
  tasks,
  meta,
  counts,
  statusFilter,
  search,
  sort,
  onStatusFilterChange,
  onSearchChange,
  onSortChange,
  onPageChange,
  onPerPageChange,
  onTaskStatusChange,
  onTaskDelete,
}: Props) {
  const columns = useMemo(
    () =>
      buildTaskColumns({
        onStatusChange: onTaskStatusChange,
        onDelete: onTaskDelete,
      }),
    [onTaskStatusChange, onTaskDelete],
  );

  const sorting: SortingState = [
    { id: sort.field, desc: sort.order === "desc" },
  ];

  const table = useReactTable({
    data: tasks,
    columns,
    state: { sorting },
    onSortingChange: (updater) => {
      const next = typeof updater === "function" ? updater(sorting) : updater;
      const first = next[0];
      if (!first) return;
      onSortChange({
        field: first.id as TaskSortField,
        order: first.desc ? "desc" : "asc",
      });
    },
    manualSorting: true,
    manualPagination: true,
    pageCount: Math.max(1, Math.ceil(meta.total / meta.perPage)),
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="space-y-4">
      <TaskToolbar
        counts={counts}
        statusFilter={statusFilter}
        search={search}
        onStatusFilterChange={onStatusFilterChange}
        onSearchChange={onSearchChange}
      />

      <div className="overflow-hidden rounded-lg border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((group) => (
              <TableRow key={group.id}>
                {group.headers.map((header) => {
                  const meta = header.column.columnDef.meta as
                    | (ColumnMeta<Task, unknown> & { className?: string })
                    | undefined;
                  return (
                    <TableHead key={header.id} className={meta?.className}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="text-muted-foreground py-12 text-center text-sm"
                >
                  No tasks match your filters.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id} className="group">
                  {row.getVisibleCells().map((cell) => {
                    const meta = cell.column.columnDef.meta as
                      | (ColumnMeta<Task, unknown> & { className?: string })
                      | undefined;
                    return (
                      <TableCell key={cell.id} className={meta?.className}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </TableCell>
                    );
                  })}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination
        page={meta.page}
        perPage={meta.perPage}
        total={meta.total}
        onPageChange={onPageChange}
        onPerPageChange={onPerPageChange}
      />
    </div>
  );
}

function TaskToolbar({
  counts,
  statusFilter,
  search,
  onStatusFilterChange,
  onSearchChange,
}: {
  counts: TaskListCounts;
  statusFilter: TaskStatus | undefined;
  search: string;
  onStatusFilterChange: (status: TaskStatus | undefined) => void;
  onSearchChange: (search: string) => void;
}) {
  const value = statusFilter ?? "all";
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <Tabs
        value={value}
        onValueChange={(v) =>
          onStatusFilterChange(v === "all" ? undefined : (v as TaskStatus))
        }
      >
        <TabsList>
          <TabsTrigger value="all">
            All{" "}
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
