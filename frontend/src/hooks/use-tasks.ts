import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { taskKeys, taskService } from "@/services/task.service";
import type {
  ListTasksParams,
  Task,
  TaskListResponse,
  TaskStatus,
  UpdateTaskRequest,
} from "@/types/task";

// Shared invalidator — every task mutation runs this on success.
// Centralizing it means there's one place to change cache behavior
// (e.g., switch to optimistic updates) for every task mutation at once.
function useInvalidateTasks() {
  const qc = useQueryClient();
  return () => qc.invalidateQueries({ queryKey: taskKeys.all });
}

export function useTasks(params: ListTasksParams) {
  return useQuery({
    queryKey: taskKeys.list(params),
    queryFn: () => taskService.listTasks(params),
    staleTime: 5_000,
    // Keep showing the previous page while a new one loads — prevents the
    // whole route from flashing "Loading..." on every sort/filter/page change.
    placeholderData: keepPreviousData,
  });
}

export function useCreateTask() {
  const invalidate = useInvalidateTasks();
  return useMutation({
    mutationFn: taskService.createTask,
    onSuccess: invalidate,
  });
}

export function useUpdateTask() {
  const invalidate = useInvalidateTasks();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateTaskRequest }) =>
      taskService.updateTask(id, data),
    onSuccess: invalidate,
  });
}

// Status flip is the most-common interaction — give it instant feedback by
// patching every cached list page in place, then roll back if the server rejects.
type ListSnapshots = [readonly unknown[], TaskListResponse | undefined][];

export function useUpdateTaskStatus() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: number; status: TaskStatus }) =>
      taskService.updateStatus(id, status),
    onMutate: async ({ id, status }) => {
      await qc.cancelQueries({ queryKey: taskKeys.lists() });

      const snapshots: ListSnapshots = qc.getQueriesData<TaskListResponse>({
        queryKey: taskKeys.lists(),
      });

      for (const [key, value] of snapshots) {
        if (!value) continue;
        qc.setQueryData<TaskListResponse>(key, {
          ...value,
          data: value.data.map((t: Task) =>
            t.id === id ? { ...t, status } : t,
          ),
        });
      }
      return { snapshots };
    },
    onError: (_err, _vars, context) => {
      // Restore every page we patched.
      context?.snapshots.forEach(([key, value]) => qc.setQueryData(key, value));
    },
    onSettled: () => qc.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

export function useDeleteTask() {
  const invalidate = useInvalidateTasks();
  return useMutation({
    mutationFn: (id: number) => taskService.deleteTask(id),
    onSuccess: invalidate,
  });
}
