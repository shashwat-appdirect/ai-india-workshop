import api from './api';

export interface Speaker {
  id?: string;
  name: string;
  bio: string;
  avatar?: string;
  linkedin?: string;
  twitter?: string;
}

export const speakerService = {
  getAll: async (): Promise<Speaker[]> => {
    const response = await api.get<Speaker[]>('/speakers');
    return Array.isArray(response.data) ? response.data : [];
  },

  create: async (speaker: Omit<Speaker, 'id'>): Promise<Speaker> => {
    const response = await api.post<Speaker>('/speakers', speaker);
    return response.data;
  },

  update: async (id: string, speaker: Partial<Speaker>): Promise<Speaker> => {
    const response = await api.put<Speaker>(`/speakers/${id}`, speaker);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/speakers/${id}`);
  },
};

