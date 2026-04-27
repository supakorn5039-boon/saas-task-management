import type { Column } from "@tanstack/react-table";
import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";

interface Props<TData, TValue> {
  column: Column<TData, TValue>;
  label: string;
}

export function DataTableSortHeader<TData, TValue>({
  column,
  label,
}: Props<TData, TValue>) {
  if (!column.getCanSort()) return <span>{label}</span>;

  const sorted = column.getIsSorted();
  const Icon =
    sorted === false ? ArrowUpDown : sorted === "asc" ? ArrowUp : ArrowDown;

  return (
    <button
      type="button"
      onClick={() => column.toggleSorting(sorted === "asc")}
      className="hover:text-foreground inline-flex items-center gap-1.5 transition-colors"
    >
      {label}
      <Icon
        className={`h-3.5 w-3.5 ${sorted ? "text-foreground" : "text-muted-foreground"}`}
      />
    </button>
  );
}
