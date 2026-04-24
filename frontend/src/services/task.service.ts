import api from "@/lib/axios";
import type {
  CreateTaskRequest,
  ListTasksParams,
  Task,
  TaskListResponse,
  TaskStatus,
} from "@/types/task";

// Query key factory — single source of truth for React Query keys.
// Hierarchy lets us invalidate everything with `taskKeys.all`.
export const taskKeys = {
  all: ["tasks"] as const,
  lists: () => [...taskKeys.all, "list"] as const,
  list: (params: ListTasksParams) => [...taskKeys.lists(), params] as const,
  details: () => [...taskKeys.all, "detail"] as const,
  detail: (id: number) => [...taskKeys.details(), id] as const,
};

export const taskService = {
  listTasks: async (params: ListTasksParams): Promise<TaskListResponse> => {
    const response = await api.get<TaskListResponse>("/tasks", {
      params: {
        page: params.page,
        per_page: params.perPage,
        status: params.status,
        search: params.search || undefined,
        sort: params.sort,
        order: params.order,
      },
    });
    return response.data;
  },

  createTask: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await api.post<Task>("/tasks", data);
    return response.data;
  },

  updateStatus: async (id: number, status: TaskStatus): Promise<Task> => {
    const response = await api.put<Task>(`/tasks/${id}`, { status });
    return response.data;
  },

  deleteTask: async (id: number): Promise<void> => {
    await api.delete(`/tasks/${id}`);
  },
};
