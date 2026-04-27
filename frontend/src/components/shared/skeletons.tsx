import { Skeleton } from "@/components/ui/skeleton";

// Generic table-shaped skeleton — used by both the tasks and users pages
// while their first fetch is in flight. Visually matches the table layout
// (header bar + N rows) so the page doesn't shift when data arrives.
interface TableSkeletonProps {
  rows?: number;
  columns?: number;
}

export function TableSkeleton({ rows = 5, columns = 4 }: TableSkeletonProps) {
  return (
    <div className="space-y-3" aria-label="Loading content">
      <div className="flex gap-2">
        <Skeleton className="h-9 w-64" />
        <div className="flex-1" />
        <Skeleton className="h-9 w-40" />
      </div>
      <div className="overflow-hidden rounded-lg border">
        <div className="flex gap-4 border-b p-3">
          {Array.from({ length: columns }).map((_, i) => (
            <Skeleton key={i} className="h-4 flex-1" />
          ))}
        </div>
        {Array.from({ length: rows }).map((_, i) => (
          <div key={i} className="flex gap-4 border-b p-3 last:border-b-0">
            {Array.from({ length: columns }).map((_, j) => (
              <Skeleton key={j} className="h-4 flex-1" />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

// Stats grid skeleton — used by the dashboard. Matches the 4-card layout
// so the page doesn't reflow when the real KPIs land.
export function StatsGridSkeleton({ count = 4 }: { count?: number }) {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: count }).map((_, i) => (
        <div
          key={i}
          className="border-l-4 border-l-muted rounded-md border p-4"
        >
          <div className="flex items-center gap-4">
            <Skeleton className="h-11 w-11 rounded-xl" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-3 w-20" />
              <Skeleton className="h-6 w-12" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

// Generic centered card-y skeleton for things like the profile page where
// there isn't a table or a grid — just a hero block.
export function CardSkeleton() {
  return (
    <div className="space-y-4 rounded-md border p-6" aria-label="Loading">
      <div className="flex items-center gap-4">
        <Skeleton className="h-20 w-20 rounded-full" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-6 w-48" />
          <Skeleton className="h-4 w-32" />
        </div>
      </div>
      <div className="border-t pt-4">
        <div className="grid grid-cols-2 gap-4">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-4 w-24" />
        </div>
      </div>
    </div>
  );
}
