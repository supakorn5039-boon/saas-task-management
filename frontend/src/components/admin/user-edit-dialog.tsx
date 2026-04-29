import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useUpdateUser } from "@/hooks/use-admin-users";
import { getApiError } from "@/lib/api-error";
import { ROLE_LABEL, ROLES, USER_STATUS_LABEL } from "@/types/admin";
import type { Role, UserProfile } from "@/types/auth";
import {
  editUserSchema,
  type EditUserInput,
} from "@/validators/admin.validator";

interface Props {
  user: UserProfile | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function UserEditDialog({ user, open, onOpenChange }: Props) {
  const updateUser = useUpdateUser();
  const form = useForm<EditUserInput>({
    resolver: zodResolver(editUserSchema),
    defaultValues: { role: "user", status: 1 },
  });

  useEffect(() => {
    if (user && open) {
      form.reset({
        role: (user.role as Role) ?? "user",
        status: user.status === 0 ? 0 : 1,
      });
    }
  }, [user, open, form]);

  if (!user) return null;

  const onSubmit = (values: EditUserInput) => {
    updateUser.mutate(
      { id: user.id, data: values },
      {
        onSuccess: () => {
          toast.success("User updated");
          onOpenChange(false);
        },
        onError: (err) =>
          toast.error(getApiError(err, "Failed to update user")),
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit user</DialogTitle>
          <DialogDescription>{user.email}</DialogDescription>
        </DialogHeader>

        <form
          id="edit-user-form"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-4 pt-2"
        >
          <div className="space-y-2">
            <Label>Role</Label>
            <Controller
              control={form.control}
              name="role"
              render={({ field }) => (
                <Select value={field.value} onValueChange={field.onChange}>
                  <SelectTrigger className="w-full">
                    <SelectValue labels={ROLE_LABEL} />
                  </SelectTrigger>
                  <SelectContent>
                    {ROLES.map((r) => (
                      <SelectItem key={r} value={r}>
                        {ROLE_LABEL[r]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-2">
            <Label>Status</Label>
            <Controller
              control={form.control}
              name="status"
              render={({ field }) => (
                <Select
                  value={String(field.value)}
                  onValueChange={(v) => field.onChange(Number(v) as 0 | 1)}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue labels={USER_STATUS_LABEL} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">{USER_STATUS_LABEL[1]}</SelectItem>
                    <SelectItem value="0">{USER_STATUS_LABEL[0]}</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
          </div>
        </form>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={updateUser.isPending}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            form="edit-user-form"
            disabled={updateUser.isPending}
          >
            {updateUser.isPending ? "Saving..." : "Save changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
