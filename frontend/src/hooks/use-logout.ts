import { useNavigate } from "@tanstack/react-router";
import { useAuthStore } from "@/store/auth.store";

export function useLogout() {
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);

  return () => {
    logout();
    navigate({ to: "/login" });
  };
}
