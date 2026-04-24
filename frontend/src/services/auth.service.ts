import api from "@/lib/axios";
import { useAuthStore } from "@/store/auth.store";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  UserProfile,
} from "@/types/auth";

export const authService = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>("/auth/login", data);
    if (response.data.token) {
      useAuthStore.getState().setAuth(response.data.token, response.data.user);
    }
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>("/auth/register", data);
    if (response.data.token) {
      useAuthStore.getState().setAuth(response.data.token, response.data.user);
    }
    return response.data;
  },

  getProfile: async (): Promise<UserProfile> => {
    const response = await api.get<UserProfile>("/user/profile");
    const token = useAuthStore.getState().token;
    if (token) {
      useAuthStore.setState({ user: response.data });
    }
    return response.data;
  },

  logout: () => {
    useAuthStore.getState().logout();
  },

  isAuthenticated: () => {
    return useAuthStore.getState().isAuthenticated();
  },
};
