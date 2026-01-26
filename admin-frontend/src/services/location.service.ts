import axiosInstance from "@/lib/axios";
import { Location } from "@/types/product.types";
import { ApiResponse, PaginatedResponse } from "@/types/global.types";

export interface CreateLocationRequest {
  name: string;
  code: string;
  email: string;
  type: string;
  address?: {
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
  };
  is_active: boolean;
}

export interface UserAccessData {
  locations: Location[];
  organizations: any[];
  companies: any[];
}

class LocationService {
  /**
   * Get all locations accessible by the current user
   */
  async getUserLocations(): Promise<Location[]> {
    const response =
      await axiosInstance.get<ApiResponse<UserAccessData>>("/users/me/access");
    return response.data.data.locations || [];
  }

  /**
   * Get user access data (locations, companies, organizations)
   */
  async getUserAccess(): Promise<UserAccessData> {
    const response =
      await axiosInstance.get<ApiResponse<UserAccessData>>("/users/me/access");
    return response.data.data;
  }

  /**
   * List locations for a company
   */
  async getLocations(
    companyId: string,
    params?: { page?: number; limit?: number },
  ): Promise<PaginatedResponse<Location>> {
    const response = await axiosInstance.get<PaginatedResponse<Location>>(
      `/companies/${companyId}/locations`,
      { params },
    );
    return response.data;
  }

  /**
   * Get a location by ID
   */
  async getLocation(id: string): Promise<Location> {
    const response = await axiosInstance.get<ApiResponse<Location>>(
      `/locations/${id}`,
    );
    return response.data.data;
  }

  /**
   * Create a new location
   */
  async createLocation(
    companyId: string,
    data: CreateLocationRequest,
  ): Promise<Location> {
    const response = await axiosInstance.post<ApiResponse<Location>>(
      `/companies/${companyId}/locations`,
      data,
    );
    return response.data.data;
  }
}

export const locationService = new LocationService();
