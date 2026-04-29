import {
  createRootRoute,
  Outlet,
  redirect,
  useLocation,
} from "@tanstack/react-router";
import { AppShell } from "@/components/layout/app-shell";
import { Toaster } from "@/components/ui/sonner";
import {
  isAuthenticated,
  selectIsAuthenticated,
  useAuthStore,
} from "@/store/auth.store";

const PUBLIC_ROUTES = ["/", "/login", "/register"];
const ADMIN_ROUTES = ["/users", "/audit-logs"];

const isPublic = (path: string) => PUBLIC_ROUTES.includes(path);
const isAdminRoute = (path: string) =>
  ADMIN_ROUTES.some((p) => path === p || path.startsWith(p + "/"));

export const Route = createRootRoute({
  beforeLoad: ({ location }) => {
    const authed = isAuthenticated();

    if (!authed && !isPublic(location.pathname)) {
      throw redirect({
        to: "/login",
        search: { redirect: location.href },
      });
    }
    if (authed && ["/", "/login", "/register"].includes(location.pathname)) {
      throw redirect({ to: "/dashboard" });
    }
    if (authed && isAdminRoute(location.pathname)) {
      const role = useAuthStore.getState().user?.role;
      if (role !== "admin") {
        throw redirect({ to: "/dashboard" });
      }
    }
  },
  component: RootComponent,
});

function RootComponent() {
  const location = useLocation();
  const authed = useAuthStore(selectIsAuthenticated);
  const useShell = authed && !isPublic(location.pathname);

  return (
    <div className="bg-background text-foreground min-h-screen font-sans antialiased">
      {useShell ? (
        <AppShell>
          <Outlet />
        </AppShell>
      ) : (
        <Outlet />
      )}
      <Toaster position="top-right" richColors />
    </div>
  );
}
