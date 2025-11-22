import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';
import { attendeeService, type Attendee } from '../services/attendeeService';
import { speakerService, type Speaker } from '../services/speakerService';
import { sessionService, type Session } from '../services/sessionService';
import { adminService } from '../services/adminService';

const COLORS = ['#0ea5e9', '#3b82f6', '#6366f1', '#8b5cf6', '#a855f7', '#d946ef', '#ec4899', '#f43f5e', '#ef4444', '#f59e0b'];

const AdminPanel = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'attendees' | 'speakers' | 'sessions'>('attendees');
  const [attendees, setAttendees] = useState<Attendee[]>([]);
  const [speakers, setSpeakers] = useState<Speaker[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [stats, setStats] = useState<{ designation: string; count: number }[]>([]);
  const [loading, setLoading] = useState(true);
  const [showSpeakerModal, setShowSpeakerModal] = useState(false);
  const [showSessionModal, setShowSessionModal] = useState(false);
  const [editingSpeaker, setEditingSpeaker] = useState<Speaker | null>(null);
  const [editingSession, setEditingSession] = useState<Session | null>(null);
  const [speakerForm, setSpeakerForm] = useState<Omit<Speaker, 'id'>>({
    name: '',
    bio: '',
    avatar: '',
    linkedin: '',
    twitter: '',
  });
  const [sessionForm, setSessionForm] = useState<Omit<Session, 'id'>>({
    title: '',
    description: '',
    time: '',
    speakers: [],
  });

  useEffect(() => {
    fetchAllData();
  }, []);

  const fetchAllData = async () => {
    try {
      setLoading(true);
      const [attendeesData, speakersData, sessionsData, statsData] = await Promise.all([
        attendeeService.getAll(),
        speakerService.getAll(),
        sessionService.getAll(),
        adminService.getStats(),
      ]);
      setAttendees(Array.isArray(attendeesData) ? attendeesData : []);
      setSpeakers(Array.isArray(speakersData) ? speakersData : []);
      setSessions(Array.isArray(sessionsData) ? sessionsData : []);
      setStats(Array.isArray(statsData?.designationBreakdown) ? statsData.designationBreakdown : []);
    } catch (err: any) {
      if (err.response?.status === 401) {
        navigate('/');
      }
      console.error('Error fetching data:', err);
      // Set empty arrays on error to prevent crashes
      setAttendees([]);
      setSpeakers([]);
      setSessions([]);
      setStats([]);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await adminService.logout();
      navigate('/');
    } catch (err) {
      console.error('Logout error:', err);
      navigate('/');
    }
  };

  const handleDeleteAttendee = async (id: string) => {
    if (!confirm('Are you sure you want to delete this attendee?')) return;
    try {
      await attendeeService.delete(id);
      await fetchAllData();
    } catch (err) {
      console.error('Error deleting attendee:', err);
      alert('Failed to delete attendee');
    }
  };

  const handleSpeakerSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingSpeaker) {
        await speakerService.update(editingSpeaker.id!, speakerForm);
      } else {
        await speakerService.create(speakerForm);
      }
      setShowSpeakerModal(false);
      setEditingSpeaker(null);
      setSpeakerForm({ name: '', bio: '', avatar: '', linkedin: '', twitter: '' });
      await fetchAllData();
    } catch (err) {
      console.error('Error saving speaker:', err);
      alert('Failed to save speaker');
    }
  };

  const handleEditSpeaker = (speaker: Speaker) => {
    setEditingSpeaker(speaker);
    setSpeakerForm({
      name: speaker.name,
      bio: speaker.bio,
      avatar: speaker.avatar || '',
      linkedin: speaker.linkedin || '',
      twitter: speaker.twitter || '',
    });
    setShowSpeakerModal(true);
  };

  const handleDeleteSpeaker = async (id: string) => {
    if (!confirm('Are you sure you want to delete this speaker?')) return;
    try {
      await speakerService.delete(id);
      await fetchAllData();
    } catch (err) {
      console.error('Error deleting speaker:', err);
      alert('Failed to delete speaker');
    }
  };

  const handleSessionSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingSession) {
        await sessionService.update(editingSession.id!, sessionForm);
      } else {
        await sessionService.create(sessionForm);
      }
      setShowSessionModal(false);
      setEditingSession(null);
      setSessionForm({ title: '', description: '', time: '', speakers: [] });
      await fetchAllData();
    } catch (err) {
      console.error('Error saving session:', err);
      alert('Failed to save session');
    }
  };

  const handleEditSession = (session: Session) => {
    setEditingSession(session);
    setSessionForm({
      title: session.title || '',
      description: session.description || '',
      time: session.time || '',
      speakers: session.speakers || [],
    });
    setShowSessionModal(true);
  };

  const handleDeleteSession = async (id: string) => {
    if (!confirm('Are you sure you want to delete this session?')) return;
    try {
      await sessionService.delete(id);
      await fetchAllData();
    } catch (err) {
      console.error('Error deleting session:', err);
      alert('Failed to delete session');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading admin panel...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Admin Panel</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Logout
          </button>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        {/* Stats Chart */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white rounded-xl shadow-lg p-6 mb-8"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Attendee Breakdown by Designation</h2>
          {stats.length > 0 ? (
            <ResponsiveContainer width="100%" height={400}>
              <PieChart>
                <Pie
                  data={stats}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={(entry) => {
                    const data = entry.payload as { designation: string; count: number };
                    const percent = entry.percent ?? 0;
                    return `${data.designation}: ${(percent * 100).toFixed(0)}%`;
                  }}
                  outerRadius={120}
                  fill="#8884d8"
                  dataKey="count"
                >
                  {stats.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-500 text-center py-8">No data available</p>
          )}
        </motion.div>

        {/* Tabs */}
        <div className="bg-white rounded-xl shadow-lg mb-8">
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              {(['attendees', 'speakers', 'sessions'] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`px-6 py-4 font-semibold text-sm border-b-2 transition-colors ${
                    activeTab === tab
                      ? 'border-primary-600 text-primary-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  {tab.charAt(0).toUpperCase() + tab.slice(1)}
                </button>
              ))}
            </nav>
          </div>

          <div className="p-6">
            {/* Attendees Tab */}
            {activeTab === 'attendees' && (
              <div>
                <div className="mb-4 flex justify-between items-center">
                  <h3 className="text-xl font-bold text-gray-900">Attendees ({attendees.length})</h3>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Name</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Email</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Designation</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Registered</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Actions</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {attendees.map((attendee) => (
                        <tr key={attendee.id} className="hover:bg-gray-50">
                          <td className="px-4 py-3 text-sm text-gray-900">{attendee.name}</td>
                          <td className="px-4 py-3 text-sm text-gray-600">{attendee.email}</td>
                          <td className="px-4 py-3 text-sm text-gray-600">{attendee.designation}</td>
                          <td className="px-4 py-3 text-sm text-gray-600">
                            {attendee.createdAt
                              ? new Date(attendee.createdAt).toLocaleDateString()
                              : 'N/A'}
                          </td>
                          <td className="px-4 py-3 text-sm">
                            <button
                              onClick={() => handleDeleteAttendee(attendee.id!)}
                              className="text-red-600 hover:text-red-800 font-semibold"
                            >
                              Delete
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                  {attendees.length === 0 && (
                    <p className="text-center text-gray-500 py-8">No attendees registered yet</p>
                  )}
                </div>
              </div>
            )}

            {/* Speakers Tab */}
            {activeTab === 'speakers' && (
              <div>
                <div className="mb-4 flex justify-between items-center">
                  <h3 className="text-xl font-bold text-gray-900">Speakers ({speakers.length})</h3>
                  <button
                    onClick={() => {
                      setEditingSpeaker(null);
                      setSpeakerForm({ name: '', bio: '', avatar: '', linkedin: '', twitter: '' });
                      setShowSpeakerModal(true);
                    }}
                    className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Add Speaker
                  </button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {speakers.map((speaker) => (
                    <div key={speaker.id} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex items-start gap-4 mb-3">
                        {speaker.avatar ? (
                          <img
                            src={speaker.avatar}
                            alt={speaker.name}
                            className="w-16 h-16 rounded-full object-cover"
                          />
                        ) : (
                          <div className="w-16 h-16 rounded-full bg-primary-200 flex items-center justify-center text-primary-700 font-bold text-xl">
                            {speaker.name.charAt(0).toUpperCase()}
                          </div>
                        )}
                        <div className="flex-1">
                          <h4 className="font-bold text-gray-900">{speaker.name}</h4>
                          <p className="text-sm text-gray-600 line-clamp-2">{speaker.bio}</p>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleEditSpeaker(speaker)}
                          className="flex-1 px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors text-sm font-semibold"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteSpeaker(speaker.id!)}
                          className="flex-1 px-3 py-2 bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors text-sm font-semibold"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  ))}
                  {speakers.length === 0 && (
                    <p className="col-span-full text-center text-gray-500 py-8">No speakers added yet</p>
                  )}
                </div>
              </div>
            )}

            {/* Sessions Tab */}
            {activeTab === 'sessions' && (
              <div>
                <div className="mb-4 flex justify-between items-center">
                  <h3 className="text-xl font-bold text-gray-900">Sessions ({sessions.length})</h3>
                  <button
                    onClick={() => {
                      setEditingSession(null);
                      setSessionForm({ title: '', description: '', time: '', speakers: [] });
                      setShowSessionModal(true);
                    }}
                    className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Add Session
                  </button>
                </div>
                <div className="space-y-4">
                  {sessions.map((session) => (
                    <div key={session.id} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex justify-between items-start mb-2">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-2">
                            <span className="text-sm font-semibold text-primary-600 bg-primary-50 px-3 py-1 rounded-full">
                              {session.time || 'TBD'}
                            </span>
                            <h4 className="text-lg font-bold text-gray-900">{session.title}</h4>
                          </div>
                          <p className="text-gray-600 mb-2">{session.description}</p>
                          <p className="text-sm text-gray-500">
                            Speakers: {(session.speakers && session.speakers.length > 0) ? session.speakers.join(', ') : 'None'}
                          </p>
                        </div>
                        <div className="flex gap-2 ml-4">
                          <button
                            onClick={() => handleEditSession(session)}
                            className="px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors text-sm font-semibold"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => handleDeleteSession(session.id!)}
                            className="px-3 py-2 bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors text-sm font-semibold"
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                  {sessions.length === 0 && (
                    <p className="text-center text-gray-500 py-8">No sessions added yet</p>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Speaker Modal */}
      {showSpeakerModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className="bg-white rounded-xl shadow-2xl p-8 max-w-2xl w-full max-h-[90vh] overflow-y-auto"
          >
            <h3 className="text-2xl font-bold text-gray-900 mb-6">
              {editingSpeaker ? 'Edit Speaker' : 'Add Speaker'}
            </h3>
            <form onSubmit={handleSpeakerSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Name *</label>
                <input
                  type="text"
                  required
                  value={speakerForm.name}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, name: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Bio *</label>
                <textarea
                  required
                  value={speakerForm.bio}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, bio: e.target.value })}
                  rows={4}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Avatar URL</label>
                <input
                  type="url"
                  value={speakerForm.avatar}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, avatar: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">LinkedIn</label>
                  <input
                    type="url"
                    value={speakerForm.linkedin}
                    onChange={(e) => setSpeakerForm({ ...speakerForm, linkedin: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">Twitter</label>
                  <input
                    type="url"
                    value={speakerForm.twitter}
                    onChange={(e) => setSpeakerForm({ ...speakerForm, twitter: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
              <div className="flex gap-3 justify-end pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowSpeakerModal(false);
                    setEditingSpeaker(null);
                  }}
                  className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                >
                  {editingSpeaker ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}

      {/* Session Modal */}
      {showSessionModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className="bg-white rounded-xl shadow-2xl p-8 max-w-2xl w-full max-h-[90vh] overflow-y-auto"
          >
            <h3 className="text-2xl font-bold text-gray-900 mb-6">
              {editingSession ? 'Edit Session' : 'Add Session'}
            </h3>
            <form onSubmit={handleSessionSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Title *</label>
                <input
                  type="text"
                  required
                  value={sessionForm.title}
                  onChange={(e) => setSessionForm({ ...sessionForm, title: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Description *</label>
                <textarea
                  required
                  value={sessionForm.description}
                  onChange={(e) => setSessionForm({ ...sessionForm, description: e.target.value })}
                  rows={4}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Time *</label>
                <input
                  type="text"
                  required
                  value={sessionForm.time}
                  onChange={(e) => setSessionForm({ ...sessionForm, time: e.target.value })}
                  placeholder="e.g., 10:00 AM - 11:00 AM"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Speakers (IDs, comma-separated)</label>
                <input
                  type="text"
                  value={(sessionForm.speakers || []).join(', ')}
                  onChange={(e) =>
                    setSessionForm({
                      ...sessionForm,
                      speakers: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                    })
                  }
                  placeholder="speaker-id-1, speaker-id-2"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Available speakers: {speakers.length > 0 ? speakers.map((s) => s.id).filter(Boolean).join(', ') : 'None'}
                </p>
              </div>
              <div className="flex gap-3 justify-end pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowSessionModal(false);
                    setEditingSession(null);
                  }}
                  className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                >
                  {editingSession ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
};

export default AdminPanel;

