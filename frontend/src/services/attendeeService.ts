import api from './api';

export interface Attendee {
  id?: string;
  name: string;
  email: string;
  designation: string;
  createdAt?: string;
}

export const attendeeService = {
  register: async (attendee: Omit<Attendee, 'id' | 'createdAt'>): Promise<Attendee> => {
    const response = await api.post<Attendee>('/attendees', attendee);
    return response.data;
  },

  getCount: async (): Promise<number> => {
    const response = await api.get<{ count: number }>('/attendees/count');
    return response.data.count;
  },

  getAll: async (): Promise<Attendee[]> => {
    const response = await api.get<Attendee[]>('/attendees');
    return Array.isArray(response.data) ? response.data : [];
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/attendees/${id}`);
  },
};

