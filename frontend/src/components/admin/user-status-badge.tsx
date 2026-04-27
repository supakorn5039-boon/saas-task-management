import { Badge } from "@/components/ui/badge";
import { USER_STATUS_LABEL, type UserStatus } from "@/types/admin";

const STATUS_STYLE: Record<UserStatus, { dot: string; chip: string }> = {
  1: {
    dot: "bg-emerald-500",
    chip: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300",
  },
  0: {
    dot: "bg-slate-400",
    chip: "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200",
  },
};

interface Props {
  status: UserStatus;
}

export function UserStatusBadge({ status }: Props) {
  const style = STATUS_STYLE[status];
  return (
    <Badge
      variant="outline"
      className={`gap-1.5 border-0 font-medium ${style.chip}`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${style.dot}`} aria-hidden />
      {USER_STATUS_LABEL[status]}
    </Badge>
  );
}
