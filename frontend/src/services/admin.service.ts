import api from "@/lib/axios";
import type { UserProfile } from "@/types/auth";
import type {
  ListUsersParams,
  UpdateUserRequest,
  UserListResponse,
} from "@/types/admin";

// Mirror taskKeys / authKeys — single source of truth for invalidation.
export const adminKeys = {
  all: ["admin"] as const,
  users: () => [...adminKeys.all, "users"] as const,
  userList: (params: ListUsersParams) =>
    [...adminKeys.users(), "list", params] as const,
};

export const adminService = {
  listUsers: async (params: ListUsersParams): Promise<UserListResponse> => {
    const response = await api.get<UserListResponse>("/admin/users", {
      params: {
        page: params.page,
        per_page: params.perPage,
        search: params.search || undefined,
        sort: params.sort,
        order: params.order,
      },
    });
    return response.data;
  },

  updateUser: async (
    id: number,
    data: UpdateUserRequest,
  ): Promise<UserProfile> => {
    const response = await api.put<UserProfile>(`/admin/users/${id}`, data);
    return response.data;
  },

  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/admin/users/${id}`);
  },
};
