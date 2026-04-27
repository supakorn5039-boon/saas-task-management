import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { TaskStatusBadge } from "@/components/tasks/task-status-badge";

describe("TaskStatusBadge", () => {
  it("renders the human-readable label for each status", () => {
    const { rerender } = render(<TaskStatusBadge status="todo" />);
    expect(screen.getByText("To do")).toBeInTheDocument();

    rerender(<TaskStatusBadge status="in_progress" />);
    expect(screen.getByText("In progress")).toBeInTheDocument();

    rerender(<TaskStatusBadge status="done" />);
    expect(screen.getByText("Done")).toBeInTheDocument();
  });
});
