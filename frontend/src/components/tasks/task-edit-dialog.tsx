import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { TaskStatusSelect } from "@/components/tasks/task-status-select";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useUpdateTask } from "@/hooks/use-tasks";
import { getApiError } from "@/lib/api-error";
import type { Task } from "@/types/task";
import {
  editTaskSchema,
  type EditTaskInput,
} from "@/validators/task.validator";

interface Props {
  task: Task | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function TaskEditDialog({ task, open, onOpenChange }: Props) {
  const updateTask = useUpdateTask();
  const form = useForm<EditTaskInput>({
    resolver: zodResolver(editTaskSchema),
    defaultValues: { title: "", description: "", status: "todo" },
  });

  // Reset form whenever the target task changes (or the dialog re-opens).
  useEffect(() => {
    if (task && open) {
      form.reset({
        title: task.title,
        description: task.description ?? "",
        status: task.status,
      });
    }
  }, [task, open, form]);

  if (!task) return null;

  const onSubmit = (values: EditTaskInput) => {
    updateTask.mutate(
      {
        id: task.id,
        data: {
          title: values.title,
          description: values.description ?? "",
          status: values.status,
        },
      },
      {
        onSuccess: () => {
          toast.success("Task updated");
          onOpenChange(false);
        },
        onError: (err) =>
          toast.error(getApiError(err, "Failed to update task")),
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit task</DialogTitle>
          <DialogDescription>
            Update the title, description, or status.
          </DialogDescription>
        </DialogHeader>

        <form
          id="edit-task-form"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-4 pt-2"
        >
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input id="title" {...form.register("title")} />
            {form.formState.errors.title && (
              <p className="text-destructive text-sm font-medium">
                {form.formState.errors.title.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              rows={4}
              {...form.register("description")}
            />
          </div>

          <div className="space-y-2">
            <Label>Status</Label>
            <TaskStatusSelect
              value={form.watch("status")}
              onChange={(s) =>
                form.setValue("status", s, { shouldDirty: true })
              }
              className="w-full"
            />
          </div>
        </form>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={updateTask.isPending}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            form="edit-task-form"
            disabled={updateTask.isPending}
          >
            {updateTask.isPending ? "Saving..." : "Save changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
