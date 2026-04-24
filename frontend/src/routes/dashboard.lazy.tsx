import { createLazyFileRoute, Link } from "@tanstack/react-router";
import { CheckCircle2, Circle, ListChecks, Loader2 } from "lucide-react";
import { KpiCard } from "@/components/dashboard/kpi-card";
import { TaskStatusBadge } from "@/components/tasks/task-status-badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTasks } from "@/hooks/use-tasks";

export const Route = createLazyFileRoute("/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  // Single fetch powers both KPIs (via counts) and recent activity (via data).
  const { data, isLoading, error } = useTasks({
    page: 1,
    perPage: 5,
    sort: "created_at",
    order: "desc",
  });

  if (isLoading)
    return <div className="p-8 text-center">Loading dashboard...</div>;
  if (error)
    return (
      <div className="text-destructive p-8 text-center">
        Error loading dashboard
      </div>
    );

  const counts = data?.counts ?? { all: 0, todo: 0, in_progress: 0, done: 0 };
  const recent = data?.data ?? [];
  const completionRate =
    counts.all > 0 ? Math.round((counts.done / counts.all) * 100) : 0;

  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <KpiCard
          label="Total tasks"
          value={counts.all}
          icon={ListChecks}
          tone="indigo"
          hint={`${completionRate}% completed`}
        />
        <KpiCard label="To do" value={counts.todo} icon={Circle} tone="slate" />
        <KpiCard
          label="In progress"
          value={counts.in_progress}
          icon={Loader2}
          tone="amber"
        />
        <KpiCard
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
                  <div className="flex items-center gap-3 min-w-0">
                    <TaskStatusBadge status={task.status ?? "todo"} />
                    <span className="truncate text-sm font-medium">
                      {task.title}
                    </span>
                  </div>
                  <span className="text-muted-foreground text-xs whitespace-nowrap">
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

function formatRelative(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const minutes = Math.floor(diff / 60_000);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  if (days < 7) return `${days}d ago`;
  return new Date(iso).toLocaleDateString();
}
