import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { authService } from "@/services/auth.service";

export const Route = createLazyFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const {
    data: user,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["profile"],
    queryFn: authService.getProfile,
  });

  if (isLoading)
    return <div className="p-8 text-center">Loading profile...</div>;
  if (error)
    return (
      <div className="text-destructive p-8 text-center">
        Error loading profile
      </div>
    );

  return (
    <div className="mx-auto max-w-2xl p-4">
      <Card>
        <CardHeader className="flex flex-row items-center gap-4">
          <Avatar className="h-20 w-20">
            <AvatarImage
              src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user?.email}`}
            />
            <AvatarFallback>
              {user?.email?.substring(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div>
            <CardTitle className="text-3xl font-bold">{user?.email}</CardTitle>
            <p className="text-muted-foreground">
              Role: {user?.role || "user"}
            </p>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 border-t pt-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-muted-foreground text-sm font-medium">
                User ID
              </p>
              <p className="font-mono text-xs">{user?.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-sm font-medium">Role</p>
              <p className="capitalize">{user?.role || "User"}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
