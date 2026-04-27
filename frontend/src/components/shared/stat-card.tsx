import type { LucideIcon } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

export type StatTone = "indigo" | "slate" | "amber" | "emerald";

interface Props {
  label: string;
  value: number | string;
  icon: LucideIcon;
  tone?: StatTone;
  hint?: string;
}

// Tinted chip + accent border per tone — keeps the look consistent across
// any page that renders stats.
const TONE_STYLES: Record<StatTone, { chip: string; border: string }> = {
  indigo: {
    chip: "bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300",
    border: "border-l-indigo-500",
  },
  slate: {
    chip: "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200",
    border: "border-l-slate-400",
  },
  amber: {
    chip: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
    border: "border-l-amber-500",
  },
  emerald: {
    chip: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300",
    border: "border-l-emerald-500",
  },
};

export function StatCard({
  label,
  value,
  icon: Icon,
  tone = "indigo",
  hint,
}: Props) {
  const t = TONE_STYLES[tone];
  return (
    <Card
      className={`overflow-hidden border-l-4 transition-shadow hover:shadow-md ${t.border}`}
    >
      <CardContent className="flex items-center gap-4 p-4">
        <div
          className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl ${t.chip}`}
        >
          <Icon className="h-5 w-5" />
        </div>
        <div className="min-w-0 flex-1">
          <div className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            {label}
          </div>
          <div className="mt-0.5 text-2xl font-bold leading-tight">{value}</div>
          {hint && (
            <p className="text-muted-foreground mt-0.5 text-xs">{hint}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
