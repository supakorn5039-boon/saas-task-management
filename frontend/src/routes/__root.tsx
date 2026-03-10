import { createRootRoute, Link, Outlet, useNavigate, redirect } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { Toaster } from '@/components/ui/sonner'
import { useAuthStore } from '@/store/auth.store'
import { Button } from '@/components/ui/button'
import { LogOut, User, CheckSquare, LayoutDashboard } from 'lucide-react'

export const Route = createRootRoute({
  beforeLoad: ({ location }) => {
    const isAuthenticated = useAuthStore.getState().isAuthenticated()
    const isPublicRoute = ['/', '/login', '/register'].includes(location.pathname)

    // Redirect to login if not authenticated and trying to access a protected route
    if (!isAuthenticated && !isPublicRoute) {
      throw redirect({
        to: '/login',
        search: {
          redirect: location.href,
        },
      })
    }

    // Redirect to tasks if authenticated and trying to access login/register
    if (isAuthenticated && ['/login', '/register'].includes(location.pathname)) {
      throw redirect({
        to: '/tasks',
      })
    }
  },
  component: RootComponent,
})

function RootComponent() {
  const navigate = useNavigate()
  const { isAuthenticated, logout, user } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate({ to: '/login' })
  }

  return (
    <div className="min-h-screen bg-background font-sans antialiased text-foreground">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 flex h-16 items-center justify-between">
          <div className="flex items-center gap-8">
            <Link to="/" className="flex items-center space-x-2 transition-opacity hover:opacity-80">
              <CheckSquare className="h-6 w-6 text-primary" />
              <span className="font-bold text-lg tracking-tight">SaaS Task</span>
            </Link>
            <nav className="hidden md:flex items-center space-x-6 text-sm font-medium">
              {isAuthenticated() ? (
                <>
                  <Link to="/tasks" className="transition-colors hover:text-primary text-muted-foreground [&.active]:text-foreground [&.active]:font-semibold">
                    <div className="flex items-center gap-2">
                      <LayoutDashboard className="h-4 w-4" />
                      Tasks
                    </div>
                  </Link>
                  <Link to="/profile" className="transition-colors hover:text-primary text-muted-foreground [&.active]:text-foreground [&.active]:font-semibold">
                    <div className="flex items-center gap-2">
                      <User className="h-4 w-4" />
                      Profile
                    </div>
                  </Link>
                </>
              ) : (
                <Link to="/" className="transition-colors hover:text-primary text-muted-foreground [&.active]:text-foreground">
                  Home
                </Link>
              )}
            </nav>
          </div>
          <div className="flex items-center gap-4">
            {isAuthenticated() ? (
              <div className="flex items-center gap-4">
                <span className="text-sm text-muted-foreground hidden sm:inline">
                  Welcome, <span className="text-foreground font-medium">{user?.email}</span>
                </span>
                <Button variant="ghost" size="sm" onClick={handleLogout} className="gap-2 font-medium">
                  <LogOut className="h-4 w-4" />
                  <span className="hidden sm:inline">Logout</span>
                </Button>
              </div>
            ) : (
              <div className="flex gap-2">
                <Link to="/login">
                  <Button variant="ghost" size="sm">Login</Button>
                </Link>
                <Link to="/register">
                  <Button size="sm">Sign Up</Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </header>
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>
      <Toaster position="top-right" richColors />
      <TanStackRouterDevtools />
    </div>
  )
}
