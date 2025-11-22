import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Hero from './components/Hero';
import SessionsSection from './components/SessionsSection';
import RegistrationForm from './components/RegistrationForm';
import LocationSection from './components/LocationSection';
import Footer from './components/Footer';
import AdminPanel from './pages/AdminPanel';

function App() {
  return (
    <Router>
      <div className="min-h-screen flex flex-col">
        <Routes>
          <Route
            path="/"
            element={
              <>
                <Hero />
                <SessionsSection />
                <RegistrationForm />
                <LocationSection />
                <Footer />
              </>
            }
          />
          <Route path="/admin" element={<AdminPanel />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
