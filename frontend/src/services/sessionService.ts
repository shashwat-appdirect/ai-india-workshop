import api from './api';

export interface Session {
  id?: string;
  title: string;
  description: string;
  time: string;
  speakers: string[]; // Speaker IDs
}

export interface SessionWithSpeakers extends Session {
  speakerDetails?: Array<{
    id: string;
    name: string;
    bio: string;
    avatar?: string;
  }>;
}

export const sessionService = {
  getAll: async (): Promise<SessionWithSpeakers[]> => {
    const response = await api.get<SessionWithSpeakers[]>('/sessions');
    return Array.isArray(response.data) ? response.data : [];
  },

  create: async (session: Omit<Session, 'id'>): Promise<Session> => {
    const response = await api.post<Session>('/sessions', session);
    return response.data;
  },

  update: async (id: string, session: Partial<Session>): Promise<Session> => {
    const response = await api.put<Session>(`/sessions/${id}`, session);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/sessions/${id}`);
  },
};

