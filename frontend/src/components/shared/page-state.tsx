// Centered error block used by every route that can fail to load.
// (Loading is handled per-page with skeletons in components/shared/skeletons.tsx
// — they preserve layout and feel snappier than a single "Loading…" line.)
interface ErrorProps {
  label: string;
}

export function PageError({ label }: ErrorProps) {
  return <div className="text-destructive p-8 text-center">{label}</div>;
}
