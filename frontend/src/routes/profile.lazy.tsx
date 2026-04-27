import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { PageError } from "@/components/shared/page-state";
import { CardSkeleton } from "@/components/shared/skeletons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
        <CardHeader className="flex flex-row items-center gap-4">
          <Avatar className="h-20 w-20">
            <AvatarImage
              src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user.email}`}
            />
            <AvatarFallback>
              {user.email.substring(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div>
            <CardTitle className="text-3xl font-bold">{user.email}</CardTitle>
            <p className="text-muted-foreground">Role: {role}</p>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 border-t pt-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-muted-foreground text-sm font-medium">
                User ID
              </p>
              <p className="font-mono text-xs">{user.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-sm font-medium">Role</p>
              <p className="capitalize">{role}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
