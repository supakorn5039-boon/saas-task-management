import * as z from "zod";

export const editUserSchema = z.object({
  role: z.enum(["admin", "manager", "user"]),
  status: z.union([z.literal(0), z.literal(1)]),
});

export type EditUserInput = z.infer<typeof editUserSchema>;
