import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { Trash2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import { taskService } from "@/services/task.service";

export const Route = createLazyFileRoute("/tasks")({
	component: Tasks,
});

function Tasks() {
	const [newTaskTitle, setNewTaskTitle] = useState("");
	const queryClient = useQueryClient();

	const {
		data: tasks,
		isLoading,
		error,
	} = useQuery({
		queryKey: ["tasks"],
		queryFn: taskService.getTasks,
	});

	const createTaskMutation = useMutation({
		mutationFn: taskService.createTask,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tasks"] });
			setNewTaskTitle("");
			toast.success("Task created");
		},
		onError: (err) => {
			const message = isAxiosError(err) ? err.response?.data?.error : undefined;
			toast.error(message || "Failed to create task");
		},
	});

	const toggleTaskMutation = useMutation({
		mutationFn: ({ id, completed }: { id: number; completed: boolean }) =>
			taskService.toggleTask(id, completed),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tasks"] });
		},
	});

	const deleteTaskMutation = useMutation({
		mutationFn: (id: number) => taskService.deleteTask(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tasks"] });
			toast.success("Task deleted");
		},
	});

	const handleAddTask = (e: React.FormEvent) => {
		e.preventDefault();
		if (!newTaskTitle.trim()) return;
		createTaskMutation.mutate({ title: newTaskTitle });
	};

	if (isLoading) return <div className="p-8 text-center">Loading tasks...</div>;
	if (error)
		return (
			<div className="text-destructive p-8 text-center">
				Error loading tasks
			</div>
		);

	return (
		<div className="mx-auto max-w-4xl space-y-6 p-4">
			<Card>
				<CardHeader>
					<CardTitle>Add New Task</CardTitle>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleAddTask} className="flex gap-2">
						<Input
							placeholder="What needs to be done?"
							value={newTaskTitle}
							onChange={(e) => setNewTaskTitle(e.target.value)}
							disabled={createTaskMutation.isPending}
						/>
						<Button type="submit" disabled={createTaskMutation.isPending}>
							{createTaskMutation.isPending ? "Adding..." : "Add"}
						</Button>
					</form>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>My Tasks</CardTitle>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead className="w-12.5">Status</TableHead>
								<TableHead>Task</TableHead>
								<TableHead className="w-25 text-right">Actions</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{tasks?.length === 0 ? (
								<TableRow>
									<TableCell
										colSpan={3}
										className="text-muted-foreground py-8 text-center"
									>
										No tasks found. Add one above!
									</TableCell>
								</TableRow>
							) : (
								tasks?.map((task) => (
									<TableRow key={task.id}>
										<TableCell>
											<Checkbox
												checked={task.completed}
												onCheckedChange={(checked) =>
													toggleTaskMutation.mutate({
														id: task.id,
														completed: !!checked,
													})
												}
											/>
										</TableCell>
										<TableCell
											className={
												task.completed
													? "text-muted-foreground line-through"
													: ""
											}
										>
											{task.title}
										</TableCell>
										<TableCell className="text-right">
											<Button
												variant="ghost"
												size="icon"
												onClick={() => deleteTaskMutation.mutate(task.id)}
												className="text-destructive hover:text-destructive/90"
											>
												<Trash2 className="h-4 w-4" />
											</Button>
										</TableCell>
									</TableRow>
								))
							)}
						</TableBody>
					</Table>
				</CardContent>
			</Card>
		</div>
	);
}
