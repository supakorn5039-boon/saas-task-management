import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { PageError } from "@/components/shared/page-state";
import { CardSkeleton } from "@/components/shared/skeletons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatRelative } from "@/lib/date";
import { authKeys, authService } from "@/services/auth.service";
import { selectUser, useAuthStore } from "@/store/auth.store";

export const Route = createLazyFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  // Show the cached user immediately; refetch in the background to pick up
  // server-side role changes.
  const cached = useAuthStore(selectUser);
  const { data, isLoading, error } = useQuery({
    queryKey: authKeys.profile(),
    queryFn: authService.getProfile,
    initialData: cached ?? undefined,
  });

  const user = data ?? cached;
  if (!user && isLoading)
    return (
      <div className="mx-auto max-w-2xl p-4">
        <CardSkeleton />
      </div>
    );
  if (!user && error) return <PageError label="Error loading profile" />;
  if (!user) return null;

  const role = user.role || "user";

  return (
    <div className="mx-auto max-w-2xl p-4">
      <Card>
        <CardHeader className="flex flex-col items-center gap-4 text-center sm:flex-row sm:items-center sm:text-left">
          <Avatar className="h-20 w-20 shrink-0">
            <AvatarImage
              src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user.email}`}
            />
            <AvatarFallback>
              {user.email.substring(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div className="min-w-0 flex-1">
            <CardTitle className="truncate text-xl font-bold sm:text-2xl md:text-3xl">
              {user.email}
            </CardTitle>
            <p className="text-muted-foreground capitalize">Role: {role}</p>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 border-t pt-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <p className="text-muted-foreground text-sm font-medium">Role</p>
              <p className="capitalize">{role}</p>
            </div>
            {user.createdAt && (
              <div>
                <p className="text-muted-foreground text-sm font-medium">
                  Member since
                </p>
                <p>{formatRelative(user.createdAt)}</p>
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
