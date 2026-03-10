import api from "@/lib/axios";

export interface Task {
	id: number;
	title: string;
	description: string;
	completed: boolean;
	createdAt: string;
}

export interface CreateTaskRequest {
	title: string;
	description?: string;
}

export const taskService = {
	getTasks: async (): Promise<Task[]> => {
		const response = await api.get<Task[]>("/tasks");
		return response.data;
	},

	createTask: async (data: CreateTaskRequest): Promise<Task> => {
		const response = await api.post<Task>("/tasks", data);
		return response.data;
	},

	toggleTask: async (id: number, completed: boolean): Promise<Task> => {
		const response = await api.patch<Task>(`/tasks/${id}`, { completed });
		return response.data;
	},

	deleteTask: async (id: number): Promise<void> => {
		await api.delete(`/tasks/${id}`);
	},
};
