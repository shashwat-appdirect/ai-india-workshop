import api from './api';

export interface AdminStats {
  designationBreakdown: Array<{
    designation: string;
    count: number;
  }>;
}

export const adminService = {
  login: async (password: string): Promise<{ success: boolean }> => {
    const response = await api.post<{ success: boolean }>('/admin/login', { password });
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/admin/logout');
  },

  getStats: async (): Promise<AdminStats> => {
    const response = await api.get<AdminStats>('/admin/stats');
    return response.data;
  },
};


