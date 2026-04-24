import { Badge } from "@/components/ui/badge";
import { TASK_STATUS_LABEL } from "@/constants/task";
import { statusStyle } from "@/styles/task-status";
import type { TaskStatus } from "@/types/task";

interface Props {
  status: TaskStatus;
  className?: string;
}

export function TaskStatusBadge({ status, className }: Props) {
  const style = statusStyle(status);
  return (
    <Badge
      variant="outline"
      className={`gap-1.5 border-0 font-medium ${style.badge} ${className ?? ""}`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${style.dot}`} aria-hidden />
      {TASK_STATUS_LABEL[status]}
    </Badge>
  );
}
