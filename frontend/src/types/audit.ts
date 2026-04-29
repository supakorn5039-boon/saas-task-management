export const AUDIT_ACTIONS = [
  "auth.login",
  "auth.login_failed",
  "auth.register",
  "auth.logout",
  "user.password_changed",
  "admin.user_updated",
  "admin.user_deleted",
  "task.created",
  "task.updated",
  "task.deleted",
] as const;

export type AuditAction = (typeof AUDIT_ACTIONS)[number];

// Human-readable labels for filter dropdowns and table cells. Keep aligned
// with the backend AuditAction* constants in database/model/audit.go.
export const AUDIT_ACTION_LABEL: Record<AuditAction, string> = {
  "auth.login": "Login",
  "auth.login_failed": "Failed login",
  "auth.register": "Registration",
  "auth.logout": "Logout",
  "user.password_changed": "Password changed",
  "admin.user_updated": "User updated",
  "admin.user_deleted": "User deleted",
  "task.created": "Task created",
  "task.updated": "Task updated",
  "task.deleted": "Task deleted",
};

export type AuditStatus = "success" | "failure";

export interface AuditLog {
  id: number;
  actorId?: number;
  actorEmail: string;
  action: string;
  targetType?: string;
  targetId?: number;
  status: AuditStatus;
  ip?: string;
  userAgent?: string;
  metadata?: Record<string, unknown>;
  createdAt: string;
}

export interface AuditListMeta {
  page: number;
  perPage: number;
  total: number;
}

export interface AuditListResponse {
  data: AuditLog[];
  meta: AuditListMeta;
}

export interface ListAuditParams {
  page?: number;
  perPage?: number;
  action?: string;
  search?: string;
  from?: string; // RFC3339
  to?: string; // RFC3339
  sort?: "created_at" | "action" | "status";
  order?: "asc" | "desc";
}
