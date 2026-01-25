import { NavLink } from "@/components/NavLink";
import { SIDEBAR_NAV } from "@/routes/sidebarNav";
import { cn } from "@/lib/utils";
import { ShieldCheck, X } from "lucide-react";
import { UserMenu } from "@/components/common/UserMenu";

interface DashboardSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export function DashboardSidebar({ isOpen, onClose }: DashboardSidebarProps) {
  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          "fixed left-0 top-0 z-50 h-screen w-[280px] border-r border-border bg-elevated text-elevated-foreground transition-transform duration-300 lg:z-30 lg:translate-x-0",
          isOpen ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <div className="flex h-full flex-col">
          {/* Header with close button on mobile */}
          <div className="flex items-center justify-between gap-3 px-6 py-5">
            <div className="flex items-center gap-3">
              <div className="grid h-10 w-10 place-items-center rounded-xl bg-brand text-brand-foreground">
                <ShieldCheck className="h-5 w-5" />
              </div>
              <div className="text-lg font-semibold tracking-tight">
                AdminHub
              </div>
            </div>
            {/* Close button - only visible on mobile */}
            <button
              onClick={onClose}
              className="grid h-9 w-9 place-items-center rounded-xl hover:bg-muted/60 lg:hidden"
              aria-label="Close menu"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          <nav className="px-3">
            <ul className="space-y-1.5">
              {SIDEBAR_NAV.map((item) => {
                const Icon = item.icon;
                return (
                  <li key={item.to}>
                    <NavLink
                      to={item.to}
                      onClick={onClose}
                      className={cn(
                        "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground transition-colors",
                        "hover:bg-muted/60 hover:text-foreground",
                      )}
                      activeClassName="bg-brand text-brand-foreground shadow-sm"
                    >
                      <Icon className="h-4 w-4" />
                      <span className="flex-1">{item.label}</span>
                      {item.badge && (
                        <span
                          className={cn(
                            "rounded-full px-2 py-0.5 text-xs",
                            "bg-brand text-brand-foreground",
                          )}
                        >
                          {item.badge}
                        </span>
                      )}
                    </NavLink>
                  </li>
                );
              })}
            </ul>
          </nav>

          {/* User Menu at bottom */}
          <div className="mt-auto border-t border-border p-4">
            <UserMenu variant="sidebar" />
          </div>
        </div>
      </aside>
    </>
  );
}
