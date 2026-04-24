import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { TASK_STATUS_LABEL, TASK_STATUSES } from "@/constants/task";
import { TASK_STATUS_STYLE } from "@/styles/task-status";
import type { TaskStatus } from "@/types/task";

interface Props {
  value: TaskStatus;
  onChange: (status: TaskStatus) => void;
  className?: string;
}

export function TaskStatusSelect({ value, onChange, className }: Props) {
  return (
    <Select value={value} onValueChange={(v) => onChange(v as TaskStatus)}>
      <SelectTrigger className={className ?? "w-36"}>
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {TASK_STATUSES.map((status) => (
          <SelectItem key={status} value={status}>
            <span className="flex items-center gap-2">
              <span
                className={`h-2 w-2 rounded-full ${TASK_STATUS_STYLE[status].dot}`}
                aria-hidden
              />
              {TASK_STATUS_LABEL[status]}
            </span>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
