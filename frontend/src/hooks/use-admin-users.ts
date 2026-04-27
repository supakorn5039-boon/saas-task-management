import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { adminKeys, adminService } from "@/services/admin.service";
import type { ListUsersParams, UpdateUserRequest } from "@/types/admin";

function useInvalidateUsers() {
  const qc = useQueryClient();
  return () => qc.invalidateQueries({ queryKey: adminKeys.users() });
}

export function useAdminUsers(params: ListUsersParams) {
  return useQuery({
    queryKey: adminKeys.userList(params),
    queryFn: () => adminService.listUsers(params),
    staleTime: 5_000,
    placeholderData: keepPreviousData,
  });
}

export function useUpdateUser() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) =>
      adminService.updateUser(id, data),
    onSuccess: invalidate,
  });
}

export function useDeleteUser() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: (id: number) => adminService.deleteUser(id),
    onSuccess: invalidate,
  });
}
