export type TaskStatus = "todo" | "in_progress" | "done";

// Priority — kept aligned with backend model.TaskPriority constants. Order
// of the array drives select-menu order and visual sort comparators.
export const TASK_PRIORITIES = ["low", "medium", "high", "urgent"] as const;
export type TaskPriority = (typeof TASK_PRIORITIES)[number];

export interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  startDate?: string | null;
  dueDate?: string | null;
  assigneeId?: number | null;
  assigneeEmail?: string;
  userId: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
  priority?: TaskPriority;
  startDate?: string;
  dueDate?: string;
  assigneeId?: number;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  // ISO datetime strings, or use clear*=true to reset to null.
  startDate?: string;
  dueDate?: string;
  assigneeId?: number;
  clearStartDate?: boolean;
  clearDueDate?: boolean;
  clearAssignee?: boolean;
}

export type TaskSortField =
  | "created_at"
  | "updated_at"
  | "title"
  | "status"
  | "priority"
  | "due_date";
export type SortOrder = "asc" | "desc";

export interface ListTasksParams {
  page?: number;
  perPage?: number;
  status?: TaskStatus;
  priority?: TaskPriority;
  assignee?: "me";
  search?: string;
  sort?: TaskSortField;
  order?: SortOrder;
}

export interface TaskListMeta {
  page: number;
  perPage: number;
  total: number;
}

export interface TaskListCounts {
  all: number;
  todo: number;
  in_progress: number;
  done: number;
}

export interface TaskListResponse {
  data: Task[];
  meta: TaskListMeta;
  counts: TaskListCounts;
}

export interface AssignableUser {
  id: number;
  email: string;
  role: string;
}

export interface AssignableUserListResponse {
  data: AssignableUser[];
}
