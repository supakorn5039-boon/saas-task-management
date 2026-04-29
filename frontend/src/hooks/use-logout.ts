import { useNavigate } from "@tanstack/react-router";
import { auditService } from "@/services/audit.service";
import { useAuthStore } from "@/store/auth.store";

// Logout sends an audit-log "auth.logout" event server-side, then clears
// local auth state and routes the user back to /login.
//
// Why best-effort? JWT is stateless — there's no server-side session to
// invalidate. The endpoint exists purely to record the event. If it fails
// (network down, token already expired) we still want the user logged out
// locally, so the request error is swallowed.
export function useLogout() {
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);

  return async () => {
    try {
      await auditService.logout();
    } catch {
      // best-effort — see comment above.
    }
    logout();
    navigate({ to: "/login" });
  };
}
