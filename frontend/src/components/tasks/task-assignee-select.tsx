import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAssignableUsers } from "@/hooks/use-tasks";

interface Props {
  value: number | undefined;
  onChange: (assigneeId: number | undefined) => void;
  className?: string;
  // Hide the "Unassigned" option — used in the Create dialog where the
  // backend defaults to the creator if nothing is sent.
  allowUnassigned?: boolean;
}

const UNASSIGNED_VALUE = "__unassigned__";

export function TaskAssigneeSelect({
  value,
  onChange,
  className,
  allowUnassigned = true,
}: Props) {
  const { data, isLoading } = useAssignableUsers();
  const users = data?.data ?? [];

  // Build a value→email label map so the trigger renders the email instead
  // of the raw user id (base-ui Select default).
  const labels: Record<string, string> = {
    [UNASSIGNED_VALUE]: "Unassigned",
  };
  for (const u of users) labels[String(u.id)] = u.email;

  const stringValue = value == null ? UNASSIGNED_VALUE : String(value);

  return (
    <Select
      value={stringValue}
      onValueChange={(v) => {
        if (v === UNASSIGNED_VALUE) {
          onChange(undefined);
        } else {
          onChange(Number(v));
        }
      }}
      disabled={isLoading}
    >
      <SelectTrigger className={className ?? "w-full"}>
        <SelectValue labels={labels} />
      </SelectTrigger>
      <SelectContent>
        {allowUnassigned && (
          <SelectItem value={UNASSIGNED_VALUE}>
            <span className="text-muted-foreground italic">Unassigned</span>
          </SelectItem>
        )}
        {users.map((u) => (
          <SelectItem key={u.id} value={String(u.id)}>
            {u.email}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
