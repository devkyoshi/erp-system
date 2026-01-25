import {
  BarChart3,
  Boxes,
  CalendarDays,
  CreditCard,
  FileText,
  LayoutDashboard,
  Settings,
  Users,
  Users2,
  ShoppingBag,
} from "lucide-react";

export type SidebarNavItem = {
  label: string;
  to: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string;
};

export const SIDEBAR_NAV: SidebarNavItem[] = [
  { label: "Dashboard", to: "/app/dashboard", icon: LayoutDashboard },
  { label: "Reports", to: "/app/reports", icon: FileText },
  { label: "Orders", to: "/app/orders", icon: ShoppingBag, badge: "12" },
  { label: "Inventory", to: "/app/inventory", icon: Boxes, badge: "3" },
  { label: "Users", to: "/app/users", icon: Users },
  { label: "Settings", to: "/app/settings", icon: Settings },
];
