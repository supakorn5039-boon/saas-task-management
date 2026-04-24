import type { TaskStatus } from "@/types/task";

export interface TaskStatusStyle {
  dot: string;
  badge: string;
  row: string;
}

export function statusStyle(status: TaskStatus | undefined): TaskStatusStyle {
  return TASK_STATUS_STYLE[status ?? "todo"] ?? TASK_STATUS_STYLE.todo;
}

export const TASK_STATUS_STYLE: Record<TaskStatus, TaskStatusStyle> = {
  todo: {
    dot: "bg-slate-400",
    badge: "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200",
    row: "border-l-slate-300 dark:border-l-slate-700",
  },
  in_progress: {
    dot: "bg-amber-500",
    badge: "bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-200",
    row: "border-l-amber-400",
  },
  done: {
    dot: "bg-emerald-500",
    badge:
      "bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-200",
    row: "border-l-emerald-400",
  },
};
