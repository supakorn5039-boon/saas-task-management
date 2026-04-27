import axios from "axios";
import { navigateTo } from "@/lib/router";
import { useAuthStore } from "@/store/auth.store";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// Add a request interceptor to include the auth token
api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Auto-logout on 401. Skip the redirect when the failing call is the login
// endpoint itself — we want the form to surface "invalid credentials" inline,
// not bounce the user back to /login mid-submit.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status;
    const url: string | undefined = error?.config?.url;
    const isAuthCall = typeof url === "string" && url.includes("/auth/");

    if (status === 401 && !isAuthCall) {
      useAuthStore.getState().logout();
      navigateTo("/login");
    }
    return Promise.reject(error);
  },
);

export default api;
