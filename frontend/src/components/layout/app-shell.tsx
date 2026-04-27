import { useLocation } from "@tanstack/react-router";
import { AppSidebar } from "@/components/layout/app-sidebar";
import { UserMenu } from "@/components/layout/user-menu";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { TooltipProvider } from "@/components/ui/tooltip";
import { selectUser, useAuthStore } from "@/store/auth.store";
import type { Role } from "./nav-config";

interface Props {
  children: React.ReactNode;
}

export function AppShell({ children }: Props) {
  const user = useAuthStore(selectUser);
  const location = useLocation();

  return (
    <TooltipProvider>
      <SidebarProvider>
        <AppSidebar
          role={user?.role as Role | undefined}
          currentPath={location.pathname}
        />
        <SidebarInset>
          <header className="bg-background/95 supports-backdrop-filter:bg-background/60 sticky top-0 z-10 flex h-14 items-center gap-2 border-b px-4 backdrop-blur">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <div className="flex-1" />
            <UserMenu />
          </header>
          <main className="flex-1 p-4 md:p-6">{children}</main>
        </SidebarInset>
      </SidebarProvider>
    </TooltipProvider>
  );
}
