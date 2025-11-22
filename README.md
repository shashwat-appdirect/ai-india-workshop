# AppDirect India AI Workshop Registration SPA

A production-ready React SPA with Golang backend for event registration, featuring session management, speaker profiles, and admin dashboard.

## Architecture

- **Frontend**: React + Vite + TypeScript + Tailwind CSS
- **Backend**: Golang (Gin framework) with REST APIs
- **Database**: Firebase Firestore with Application Default Credentials (Cloud Run) or service account (local)
- **Deployment**: Unified Docker container (frontend + backend) ready for Google Cloud Run

## Features

- 🎯 Hero section with animated CTAs
- 📅 Sessions & Speakers grid with responsive design
- 📝 Registration form with live attendee count
- 📍 Location section with embedded Google Maps
- 🔐 Password-protected admin panel
- 📊 Analytics dashboard with pie charts
- ✏️ CRUD operations for speakers and sessions
- 👥 Attendee management

## Prerequisites

- Node.js 20+ and npm
- Go 1.23+
- Docker and Docker Compose (optional)
- Firebase service account JSON file (for local development)
- Google Cloud Project with Firestore enabled (for Cloud Run deployment)

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd ai-india-workshop
```

### 2. Backend Setup

```bash
cd backend
go mod download
```

### 3. Frontend Setup

```bash
cd frontend
npm install
```

### 4. Environment Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Update the following variables:
- `FIREBASE_SERVICE_ACCOUNT_PATH`: Path to your Firebase service account JSON (optional for Cloud Run, required for local)
- `FIRESTORE_SUBCOLLECTION_ID`: Your Firestore subcollection identifier
- `ADMIN_PASSWORD`: Admin login password
- `SESSION_SECRET`: Session secret (min 32 characters)
- `FRONTEND_URL`: Frontend URL for CORS (defaults to http://localhost:5173)
- `STATIC_DIR`: Directory for static files (set automatically in Docker, optional for local)

For frontend, copy `frontend/.env.example` to `frontend/.env`:

```bash
cd frontend
cp .env.example .env
```

### 5. Firebase Setup

1. Create a Firebase project
2. Enable Firestore Database
3. Download service account JSON
4. Place it in the project root as `firebase-service-account.json`
5. Update `FIRESTORE_SUBCOLLECTION_ID` in `.env`

## Running Locally

### Development Mode

**Backend:**
```bash
cd backend
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend
npm run dev
```

The application will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

### Using Makefile

The project includes a Makefile for easy execution:

```bash
# Install all dependencies
make install

# Build both frontend and backend
make build

# Run both services in development mode
make run

# Run tests
make test

# Build Docker image
make docker-build

# Start with Docker Compose
make docker-up

# Stop Docker Compose services
make docker-down

# View Docker logs
make docker-logs
```

### Docker Compose

```bash
docker-compose up --build
```

The application will be available at:
- Application: http://localhost:8080 (frontend + backend unified)

### Docker Build

Build the unified Docker image:

```bash
docker build -t ai-india-workshop:latest .
```

Run the container:

```bash
docker run -p 8080:8080 \
  -e FIRESTORE_SUBCOLLECTION_ID=ai-india-workshop-2024 \
  -e ADMIN_PASSWORD=your-password \
  -e SESSION_SECRET=your-secret-min-32-chars \
  ai-india-workshop:latest
```

For local development with service account file:

```bash
docker run -p 8080:8080 \
  -e FIREBASE_SERVICE_ACCOUNT_PATH=/app/firebase-service-account.json \
  -e FIRESTORE_SUBCOLLECTION_ID=ai-india-workshop-2024 \
  -e ADMIN_PASSWORD=your-password \
  -e SESSION_SECRET=your-secret-min-32-chars \
  -v $(pwd)/firebase-service-account.json:/app/firebase-service-account.json:ro \
  ai-india-workshop:latest
```

## Google Cloud Run Deployment

The application is configured to deploy to Google Cloud Run without requiring a Firebase service account file. It uses Application Default Credentials (ADC) provided by Cloud Run.

### Prerequisites

1. Google Cloud Project with billing enabled
2. Firestore Database enabled
3. Cloud Run API enabled
4. Service account with Firestore permissions

### Deployment Steps

1. **Build and push the Docker image to Google Container Registry:**

```bash
# Set your project ID
export PROJECT_ID=your-project-id
export SERVICE_NAME=ai-india-workshop

# Build and tag the image
docker build -t gcr.io/$PROJECT_ID/$SERVICE_NAME:latest .

# Push to GCR
docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:latest
```

2. **Deploy to Cloud Run:**

```bash
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="FIRESTORE_SUBCOLLECTION_ID=ai-india-workshop-2024" \
  --set-env-vars="ADMIN_PASSWORD=your-secure-password" \
  --set-env-vars="SESSION_SECRET=your-secret-min-32-characters-long" \
  --set-env-vars="PORT=8080" \
  --set-env-vars="STATIC_DIR=/app/static"
```

3. **Set up Service Account (for Firestore access):**

```bash
# Create a service account (if not exists)
gcloud iam service-accounts create ai-workshop-sa \
  --display-name="AI Workshop Service Account"

# Grant Firestore permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:ai-workshop-sa@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/datastore.user"

# Attach service account to Cloud Run service
gcloud run services update $SERVICE_NAME \
  --service-account=ai-workshop-sa@$PROJECT_ID.iam.gserviceaccount.com \
  --region us-central1
```

### Important Notes for Cloud Run

- **No Firebase Service Account File Required**: The application automatically uses Application Default Credentials when `FIREBASE_SERVICE_ACCOUNT_PATH` is not set
- **Environment Variables**: Set all required environment variables via Cloud Run configuration
- **Port**: Cloud Run sets the `PORT` environment variable automatically, but you can override it
- **Static Files**: The `STATIC_DIR` environment variable is set in the Dockerfile to `/app/static`
- **CORS**: Update `FRONTEND_URL` environment variable to your Cloud Run service URL

## API Endpoints

### Public Endpoints

- `POST /api/attendees` - Register new attendee
- `GET /api/attendees/count` - Get attendee count
- `GET /api/speakers` - List all speakers
- `GET /api/sessions` - List all sessions
- `POST /api/admin/login` - Admin login
- `POST /api/admin/logout` - Admin logout

### Admin Endpoints (Requires Authentication)

- `GET /api/attendees` - List all attendees
- `DELETE /api/attendees/:id` - Delete attendee
- `POST /api/speakers` - Create speaker
- `PUT /api/speakers/:id` - Update speaker
- `DELETE /api/speakers/:id` - Delete speaker
- `POST /api/sessions` - Create session
- `PUT /api/sessions/:id` - Update session
- `DELETE /api/sessions/:id` - Delete session
- `GET /api/admin/stats` - Get statistics

## Project Structure

```
ai-india-workshop/
├── frontend/              # React SPA
│   ├── src/
│   │   ├── components/    # Reusable components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API service layer
│   │   └── App.tsx
│   └── Dockerfile         # Frontend-only Dockerfile (legacy)
├── backend/               # Golang REST API
│   ├── cmd/server/        # Server entry point
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers
│   │   ├── models/        # Data models
│   │   ├── repository/    # Firestore repository
│   │   └── middleware/    # Auth middleware
│   └── Dockerfile         # Backend-only Dockerfile (legacy)
├── Dockerfile             # Unified multi-stage Dockerfile (production)
├── docker-compose.yml     # Docker Compose configuration
├── Makefile              # Makefile for easy execution
├── .dockerignore         # Docker ignore file
├── .env.example          # Environment variable template
└── README.md
```

## Security Notes

- Never commit `.env` files or Firebase service account JSON
- Use strong passwords and secrets in production
- Configure CORS properly for production
- Use HTTPS in production
- Regularly rotate secrets and passwords

## License

Copyright © 2024 AppDirect India

