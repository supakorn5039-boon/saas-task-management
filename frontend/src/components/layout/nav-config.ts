import {
  LayoutDashboard,
  ListChecks,
  Settings,
  User,
  Users,
  type LucideIcon,
} from "lucide-react";

export type Role = "admin" | "manager" | "user";

export interface NavItem {
  label: string;
  to: string;
  icon: LucideIcon;
  // If undefined, visible to all authenticated users.
  // Otherwise, only roles in this list see the item.
  roles?: Role[];
}

export interface NavSection {
  label?: string;
  items: NavItem[];
}

export const NAV_SECTIONS: NavSection[] = [
  {
    label: "Workspace",
    items: [
      { label: "Dashboard", to: "/dashboard", icon: LayoutDashboard },
      { label: "Tasks", to: "/tasks", icon: ListChecks },
    ],
  },
  {
    label: "Admin",
    items: [
      { label: "Users", to: "/users", icon: Users, roles: ["admin"] },
      { label: "Settings", to: "/settings", icon: Settings, roles: ["admin"] },
    ],
  },
  {
    label: "Account",
    items: [{ label: "Profile", to: "/profile", icon: User }],
  },
];

export function visibleSections(role: Role | undefined): NavSection[] {
  return NAV_SECTIONS.map((section) => ({
    ...section,
    items: section.items.filter(
      (item) => !item.roles || (role && item.roles.includes(role)),
    ),
  })).filter((section) => section.items.length > 0);
}
