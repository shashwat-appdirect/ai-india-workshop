import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { sessionService, type SessionWithSpeakers } from '../services/sessionService';
import { speakerService } from '../services/speakerService';

const SessionsSection = () => {
  const [sessions, setSessions] = useState<SessionWithSpeakers[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [sessionsData, speakersData] = await Promise.all([
          sessionService.getAll(),
          speakerService.getAll(),
        ]);

        // Handle null responses
        const safeSessions = Array.isArray(sessionsData) ? sessionsData : [];
        const safeSpeakers = Array.isArray(speakersData) ? speakersData : [];

        // Enrich sessions with speaker details
        const enrichedSessions: SessionWithSpeakers[] = safeSessions.map((session) => {
          const speakerDetails = (session.speakers || [])
            .map((speakerId) => {
              const speaker = safeSpeakers.find((s) => s.id === speakerId);
              if (speaker && speaker.id) {
                return {
                  id: speaker.id,
                  name: speaker.name,
                  bio: speaker.bio,
                  ...(speaker.avatar && { avatar: speaker.avatar }),
                };
              }
              return null;
            })
            .filter((s): s is { id: string; name: string; bio: string; avatar?: string } => s !== null);
          
          return {
            ...session,
            speakerDetails,
          };
        });
        setSessions(enrichedSessions);
      } catch (err) {
        setError('Failed to load sessions. Please try again later.');
        console.error('Error fetching sessions:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <section id="sessions" className="py-20 bg-white">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">Sessions & Speakers</h2>
            <p className="text-gray-600">Loading sessions...</p>
          </div>
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section id="sessions" className="py-20 bg-white">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <p className="text-red-600">{error}</p>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section id="sessions" className="py-20 bg-white">
      <div className="container mx-auto px-4">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="text-center mb-12"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Sessions & Speakers
          </h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Explore our curated lineup of AI experts and thought-provoking sessions
          </p>
        </motion.div>

        {sessions.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No sessions available at the moment.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {sessions.map((session, index) => (
              <motion.div
                key={session.id || index}
                initial={{ opacity: 0, y: 30 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                className="bg-white rounded-xl shadow-lg hover:shadow-2xl transition-shadow duration-300 overflow-hidden border border-gray-100"
              >
                <div className="p-6">
                  <div className="mb-4">
                    <span className="text-sm font-semibold text-primary-600 bg-primary-50 px-3 py-1 rounded-full">
                      {session.time}
                    </span>
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-3">{session.title}</h3>
                  <p className="text-gray-600 mb-6 line-clamp-3">{session.description}</p>

                  {session.speakerDetails && session.speakerDetails.length > 0 && (
                    <div className="border-t pt-4">
                      <p className="text-sm font-semibold text-gray-700 mb-3">Speakers:</p>
                      <div className="flex flex-wrap gap-3">
                        {session.speakerDetails.map((speaker) => (
                          <div
                            key={speaker.id}
                            className="flex items-center gap-2 bg-gray-50 rounded-lg px-3 py-2"
                          >
                            {speaker.avatar ? (
                              <img
                                src={speaker.avatar}
                                alt={speaker.name}
                                className="w-8 h-8 rounded-full object-cover"
                              />
                            ) : (
                              <div className="w-8 h-8 rounded-full bg-primary-200 flex items-center justify-center text-primary-700 font-semibold text-sm">
                                {speaker.name.charAt(0).toUpperCase()}
                              </div>
                            )}
                            <span className="text-sm font-medium text-gray-700">
                              {speaker.name}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
};

export default SessionsSection;

