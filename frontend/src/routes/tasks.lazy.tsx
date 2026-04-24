import { createLazyFileRoute } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { useState } from "react";
import { toast } from "sonner";
import { TaskTable } from "@/components/tasks/task-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  useCreateTask,
  useDeleteTask,
  useTasks,
  useUpdateTaskStatus,
} from "@/hooks/use-tasks";
import type { SortOrder, TaskSortField, TaskStatus } from "@/types/task";

export const Route = createLazyFileRoute("/tasks")({
  component: Tasks,
});

const DEFAULT_PER_PAGE = 10;

function Tasks() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  // Server-driven list state
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(DEFAULT_PER_PAGE);
  const [statusFilter, setStatusFilter] = useState<TaskStatus | undefined>();
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<{ field: TaskSortField; order: SortOrder }>({
    field: "created_at",
    order: "desc",
  });

  const { data, isLoading, isFetching, error } = useTasks({
    page,
    perPage,
    status: statusFilter,
    search,
    sort: sort.field,
    order: sort.order,
  });

  const createTask = useCreateTask();
  const updateStatus = useUpdateTaskStatus();
  const deleteTask = useDeleteTask();

  const handleAddTask = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    createTask.mutate(
      { title, description: description.trim() || undefined },
      {
        onSuccess: () => {
          setTitle("");
          setDescription("");
          toast.success("Task created");
        },
        onError: (err) => {
          const msg = isAxiosError(err) ? err.response?.data?.error : undefined;
          toast.error(msg || "Failed to create task");
        },
      },
    );
  };

  // Reset to page 1 whenever filters change so the user doesn't end up on an empty page
  const resetAndRun =
    <T,>(setter: (v: T) => void) =>
    (value: T) => {
      setPage(1);
      setter(value);
    };

  if (isLoading) return <div className="p-8 text-center">Loading tasks...</div>;
  if (error)
    return (
      <div className="text-destructive p-8 text-center">
        Error loading tasks
      </div>
    );

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-4">
      <Card>
        <CardHeader>
          <CardTitle>Add New Task</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAddTask} className="space-y-3">
            <Input
              placeholder="What needs to be done?"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={createTask.isPending}
            />
            <Textarea
              placeholder="Description (optional)"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={createTask.isPending}
              rows={3}
            />
            <div className="flex justify-end">
              <Button type="submit" disabled={createTask.isPending}>
                {createTask.isPending ? "Adding..." : "Add"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            My Tasks
            {isFetching && !isLoading && (
              <span className="text-muted-foreground text-xs font-normal">
                Updating…
              </span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <TaskTable
            tasks={data?.data ?? []}
            meta={data?.meta ?? { page, perPage, total: 0 }}
            counts={
              data?.counts ?? { all: 0, todo: 0, in_progress: 0, done: 0 }
            }
            statusFilter={statusFilter}
            search={search}
            sort={sort}
            onStatusFilterChange={resetAndRun(setStatusFilter)}
            onSearchChange={resetAndRun(setSearch)}
            onSortChange={resetAndRun(setSort)}
            onPageChange={setPage}
            onPerPageChange={resetAndRun(setPerPage)}
            onTaskStatusChange={(id, status) =>
              updateStatus.mutate(
                { id, status },
                { onError: () => toast.error("Failed to update status") },
              )
            }
            onTaskDelete={(id) =>
              deleteTask.mutate(id, {
                onSuccess: () => toast.success("Task deleted"),
                onError: () => toast.error("Failed to delete task"),
              })
            }
          />
        </CardContent>
      </Card>
    </div>
  );
}
