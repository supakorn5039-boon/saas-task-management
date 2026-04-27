import type { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";
import { DataTableSortHeader } from "@/components/shared/data-table/data-table-sort-header";
import { TaskStatusBadge } from "@/components/tasks/task-status-badge";
import { TaskStatusSelect } from "@/components/tasks/task-status-select";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { formatRelative } from "@/lib/date";
import type { Task, TaskStatus } from "@/types/task";

export interface TaskColumnHandlers {
  onStatusChange: (id: number, status: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
}

export function buildTaskColumns(
  handlers: TaskColumnHandlers,
): ColumnDef<Task>[] {
  return [
    {
      id: "status",
      accessorKey: "status",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Status" />
      ),
      enableSorting: true,
      meta: { className: "w-36" },
      cell: ({ row }) => (
        <TaskStatusBadge status={row.original.status ?? "todo"} />
      ),
    },
    {
      id: "title",
      accessorKey: "title",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Title" />
      ),
      enableSorting: true,
      cell: ({ row }) => {
        const task = row.original;
        const isDone = task.status === "done";
        return (
          <div>
            <div
              className={`font-medium ${isDone ? "text-muted-foreground line-through" : ""}`}
            >
              {task.title}
            </div>
            {task.description && (
              <div className="text-muted-foreground mt-0.5 line-clamp-1 text-sm">
                {task.description}
              </div>
            )}
          </div>
        );
      },
    },
    {
      id: "created_at",
      accessorKey: "createdAt",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Created" />
      ),
      enableSorting: true,
      meta: { className: "w-32" },
      cell: ({ row }) => (
        <span className="text-muted-foreground text-sm">
          {formatRelative(row.original.createdAt)}
        </span>
      ),
    },
    {
      id: "actions",
      header: () => <div className="text-right">Actions</div>,
      enableSorting: false,
      meta: { className: "w-44 text-right" },
      cell: ({ row }) => {
        const task = row.original;
        const status = task.status ?? "todo";
        return (
          <div className="inline-flex items-center gap-1 opacity-60 transition-opacity group-hover:opacity-100">
            <TaskStatusSelect
              value={status}
              onChange={(s) => handlers.onStatusChange(task.id, s)}
              className="h-8 w-28 text-xs"
            />
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handlers.onEdit(task)}
                    className="text-muted-foreground hover:text-foreground h-8 w-8"
                    aria-label="Edit task"
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                }
              />
              <TooltipContent>Edit task</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handlers.onDelete(task)}
                    className="text-muted-foreground hover:text-destructive h-8 w-8"
                    aria-label="Delete task"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                }
              />
              <TooltipContent>Delete task</TooltipContent>
            </Tooltip>
          </div>
        );
      },
    },
  ];
}
