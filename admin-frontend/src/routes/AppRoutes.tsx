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
import { ProductsPage } from "@/pages/products/Products";
import { ProductFormPage } from "@/pages/products/ProductForm";
import { CategoriesPage } from "@/pages/categories/CategoryList";
import { CategoryFormPage } from "@/pages/categories/CategoryForm";
import { BrandsPage } from "@/pages/brands/BrandList";
import { BrandFormPage } from "@/pages/brands/BrandForm";
import CompanyList from "@/pages/settings/companies/CompanyList";
import CompanyDetail from "@/pages/settings/companies/CompanyDetail";


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

        <Route path="products" element={<ProductsPage />} />
        <Route path="products/new" element={<ProductFormPage />} />
        <Route path="products/:id/edit" element={<ProductFormPage />} />

        <Route path="categories" element={<CategoriesPage />} />
        <Route path="categories/new" element={<CategoryFormPage />} />
        <Route path="categories/:id/edit" element={<CategoryFormPage />} />

        <Route path="brands" element={<BrandsPage />} />
        <Route path="brands/new" element={<BrandFormPage />} />
        <Route path="brands/:id/edit" element={<BrandFormPage />} />

        <Route path="transactions" element={<TransactionsPage />} />
        <Route path="analytics" element={<AnalyticsPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="companies" element={<CompanyList />} />
        <Route path="companies/:id" element={<CompanyDetail />} />
        <Route path="profile" element={<ProfileSettings />} />
        <Route index element={<Navigate to="/app/dashboard" replace />} />
      </Route>

      {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}
