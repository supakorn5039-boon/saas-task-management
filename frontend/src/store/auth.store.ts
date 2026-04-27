import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { UserProfile } from "@/types/auth";

interface AuthState {
  token: string | null;
  user: UserProfile | null;
  setAuth: (token: string, user: UserProfile) => void;
  setUser: (user: UserProfile) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      setAuth: (token, user) => set({ token, user }),
      setUser: (user) => set({ user }),
      logout: () => set({ token: null, user: null }),
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({ token: state.token, user: state.user }),
    },
  ),
);

// Convenience selectors — encourage components to subscribe to the smallest
// slice they need (avoids re-renders on unrelated state changes).
export const selectIsAuthenticated = (s: AuthState) => !!s.token;
export const selectUser = (s: AuthState) => s.user;

// Read auth status from anywhere (router beforeLoad, axios interceptor) without
// a hook. Mirrors `useAuthStore.getState().token` but reads through the selector.
export function isAuthenticated(): boolean {
  return selectIsAuthenticated(useAuthStore.getState());
}
