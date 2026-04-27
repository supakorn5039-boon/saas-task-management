import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnMeta,
  type SortingState,
} from "@tanstack/react-table";
import { Search } from "lucide-react";
import { useMemo } from "react";
import { buildUserColumns } from "@/components/admin/user-columns";
import { DataTablePagination } from "@/components/shared/data-table/data-table-pagination";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type {
  UserListActions,
  UserListState,
} from "@/hooks/use-user-list-state";
import type { UserListMeta, UserSortField } from "@/types/admin";
import type { UserProfile } from "@/types/auth";

interface Props {
  users: UserProfile[];
  meta: UserListMeta;
  state: UserListState;
  actions: UserListActions;
  selfId: number | undefined;
  onUserEdit: (user: UserProfile) => void;
  onUserDelete: (user: UserProfile) => void;
}

type ColumnClassMeta = ColumnMeta<UserProfile, unknown> & {
  className?: string;
};

export function UserTable({
  users,
  meta,
  state,
  actions,
  selfId,
  onUserEdit,
  onUserDelete,
}: Props) {
  const columns = useMemo(
    () =>
      buildUserColumns({
        onEdit: onUserEdit,
        onDelete: onUserDelete,
        selfId,
      }),
    [onUserEdit, onUserDelete, selfId],
  );

  const sorting: SortingState = [
    { id: state.sort.field, desc: state.sort.order === "desc" },
  ];

  const table = useReactTable({
    data: users,
    columns,
    state: { sorting },
    onSortingChange: (updater) => {
      const next = typeof updater === "function" ? updater(sorting) : updater;
      const first = next[0];
      if (!first) return;
      actions.setSort({
        field: first.id as UserSortField,
        order: first.desc ? "desc" : "asc",
      });
    },
    manualSorting: true,
    manualPagination: true,
    pageCount: Math.max(1, Math.ceil(meta.total / meta.perPage)),
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="space-y-4">
      <div className="relative w-full sm:w-80">
        <Search className="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
        <Input
          placeholder="Search by email..."
          value={state.search}
          onChange={(e) => actions.setSearch(e.target.value)}
          className="pl-9"
        />
      </div>

      <div className="overflow-hidden rounded-lg border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((group) => (
              <TableRow key={group.id}>
                {group.headers.map((header) => {
                  const m = header.column.columnDef.meta as
                    | ColumnClassMeta
                    | undefined;
                  return (
                    <TableHead key={header.id} className={m?.className}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="text-muted-foreground py-12 text-center text-sm"
                >
                  No users match your search.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id} className="group">
                  {row.getVisibleCells().map((cell) => {
                    const m = cell.column.columnDef.meta as
                      | ColumnClassMeta
                      | undefined;
                    return (
                      <TableCell key={cell.id} className={m?.className}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </TableCell>
                    );
                  })}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination
        page={meta.page}
        perPage={meta.perPage}
        total={meta.total}
        onPageChange={actions.setPage}
        onPerPageChange={actions.setPerPage}
      />
    </div>
  );
}
