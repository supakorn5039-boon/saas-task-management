import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { taskKeys, taskService } from "@/services/task.service";
import type { ListTasksParams, TaskStatus } from "@/types/task";

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

export function useUpdateTaskStatus() {
  const invalidate = useInvalidateTasks();
  return useMutation({
    mutationFn: ({ id, status }: { id: number; status: TaskStatus }) =>
      taskService.updateStatus(id, status),
    onSuccess: invalidate,
  });
}

export function useDeleteTask() {
  const invalidate = useInvalidateTasks();
  return useMutation({
    mutationFn: (id: number) => taskService.deleteTask(id),
    onSuccess: invalidate,
  });
}
