import { createLazyFileRoute, Link } from "@tanstack/react-router";
import { CheckCircle2, Circle, ListChecks, Loader2 } from "lucide-react";
import { PageError } from "@/components/shared/page-state";
import { StatsGridSkeleton } from "@/components/shared/skeletons";
import { Skeleton } from "@/components/ui/skeleton";
import { StatCard } from "@/components/shared/stat-card";
import { TaskStatusBadge } from "@/components/tasks/task-status-badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTasks } from "@/hooks/use-tasks";
import { formatRelative } from "@/lib/date";

export const Route = createLazyFileRoute("/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  // Single fetch powers both the stats (via counts) and recent activity (via data).
  const { data, isLoading, error } = useTasks({
    page: 1,
    perPage: 5,
    sort: "created_at",
    order: "desc",
  });

  if (error) return <PageError label="Error loading dashboard" />;

  if (isLoading) {
    return (
      <div className="mx-auto max-w-6xl space-y-6">
        <StatsGridSkeleton />
        <div className="rounded-md border p-6 space-y-3">
          <Skeleton className="h-6 w-40" />
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-8 w-full" />
          ))}
        </div>
      </div>
    );
  }

  const counts = data?.counts ?? { all: 0, todo: 0, in_progress: 0, done: 0 };
  const recent = data?.data ?? [];
  const completionRate =
    counts.all > 0 ? Math.round((counts.done / counts.all) * 100) : 0;

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="Total tasks"
          value={counts.all}
          icon={ListChecks}
          tone="indigo"
          hint={`${completionRate}% completed`}
        />
        <StatCard
          label="To do"
          value={counts.todo}
          icon={Circle}
          tone="slate"
        />
        <StatCard
          label="In progress"
          value={counts.in_progress}
          icon={Loader2}
          tone="amber"
        />
        <StatCard
          label="Done"
          value={counts.done}
          icon={CheckCircle2}
          tone="emerald"
        />
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Recent activity</CardTitle>
          <Button
            render={<Link to="/tasks" />}
            nativeButton={false}
            variant="ghost"
            size="sm"
          >
            View all
          </Button>
        </CardHeader>
        <CardContent>
          {recent.length === 0 ? (
            <div className="text-muted-foreground py-8 text-center text-sm">
              No tasks yet. Head over to{" "}
              <Link to="/tasks" className="underline">
                Tasks
              </Link>{" "}
              to add one.
            </div>
          ) : (
            <ul className="space-y-2">
              {recent.map((task) => (
                <li
                  key={task.id}
                  className="hover:bg-accent/30 flex items-center justify-between rounded-md p-2 transition-colors"
                >
                  <div className="flex min-w-0 items-center gap-3">
                    <TaskStatusBadge status={task.status ?? "todo"} />
                    <span className="truncate text-sm font-medium">
                      {task.title}
                    </span>
                  </div>
                  <span className="text-muted-foreground whitespace-nowrap text-xs">
                    {formatRelative(task.createdAt)}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
