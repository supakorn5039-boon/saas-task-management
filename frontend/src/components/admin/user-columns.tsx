import type { ColumnDef } from "@tanstack/react-table";
import { Pencil, Trash2 } from "lucide-react";
import { DataTableSortHeader } from "@/components/shared/data-table/data-table-sort-header";
import { UserRoleBadge } from "@/components/admin/user-role-badge";
import { UserStatusBadge } from "@/components/admin/user-status-badge";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { formatRelative } from "@/lib/date";
import type { UserStatus } from "@/types/admin";
import type { Role, UserProfile } from "@/types/auth";

export interface UserColumnHandlers {
  onEdit: (user: UserProfile) => void;
  onDelete: (user: UserProfile) => void;
  // Hide the action buttons for the row that represents the actor — admins
  // can't edit or delete their own row through this table.
  selfId: number | undefined;
}

export function buildUserColumns(
  handlers: UserColumnHandlers,
): ColumnDef<UserProfile>[] {
  return [
    {
      id: "email",
      accessorKey: "email",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Email" />
      ),
      enableSorting: true,
      cell: ({ row }) => {
        const user = row.original;
        const isSelf = handlers.selfId === user.id;
        return (
          <div className="flex flex-col">
            <span className="font-medium">{user.email}</span>
            {isSelf && (
              <span className="text-muted-foreground text-xs">(you)</span>
            )}
          </div>
        );
      },
    },
    {
      id: "role",
      accessorKey: "role",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Role" />
      ),
      enableSorting: true,
      meta: { className: "w-32" },
      cell: ({ row }) => <UserRoleBadge role={row.original.role as Role} />,
    },
    {
      id: "status",
      accessorKey: "status",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Status" />
      ),
      enableSorting: true,
      meta: { className: "w-32" },
      cell: ({ row }) => (
        <UserStatusBadge status={row.original.status as UserStatus} />
      ),
    },
    {
      id: "created_at",
      accessorKey: "createdAt",
      header: ({ column }) => (
        <DataTableSortHeader column={column} label="Joined" />
      ),
      enableSorting: true,
      meta: { className: "w-32" },
      cell: ({ row }) => (
        <span className="text-muted-foreground text-sm">
          {row.original.createdAt
            ? formatRelative(row.original.createdAt)
            : "—"}
        </span>
      ),
    },
    {
      id: "actions",
      header: () => <div className="text-right">Actions</div>,
      enableSorting: false,
      meta: { className: "w-28 text-right" },
      cell: ({ row }) => {
        const user = row.original;
        const isSelf = handlers.selfId === user.id;
        if (isSelf) {
          return (
            <span className="text-muted-foreground text-xs italic">
              No actions
            </span>
          );
        }
        return (
          <div className="inline-flex items-center gap-1 opacity-60 transition-opacity group-hover:opacity-100">
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handlers.onEdit(user)}
                    className="text-muted-foreground hover:text-foreground h-8 w-8"
                    aria-label="Edit user"
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                }
              />
              <TooltipContent>Edit user</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handlers.onDelete(user)}
                    className="text-muted-foreground hover:text-destructive h-8 w-8"
                    aria-label="Delete user"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                }
              />
              <TooltipContent>Delete user</TooltipContent>
            </Tooltip>
          </div>
        );
      },
    },
  ];
}
