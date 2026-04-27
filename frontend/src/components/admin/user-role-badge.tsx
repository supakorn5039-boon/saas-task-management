import { Badge } from "@/components/ui/badge";
import { ROLE_LABEL } from "@/types/admin";
import type { Role } from "@/types/auth";

const ROLE_STYLE: Record<Role, string> = {
  admin:
    "bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300",
  manager: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
  user: "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200",
};

interface Props {
  role: Role;
}

export function UserRoleBadge({ role }: Props) {
  return (
    <Badge
      variant="outline"
      className={`border-0 font-medium ${ROLE_STYLE[role]}`}
    >
      {ROLE_LABEL[role]}
    </Badge>
  );
}
