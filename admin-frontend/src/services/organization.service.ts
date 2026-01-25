import axiosInstance from "@/lib/axios";
import { ApiResponse } from "@/types/global.types";

export interface OrganizationOption {
  id: string;
  name: string;
}

class OrganizationService {
  /**
   * Get all organizations (options only)
   */
  async getOrganizations(): Promise<OrganizationOption[]> {
    const response = await axiosInstance.get<ApiResponse<OrganizationOption[]>>(
      "/organizations/options",
    );

    // The response structure is:
    // response.data (ApiResponse) -> data (OrganizationOption[])
    return response.data.data || [];
  }
}

export const organizationService = new OrganizationService();
