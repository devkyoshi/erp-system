import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string;
  requiredPermissions?: string[];
}

export function ProtectedRoute({
  children,
  requiredRole,
  requiredPermissions,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-brand border-t-transparent"></div>
          <p className="mt-4 text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    // Redirect to login page with return url
    return <Navigate to="/auth/login" state={{ from: location }} replace />;
  }

  // Check role if required
//   if (requiredRole && user?.role !== requiredRole) {
//     return (
//       <div className="flex h-screen items-center justify-center">
//         <div className="text-center">
//           <h1 className="text-2xl font-semibold">Access Denied</h1>
//           <p className="mt-2 text-muted-foreground">
//             You don't have permission to access this page.
//           </p>
//         </div>
//       </div>
//     );
//   }

  // Check permissions if required
//   if (requiredPermissions && requiredPermissions.length > 0) {
//     const hasPermission = requiredPermissions.every((permission) =>
//       user?.permissions?.includes(permission),
//     );

//     if (!hasPermission) {
//       return (
//         <div className="flex h-screen items-center justify-center">
//           <div className="text-center">
//             <h1 className="text-2xl font-semibold">Access Denied</h1>
//             <p className="mt-2 text-muted-foreground">
//               You don't have the required permissions to access this page.
//             </p>
//           </div>
//         </div>
//       );
//     }
//   }

  return <>{children}</>;
}
