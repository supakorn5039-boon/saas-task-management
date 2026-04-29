import { AlertCircle, CalendarClock } from "lucide-react";
import type { TaskStatus } from "@/types/task";

interface Props {
  due: string | null | undefined;
  status: TaskStatus;
  className?: string;
}

// Renders a due date as "Due in 3d" / "Overdue 2d" / a date — and turns red
// when the deadline has passed AND the task isn't done yet. Done tasks never
// show as overdue (you delivered, no nag).
export function TaskDueDate({ due, status, className }: Props) {
  if (!due) {
    return (
      <span className={`text-muted-foreground text-sm ${className ?? ""}`}>
        —
      </span>
    );
  }
  const date = new Date(due);
  const ms = date.getTime() - Date.now();
  const days = Math.round(ms / (1000 * 60 * 60 * 24));
  const overdue = ms < 0 && status !== "done";

  const text = formatDelta(days, date);
  const Icon = overdue ? AlertCircle : CalendarClock;

  return (
    <span
      className={`inline-flex items-center gap-1 text-sm ${overdue ? "text-rose-600 dark:text-rose-400 font-medium" : "text-muted-foreground"} ${className ?? ""}`}
      title={date.toLocaleString()}
    >
      <Icon className="h-3.5 w-3.5" />
      {text}
    </span>
  );
}

function formatDelta(days: number, date: Date): string {
  if (days === 0) return "Today";
  if (days === 1) return "Tomorrow";
  if (days === -1) return "Yesterday";
  if (days > 1 && days <= 7) return `Due in ${days}d`;
  if (days < -1 && days >= -7) return `Overdue ${-days}d`;
  return date.toLocaleDateString();
}
