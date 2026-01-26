import { Navigate, Route, Routes } from "react-router-dom";
import NotFound from "@/pages/NotFound";
import { AuthLayout } from "@/components/layout/AuthLayout";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

import { WelcomePage } from "@/pages/auth/Welcome";
import { LoginPage } from "@/pages/auth/Login";
import { RegisterPage } from "@/pages/auth/Register";

import { DashboardPage } from "@/pages/dashboard/Dashboard";
import { CalendarPage } from "@/pages/dashboard/Calendar";
import { InventoryPage } from "@/pages/dashboard/Inventory";
import { OrdersPage } from "@/pages/dashboard/Orders";
import { ReportsPage } from "@/pages/dashboard/Reports";
import { UsersPage } from "@/pages/users/Users";
import { TransactionsPage } from "@/pages/dashboard/Transactions";
import { AnalyticsPage } from "@/pages/dashboard/Analytics";
import { SettingsPage } from "@/pages/settings/Settings";
import { ProfileSettings } from "@/pages/settings/ProfileSettings";
import { UserRoleList } from "@/pages/users/UserRoleList";
import { UserRoleFormPage } from "@/pages/users/UserRoleFormPage";

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/app/dashboard" replace />} />

      <Route path="/auth" element={<AuthLayout />}>
        <Route index element={<WelcomePage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
      </Route>

      <Route
        path="/app"
        element={
          <ProtectedRoute requiredRole="ADMIN">
            <DashboardLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="reports" element={<ReportsPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="inventory" element={<InventoryPage />} />
        <Route path="calendar" element={<CalendarPage />} />

        <Route path="users" element={<UsersPage />} />
        <Route path="users/roles" element={<UserRoleList />} />
        <Route path="users/roles/new" element={<UserRoleFormPage />} />
        <Route path="users/roles/:id/edit" element={<UserRoleFormPage />} />
        <Route path="transactions" element={<TransactionsPage />} />
        <Route path="analytics" element={<AnalyticsPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="profile" element={<ProfileSettings />} />
        <Route index element={<Navigate to="/app/dashboard" replace />} />
      </Route>

      {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}
