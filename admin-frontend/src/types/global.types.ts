// Backend API response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
  error?: string;
  code?: string;
}