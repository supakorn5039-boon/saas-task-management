import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { TASK_PRIORITY_LABEL } from "@/constants/task";
import { TASK_PRIORITIES, type TaskPriority } from "@/types/task";

interface Props {
  value: TaskPriority;
  onChange: (priority: TaskPriority) => void;
  className?: string;
}

export function TaskPrioritySelect({ value, onChange, className }: Props) {
  return (
    <Select value={value} onValueChange={(v) => onChange(v as TaskPriority)}>
      <SelectTrigger className={className ?? "w-full"}>
        <SelectValue labels={TASK_PRIORITY_LABEL} />
      </SelectTrigger>
      <SelectContent>
        {TASK_PRIORITIES.map((p) => (
          <SelectItem key={p} value={p}>
            {TASK_PRIORITY_LABEL[p]}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
