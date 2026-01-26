import { useNavigate, useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

export function ProductFormPage() {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditMode = !!id;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/products")}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            {isEditMode ? "Edit Product" : "Add New Product"}
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditMode
              ? "Update product information and pricing"
              : "Create a new product in your catalog"}
          </p>
        </div>
      </div>

      {/* Placeholder Content */}
      <div className="rounded-lg border border-border bg-elevated p-12 text-center">
        <div className="mx-auto max-w-md space-y-4">
          <div className="mx-auto h-16 w-16 rounded-full bg-muted flex items-center justify-center">
            <span className="text-2xl">🚧</span>
          </div>
          <h3 className="text-lg font-semibold">Product Form Coming Soon</h3>
          <p className="text-sm text-muted-foreground">
            The product creation and editing form is currently under
            development. This will include tabs for basic information, pricing,
            inventory settings, and more.
          </p>
          <Button onClick={() => navigate("/app/products")}>
            Back to Products
          </Button>
        </div>
      </div>
    </div>
  );
}
