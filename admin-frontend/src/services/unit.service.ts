import axiosInstance from "@/lib/axios";
import { Unit } from "@/types/product.types";
import { ApiResponse } from "@/types/global.types";

class UnitService {
  /**
   * Get all units
   */
  async getUnits(): Promise<Unit[]> {
    const response = await axiosInstance.get<ApiResponse<Unit[]>>("/units");
    return response.data.data;
  }
}

export const unitService = new UnitService();
