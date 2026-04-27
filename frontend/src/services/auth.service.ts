import api from "@/lib/axios";
import { useAuthStore } from "@/store/auth.store";
import type { AuthResponse, Credentials, UserProfile } from "@/types/auth";

// Mirror taskKeys — single source of truth for any auth/user-related query.
export const authKeys = {
  all: ["auth"] as const,
  profile: () => [...authKeys.all, "profile"] as const,
};

export const authService = {
  login: async (data: Credentials): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>("/auth/login", data);
    if (response.data.token) {
      useAuthStore.getState().setAuth(response.data.token, response.data.user);
    }
    return response.data;
  },

  register: async (data: Credentials): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>("/auth/register", data);
    if (response.data.token) {
      useAuthStore.getState().setAuth(response.data.token, response.data.user);
    }
    return response.data;
  },

  getProfile: async (): Promise<UserProfile> => {
    const response = await api.get<UserProfile>("/user/profile");
    // Refresh the cached user (role may have changed server-side); token stays
    // as-is. Going through the public setter avoids touching internal store state.
    useAuthStore.getState().setUser(response.data);
    return response.data;
  },

  changePassword: async (data: {
    currentPassword: string;
    newPassword: string;
  }): Promise<void> => {
    await api.put("/user/password", data);
  },
};
