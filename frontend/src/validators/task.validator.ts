import * as z from "zod";
import { TASK_PRIORITIES } from "@/types/task";

const dateField = z
  .string()
  .trim()
  .optional()
  .or(z.literal("").transform(() => undefined))
  .refine((v) => !v || !Number.isNaN(new Date(v).getTime()), "Invalid date");

const taskFieldsBase = {
  title: z.string().trim().min(1, "Title is required"),
  description: z.string().trim().optional(),
  priority: z.enum(TASK_PRIORITIES).optional(),
  startDate: dateField,
  dueDate: dateField,
  assigneeId: z.number().int().positive().optional(),
};

export const newTaskSchema = z.object(taskFieldsBase).superRefine((v, ctx) => {
  if (v.startDate && v.dueDate && new Date(v.dueDate) < new Date(v.startDate)) {
    ctx.addIssue({
      code: "custom",
      path: ["dueDate"],
      message: "Due date must be on or after start date",
    });
  }
});

export type NewTaskInput = z.infer<typeof newTaskSchema>;

export const editTaskSchema = z
  .object({
    ...taskFieldsBase,
    status: z.enum(["todo", "in_progress", "done"]),
  })
  .superRefine((v, ctx) => {
    if (
      v.startDate &&
      v.dueDate &&
      new Date(v.dueDate) < new Date(v.startDate)
    ) {
      ctx.addIssue({
        code: "custom",
        path: ["dueDate"],
        message: "Due date must be on or after start date",
      });
    }
  });

export type EditTaskInput = z.infer<typeof editTaskSchema>;
