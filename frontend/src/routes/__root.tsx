import {
  createRootRoute,
  Outlet,
  redirect,
  useLocation,
} from "@tanstack/react-router";
import { AppShell } from "@/components/layout/app-shell";
import { Toaster } from "@/components/ui/sonner";
import { useAuthStore } from "@/store/auth.store";

const PUBLIC_ROUTES = ["/", "/login", "/register"];

export const Route = createRootRoute({
  beforeLoad: ({ location }) => {
    const isAuthenticated = useAuthStore.getState().isAuthenticated();
    const isPublic = PUBLIC_ROUTES.includes(location.pathname);

    if (!isAuthenticated && !isPublic) {
      throw redirect({
        to: "/login",
        search: { redirect: location.href },
      });
    }

    if (
      isAuthenticated &&
      ["/login", "/register"].includes(location.pathname)
    ) {
      throw redirect({ to: "/dashboard" });
    }
  },
  component: RootComponent,
});

function RootComponent() {
  const location = useLocation();
  const { isAuthenticated } = useAuthStore();
  const useShell =
    isAuthenticated() && !PUBLIC_ROUTES.includes(location.pathname);

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
