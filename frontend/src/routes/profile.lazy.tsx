import { createLazyFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { authService } from '@/services/auth.service'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'

export const Route = createLazyFileRoute('/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const { data: user, isLoading, error } = useQuery({
    queryKey: ['profile'],
    queryFn: authService.getProfile,
  })

  if (isLoading) return <div className="p-8 text-center">Loading profile...</div>
  if (error) return <div className="p-8 text-center text-destructive">Error loading profile</div>

  return (
    <div className="p-4 max-w-2xl mx-auto">
      <Card>
        <CardHeader className="flex flex-row items-center gap-4">
          <Avatar className="h-20 w-20">
            <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user?.email}`} />
            <AvatarFallback>{user?.email?.substring(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>
          <div>
            <CardTitle className="text-3xl font-bold">{user?.email}</CardTitle>
            <p className="text-muted-foreground">Role: {user?.role || 'user'}</p>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 pt-4 border-t">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground">User ID</p>
              <p className="font-mono text-xs">{user?.id}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-muted-foreground">Role</p>
              <p className="capitalize">{user?.role || 'User'}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
