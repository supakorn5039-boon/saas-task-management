import { zodResolver } from "@hookform/resolvers/zod";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";
import { PageError } from "@/components/shared/page-state";
import { TableSkeleton } from "@/components/shared/skeletons";
import { TaskEditDialog } from "@/components/tasks/task-edit-dialog";
import { TaskTable } from "@/components/tasks/task-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useTaskListState } from "@/hooks/use-task-list-state";
import {
  useCreateTask,
  useDeleteTask,
  useTasks,
  useUpdateTaskStatus,
} from "@/hooks/use-tasks";
import { getApiError } from "@/lib/api-error";
import type { Task } from "@/types/task";
import { newTaskSchema, type NewTaskInput } from "@/validators/task.validator";

export const Route = createLazyFileRoute("/tasks")({
  component: Tasks,
});

function Tasks() {
  const list = useTaskListState();

  const { data, isLoading, isFetching, error } = useTasks({
    page: list.page,
    perPage: list.perPage,
    status: list.status,
    search: list.search,
    sort: list.sort.field,
    order: list.sort.order,
  });

  const createTask = useCreateTask();
  const updateStatus = useUpdateTaskStatus();
  const deleteTask = useDeleteTask();

  const [editing, setEditing] = useState<Task | null>(null);
  const [deleting, setDeleting] = useState<Task | null>(null);

  const form = useForm<NewTaskInput>({
    resolver: zodResolver(newTaskSchema),
    defaultValues: { title: "", description: "" },
  });

  const onAddTask = (values: NewTaskInput) => {
    createTask.mutate(
      { title: values.title, description: values.description || undefined },
      {
        onSuccess: () => {
          form.reset();
          toast.success("Task created");
        },
        onError: (err) =>
          toast.error(getApiError(err, "Failed to create task")),
      },
    );
  };

  const onConfirmDelete = async () => {
    if (!deleting) return;
    await new Promise<void>((resolve) => {
      deleteTask.mutate(deleting.id, {
        onSuccess: () => {
          toast.success("Task deleted");
          resolve();
        },
        onError: (err) => {
          toast.error(getApiError(err, "Failed to delete task"));
          resolve();
        },
      });
    });
  };

  if (error) return <PageError label="Error loading tasks" />;

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-4">
      <Card>
        <CardHeader>
          <CardTitle>Add New Task</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={form.handleSubmit(onAddTask)} className="space-y-3">
            <div className="space-y-1">
              <Input
                placeholder="What needs to be done?"
                disabled={createTask.isPending}
                {...form.register("title")}
              />
              {form.formState.errors.title && (
                <p className="text-destructive text-sm font-medium">
                  {form.formState.errors.title.message}
                </p>
              )}
            </div>
            <Textarea
              placeholder="Description (optional)"
              disabled={createTask.isPending}
              rows={3}
              {...form.register("description")}
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
          {isLoading ? (
            <TableSkeleton rows={5} columns={4} />
          ) : (
            <TaskTable
              tasks={data?.data ?? []}
              meta={
                data?.meta ?? {
                  page: list.page,
                  perPage: list.perPage,
                  total: 0,
                }
              }
              counts={
                data?.counts ?? { all: 0, todo: 0, in_progress: 0, done: 0 }
              }
              state={list}
              actions={list}
              onTaskStatusChange={(id, status) =>
                updateStatus.mutate(
                  { id, status },
                  {
                    onError: (err) =>
                      toast.error(getApiError(err, "Failed to update status")),
                  },
                )
              }
              onTaskEdit={setEditing}
              onTaskDelete={setDeleting}
            />
          )}
        </CardContent>
      </Card>

      <TaskEditDialog
        task={editing}
        open={editing !== null}
        onOpenChange={(open) => !open && setEditing(null)}
      />

      <ConfirmDialog
        open={deleting !== null}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Delete this task?"
        description={
          deleting
            ? `"${deleting.title}" will be removed permanently. This can't be undone.`
            : undefined
        }
        confirmLabel="Delete"
        destructive
        onConfirm={onConfirmDelete}
      />
    </div>
  );
}
