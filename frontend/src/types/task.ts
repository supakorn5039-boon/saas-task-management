export type TaskStatus = "todo" | "in_progress" | "done";

export interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  createdAt: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
}

export type TaskSortField = "created_at" | "updated_at" | "title" | "status";
export type SortOrder = "asc" | "desc";

export interface ListTasksParams {
  page?: number;
  perPage?: number;
  status?: TaskStatus;
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
