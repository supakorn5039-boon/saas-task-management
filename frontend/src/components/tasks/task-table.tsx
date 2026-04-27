import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnMeta,
  type SortingState,
} from "@tanstack/react-table";
import { useMemo } from "react";
import { DataTablePagination } from "@/components/shared/data-table/data-table-pagination";
import { TaskToolbar } from "@/components/tasks/task-toolbar";
import { buildTaskColumns } from "@/components/tasks/task-columns";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type {
  TaskListActions,
  TaskListState,
} from "@/hooks/use-task-list-state";
import type {
  Task,
  TaskListCounts,
  TaskListMeta,
  TaskSortField,
  TaskStatus,
} from "@/types/task";

interface Props {
  tasks: Task[];
  meta: TaskListMeta;
  counts: TaskListCounts;
  state: TaskListState;
  actions: TaskListActions;
  onTaskStatusChange: (id: number, status: TaskStatus) => void;
  onTaskEdit: (task: Task) => void;
  onTaskDelete: (task: Task) => void;
}

// Type-safe accessor for the optional `className` we tag on column meta —
// avoids repeating the `as` cast in both header and cell loops below.
type ColumnClassMeta = ColumnMeta<Task, unknown> & { className?: string };

export function TaskTable({
  tasks,
  meta,
  counts,
  state,
  actions,
  onTaskStatusChange,
  onTaskEdit,
  onTaskDelete,
}: Props) {
  const columns = useMemo(
    () =>
      buildTaskColumns({
        onStatusChange: onTaskStatusChange,
        onEdit: onTaskEdit,
        onDelete: onTaskDelete,
      }),
    [onTaskStatusChange, onTaskEdit, onTaskDelete],
  );

  const sorting: SortingState = [
    { id: state.sort.field, desc: state.sort.order === "desc" },
  ];

  const table = useReactTable({
    data: tasks,
    columns,
    state: { sorting },
    onSortingChange: (updater) => {
      const next = typeof updater === "function" ? updater(sorting) : updater;
      const first = next[0];
      if (!first) return;
      actions.setSort({
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
        statusFilter={state.status}
        search={state.search}
        onStatusFilterChange={actions.setStatus}
        onSearchChange={actions.setSearch}
      />

      <div className="overflow-hidden rounded-lg border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((group) => (
              <TableRow key={group.id}>
                {group.headers.map((header) => {
                  const m = header.column.columnDef.meta as
                    | ColumnClassMeta
                    | undefined;
                  return (
                    <TableHead key={header.id} className={m?.className}>
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
                    const m = cell.column.columnDef.meta as
                      | ColumnClassMeta
                      | undefined;
                    return (
                      <TableCell key={cell.id} className={m?.className}>
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
        onPageChange={actions.setPage}
        onPerPageChange={actions.setPerPage}
      />
    </div>
  );
}
