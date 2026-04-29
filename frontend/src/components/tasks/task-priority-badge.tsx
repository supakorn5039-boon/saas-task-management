import { ArrowDown, ArrowUp, Equal, Flame } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { TASK_PRIORITY_LABEL } from "@/constants/task";
import type { TaskPriority } from "@/types/task";

const STYLE: Record<TaskPriority, { chip: string; icon: LucideIcon }> = {
  low: {
    chip: "bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-300",
    icon: ArrowDown,
  },
  medium: {
    chip: "bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300",
    icon: Equal,
  },
  high: {
    chip: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
    icon: ArrowUp,
  },
  urgent: {
    chip: "bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300",
    icon: Flame,
  },
};

interface Props {
  priority: TaskPriority;
  className?: string;
}

export function TaskPriorityBadge({ priority, className }: Props) {
  const cfg = STYLE[priority];
  const Icon = cfg.icon;
  return (
    <Badge
      variant="outline"
      className={`gap-1 border-0 font-medium ${cfg.chip} ${className ?? ""}`}
    >
      <Icon className="h-3 w-3" />
      {TASK_PRIORITY_LABEL[priority]}
    </Badge>
  );
}
