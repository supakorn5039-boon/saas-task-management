import { createLazyFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { toast } from "sonner";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";
import { PageError } from "@/components/shared/page-state";
import { TableSkeleton } from "@/components/shared/skeletons";
import { UserEditDialog } from "@/components/admin/user-edit-dialog";
import { UserTable } from "@/components/admin/user-table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAdminUsers, useDeleteUser } from "@/hooks/use-admin-users";
import { useUserListState } from "@/hooks/use-user-list-state";
import { getApiError } from "@/lib/api-error";
import { selectUser, useAuthStore } from "@/store/auth.store";
import type { UserProfile } from "@/types/auth";

export const Route = createLazyFileRoute("/users")({
  component: UsersPage,
});

function UsersPage() {
  const list = useUserListState();
  const me = useAuthStore(selectUser);

  const { data, isLoading, isFetching, error } = useAdminUsers({
    page: list.page,
    perPage: list.perPage,
    search: list.search,
    sort: list.sort.field,
    order: list.sort.order,
  });

  const deleteUser = useDeleteUser();

  const [editing, setEditing] = useState<UserProfile | null>(null);
  const [deleting, setDeleting] = useState<UserProfile | null>(null);

  const onConfirmDelete = async () => {
    if (!deleting) return;
    await new Promise<void>((resolve) => {
      deleteUser.mutate(deleting.id, {
        onSuccess: () => {
          toast.success("User deleted");
          resolve();
        },
        onError: (err) => {
          toast.error(getApiError(err, "Failed to delete user"));
          resolve();
        },
      });
    });
  };

  if (error) return <PageError label="Error loading users" />;

  return (
    <div className="mx-auto max-w-6xl space-y-6 p-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            Users
            {isFetching && !isLoading && (
              <span className="text-muted-foreground text-xs font-normal">
                Updating…
              </span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <TableSkeleton rows={5} columns={5} />
          ) : (
            <UserTable
              users={data?.data ?? []}
              meta={
                data?.meta ?? {
                  page: list.page,
                  perPage: list.perPage,
                  total: 0,
                }
              }
              state={list}
              actions={list}
              selfId={me?.id}
              onUserEdit={setEditing}
              onUserDelete={setDeleting}
            />
          )}
        </CardContent>
      </Card>

      <UserEditDialog
        user={editing}
        open={editing !== null}
        onOpenChange={(open) => !open && setEditing(null)}
      />

      <ConfirmDialog
        open={deleting !== null}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Delete this user?"
        description={
          deleting
            ? `${deleting.email} will be removed. They will lose access to the app immediately.`
            : undefined
        }
        confirmLabel="Delete"
        destructive
        onConfirm={onConfirmDelete}
      />
    </div>
  );
}
