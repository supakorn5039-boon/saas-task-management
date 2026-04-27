import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import {
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from "@/components/ui/dropdown-menu";

// Three options instead of a binary toggle: explicit light/dark + "system"
// (follows OS preference). Designed to drop into an existing DropdownMenu.
//
// Why no Label / Separator inside the submenu: base-ui's Menu requires
// Menu.Label and Menu.Separator to live inside Menu.Group. Wrapping them just
// to add decoration isn't worth it — the radio group's three items are clear
// enough on their own.
export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  // next-themes can't know the resolved theme until after hydration; render a
  // static "Theme" trigger label until then to avoid a hydration mismatch.
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  const current = mounted ? (theme ?? "system") : "system";

  return (
    <DropdownMenuSub>
      <DropdownMenuSubTrigger>
        <ThemeIcon theme={current} />
        <span className="ml-2">Theme</span>
      </DropdownMenuSubTrigger>
      <DropdownMenuSubContent>
        <DropdownMenuRadioGroup value={current} onValueChange={setTheme}>
          <DropdownMenuRadioItem value="light">
            <Sun className="mr-2 h-4 w-4" /> Light
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="dark">
            <Moon className="mr-2 h-4 w-4" /> Dark
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="system">
            <Monitor className="mr-2 h-4 w-4" /> System
          </DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuSubContent>
    </DropdownMenuSub>
  );
}

// Tiny standalone button variant — handy for log-in/register pages where
// the dropdown menu isn't available yet.
export function ThemeQuickToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  const next = theme === "dark" ? "light" : "dark";
  return (
    <button
      type="button"
      onClick={() => setTheme(next)}
      aria-label={`Switch to ${next} mode`}
      className="hover:bg-accent inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors"
    >
      {!mounted ? null : theme === "dark" ? (
        <Sun className="h-4 w-4" />
      ) : (
        <Moon className="h-4 w-4" />
      )}
    </button>
  );
}

function ThemeIcon({ theme }: { theme: string }) {
  if (theme === "dark") return <Moon className="h-4 w-4" />;
  if (theme === "light") return <Sun className="h-4 w-4" />;
  return <Monitor className="h-4 w-4" />;
}
