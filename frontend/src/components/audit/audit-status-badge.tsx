import { CheckCircle2, XCircle } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import type { AuditStatus } from "@/types/audit";

interface Props {
  status: AuditStatus;
}

export function AuditStatusBadge({ status }: Props) {
  if (status === "failure") {
    return (
      <Badge
        variant="outline"
        className="gap-1 border-0 bg-rose-100 font-medium text-rose-700 dark:bg-rose-950 dark:text-rose-300"
      >
        <XCircle className="h-3 w-3" />
        Failure
      </Badge>
    );
  }
  return (
    <Badge
      variant="outline"
      className="gap-1 border-0 bg-emerald-100 font-medium text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300"
    >
      <CheckCircle2 className="h-3 w-3" />
      Success
    </Badge>
  );
}
