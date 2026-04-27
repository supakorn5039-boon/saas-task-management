import * as z from "zod";

export const newTaskSchema = z.object({
  title: z.string().trim().min(1, "Title is required"),
  description: z.string().trim().optional(),
});

export type NewTaskInput = z.infer<typeof newTaskSchema>;

export const editTaskSchema = z.object({
  title: z.string().trim().min(1, "Title is required"),
  description: z.string().trim().optional(),
  status: z.enum(["todo", "in_progress", "done"]),
});

export type EditTaskInput = z.infer<typeof editTaskSchema>;
