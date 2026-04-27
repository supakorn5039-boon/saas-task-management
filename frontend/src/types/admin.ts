import type { Role, UserProfile } from "@/types/auth";

export const ROLES: readonly Role[] = ["admin", "manager", "user"];

export const ROLE_LABEL: Record<Role, string> = {
  admin: "Admin",
  manager: "Manager",
  user: "User",
};

// Backend stores user.status as int. 1 = active, 0 = inactive (per the seeder).
export type UserStatus = 0 | 1;

export const USER_STATUS_LABEL: Record<UserStatus, string> = {
  1: "Active",
  0: "Inactive",
};

export type UserSortField = "created_at" | "email" | "role" | "status";

export interface ListUsersParams {
  page?: number;
  perPage?: number;
  search?: string;
  sort?: UserSortField;
  order?: "asc" | "desc";
}

export interface UserListMeta {
  page: number;
  perPage: number;
  total: number;
}

export interface UserListResponse {
  data: UserProfile[];
  meta: UserListMeta;
}

export interface UpdateUserRequest {
  role?: Role;
  status?: UserStatus;
}
