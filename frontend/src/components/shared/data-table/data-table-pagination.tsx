import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const DEFAULT_PER_PAGE_OPTIONS = [10, 20, 50, 100];

interface Props {
  page: number;
  perPage: number;
  total: number;
  perPageOptions?: number[];
  onPageChange: (page: number) => void;
  onPerPageChange: (perPage: number) => void;
}

export function DataTablePagination({
  page,
  perPage,
  total,
  perPageOptions = DEFAULT_PER_PAGE_OPTIONS,
  onPageChange,
  onPerPageChange,
}: Props) {
  const totalPages = Math.max(1, Math.ceil(total / perPage));
  const start = (page - 1) * perPage + 1;
  const end = Math.min(page * perPage, total);
  const isFirst = page <= 1;
  const isLast = page >= totalPages;

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="text-muted-foreground text-sm">
        {total === 0 ? "No results" : `Showing ${start}–${end} of ${total}`}
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2 text-sm">
          <span className="text-muted-foreground">Rows per page</span>
          <Select
            value={String(perPage)}
            onValueChange={(v) => onPerPageChange(Number(v))}
          >
            <SelectTrigger className="h-8 w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {perPageOptions.map((n) => (
                <SelectItem key={n} value={String(n)}>
                  {n}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="text-muted-foreground text-sm">
          Page {page} of {totalPages}
        </div>

        <div className="flex items-center gap-1">
          <PageButton
            onClick={() => onPageChange(1)}
            disabled={isFirst}
            aria-label="First page"
          >
            <ChevronsLeft className="h-4 w-4" />
          </PageButton>
          <PageButton
            onClick={() => onPageChange(page - 1)}
            disabled={isFirst}
            aria-label="Previous page"
          >
            <ChevronLeft className="h-4 w-4" />
          </PageButton>
          <PageButton
            onClick={() => onPageChange(page + 1)}
            disabled={isLast}
            aria-label="Next page"
          >
            <ChevronRight className="h-4 w-4" />
          </PageButton>
          <PageButton
            onClick={() => onPageChange(totalPages)}
            disabled={isLast}
            aria-label="Last page"
          >
            <ChevronsRight className="h-4 w-4" />
          </PageButton>
        </div>
      </div>
    </div>
  );
}

function PageButton(props: React.ComponentProps<typeof Button>) {
  return (
    <Button variant="outline" size="icon" className="h-8 w-8" {...props} />
  );
}
