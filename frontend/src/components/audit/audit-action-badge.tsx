import { Badge } from "@/components/ui/badge";
import { AUDIT_ACTION_LABEL, type AuditAction } from "@/types/audit";

// Color action by domain — auth/user/admin/task. Falls back to a neutral
// style for anything unrecognized so the table never blows up if the backend
// adds a new action before the frontend ships its constant.
const STYLES: Record<string, string> = {
  auth: "bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300",
  user: "bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-300",
  admin: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
  task: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300",
};
const FALLBACK =
  "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200";

interface Props {
  action: string;
}

export function AuditActionBadge({ action }: Props) {
  const domain = action.split(".")[0] ?? "";
  const style = STYLES[domain] ?? FALLBACK;
  const label =
    AUDIT_ACTION_LABEL[action as AuditAction] ?? action.replace(".", " · ");
  return (
    <Badge variant="outline" className={`border-0 font-medium ${style}`}>
      {label}
    </Badge>
  );
}
