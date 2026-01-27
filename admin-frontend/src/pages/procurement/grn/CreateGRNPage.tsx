import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useAuth } from "@/contexts/AuthContext";
import { procurementService } from "@/services/procurement.service";
import { locationService } from "@/services/location.service";
import { PurchaseOrder } from "@/types/procurement.types";
import { Location } from "@/types/product.types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CalendarIcon, ArrowLeft, Save, AlertCircle } from "lucide-react";
import { format } from "date-fns";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import { useToast } from "@/components/ui/use-toast";

// Schema
const createGRNSchema = z.object({
  purchase_order_id: z.string().min(1, "Purchase Order ID is required"),
  location_id: z.string().min(1, "Location is required"),
  receipt_date: z.date({
    required_error: "Receipt date is required",
  }),
  notes: z.string().optional(),
  items: z
    .array(
      z.object({
        product_id: z.string(),
        sku: z.string().optional(),
        description: z.string().optional(),
        ordered_quantity: z.number(), // for reference
        received_quantity: z.coerce.number().min(0, "Must be positive"),
        batch_number: z.string().optional(),
        expiry_date: z.string().optional(), // simplified as string for now, or Date
      }),
    )
    .min(1, "At least one item is required"),
});

type CreateGRNValues = z.infer<typeof createGRNSchema>;

export function CreateGRNPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const poId = searchParams.get("po_id");
  const { user } = useAuth();
  const { toast } = useToast();

  const [po, setPo] = useState<PurchaseOrder | null>(null);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const form = useForm<CreateGRNValues>({
    resolver: zodResolver(createGRNSchema),
    defaultValues: {
      receipt_date: new Date(),
      items: [],
      purchase_order_id: poId || "",
    },
  });

  const { fields } = useFieldArray({
    control: form.control,
    name: "items",
  });

  useEffect(() => {
    if (poId) {
      fetchPO(poId);
    }
    if (user?.organization_id) {
      fetchLocations(user.organization_id);
    }
  }, [poId, user?.organization_id]);

  const fetchLocations = async (orgId: string) => {
    try {
      const locations = await locationService.getOrganizationLocations(orgId);
      setLocations(locations || []);
    } catch (error) {
      console.error("Failed to fetch locations", error);
    }
  };

  const fetchPO = async (id: string) => {
    try {
      setLoading(true);
      const data = await procurementService.getPurchaseOrder(id);
      setPo(data);
      form.setValue("purchase_order_id", data.id);
      // Pre-fill location if PO has delivery location
      if (data.delivery_location_id) {
        form.setValue("location_id", data.delivery_location_id);
      }

      // Pre-fill items
      const grnItems = data.items.map((item) => ({
        product_id: item.product_id,
        sku: item.sku,
        description: item.description,
        ordered_quantity: item.quantity,
        received_quantity: item.quantity, // Default to receiving full amount
        batch_number: "",
        expiry_date: "",
      }));

      form.setValue("items", grnItems);
    } catch (error) {
      console.error("Failed to fetch PO", error);
      toast({
        title: "Error",
        description: "Failed to load Purchase Order details",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const onSubmit = async (values: CreateGRNValues) => {
    if (!user?.organization_id) return;
    try {
      setSubmitting(true);

      const payload = {
        purchase_order_id: values.purchase_order_id,
        location_id: values.location_id,
        receipt_date: values.receipt_date.toISOString(),
        notes: values.notes,
        items: values.items.map((item) => ({
          product_id: item.product_id,
          ordered_quantity: Number(item.ordered_quantity),
          received_quantity: Number(item.received_quantity),
          batch_number: item.batch_number,
          ...(item.expiry_date ? { expiry_date: item.expiry_date } : {}),
        })),
      };

      await procurementService.createGRN(user.organization_id, payload);

      toast({
        title: "Success",
        description: "Goods Receipt Note created successfully",
      });
      navigate("/app/procurement/grn");
    } catch (error) {
      console.error("Failed to create GRN", error);
      toast({
        title: "Error",
        description: "Failed to create GRN",
        variant: "destructive",
      });
    } finally {
      setSubmitting(false);
    }
  };

  if (!poId) {
    return (
      <div className="p-8 text-center space-y-4">
        <AlertCircle className="h-10 w-10 text-muted-foreground mx-auto" />
        <h3 className="text-lg font-medium">No Purchase Order Selected</h3>
        <p className="text-muted-foreground">
          Please select a Purchase Order to receive goods against.
        </p>
        <Button onClick={() => navigate("/app/procurement/orders")}>
          Go to Orders
        </Button>
      </div>
    );
  }

  if (loading) {
    return <div className="p-8 text-center">Loading PO details...</div>;
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Receive Goods</h2>
            <p className="text-muted-foreground">
              Create GRN for PO #{po?.po_number}
            </p>
          </div>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="md:col-span-2 space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Items Received</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-6">
                    {fields.map((field, index) => (
                      <div
                        key={field.id}
                        className="flex flex-col gap-4 border-b pb-4 last:border-0 last:pb-0"
                      >
                        <div className="flex justify-between items-start">
                          <div>
                            <p className="font-medium">
                              Product{" "}
                              {form
                                .getValues(`items.${index}.product_id`)
                                .slice(-6)}
                            </p>
                            <p className="text-sm text-muted-foreground">
                              Ordered:{" "}
                              {form.getValues(
                                `items.${index}.ordered_quantity`,
                              )}
                            </p>
                          </div>
                        </div>

                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                          <FormField
                            control={form.control}
                            name={`items.${index}.received_quantity`}
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Received Qty</FormLabel>
                                <FormControl>
                                  <Input type="number" min="0" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                          <FormField
                            control={form.control}
                            name={`items.${index}.batch_number`}
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Batch # (Opt)</FormLabel>
                                <FormControl>
                                  <Input {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Receipt Details</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <FormField
                    control={form.control}
                    name="location_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Receiving Location</FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          defaultValue={field.value}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select location" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {locations.map((loc) => (
                              <SelectItem key={loc.id} value={loc.id}>
                                {loc.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="receipt_date"
                    render={({ field }) => (
                      <FormItem className="flex flex-col">
                        <FormLabel>Receipt Date</FormLabel>
                        <Popover>
                          <PopoverTrigger asChild>
                            <FormControl>
                              <Button
                                variant={"outline"}
                                className={cn(
                                  "w-full pl-3 text-left font-normal",
                                  !field.value && "text-muted-foreground",
                                )}
                              >
                                {field.value ? (
                                  format(field.value, "PPP")
                                ) : (
                                  <span>Pick a date</span>
                                )}
                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                              </Button>
                            </FormControl>
                          </PopoverTrigger>
                          <PopoverContent className="w-auto p-0" align="start">
                            <Calendar
                              mode="single"
                              selected={field.value}
                              onSelect={field.onChange}
                              disabled={
                                (date) => date > new Date() // Cannot be in future ideally
                              }
                              initialFocus
                            />
                          </PopoverContent>
                        </Popover>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="notes"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Notes</FormLabel>
                        <FormControl>
                          <Input placeholder="Receipt notes..." {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <Button
                    type="submit"
                    className="w-full"
                    disabled={submitting}
                  >
                    {submitting ? (
                      "Creating..."
                    ) : (
                      <>
                        <Save className="mr-2 h-4 w-4" /> Create GRN
                      </>
                    )}
                  </Button>
                </CardContent>
              </Card>
            </div>
          </div>
        </form>
      </Form>
    </div>
  );
}
