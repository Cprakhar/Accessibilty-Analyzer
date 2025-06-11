# Accessibility Analyser

A full-stack web application for analyzing the accessibility of web pages, providing improvement suggestions, and generating analytics dashboards.

## Project Structure

```
Accessibilty Analyser/
├── backend/   # Go (Gin) API server, MongoDB, JWT auth, logging, etc.
├── frontend/  # Next.js, React, Tailwind CSS UI
├── tech-spec.md
```

## Features
- User authentication (JWT)
- Accessibility analysis (Lighthouse/axe-core integration planned)
- Suggestions via LLM (Hugging Face API planned)
- Analytics dashboards (Tableau integration planned)
- Audit logging
- CORS, input validation, and security best practices

## Getting Started

### Prerequisites
- Go 1.20+
- Node.js 18+
- MongoDB (local or Atlas free tier)

### Backend Setup
See [backend/README.md](backend/README.md) for detailed instructions.

### Frontend Setup
From the `frontend/` directory:
```bash
npm install
npm run dev
```
The app will be available at http://localhost:3000

### Environment Variables
- Copy `backend/.env.example` to `backend/.env` and fill in your values
- Set up MongoDB and JWT secret as described in backend/README.md

## Monorepo Notes
- This repository contains both backend and frontend code for easier development and deployment.
- Use a single git repository at the project root.

## License
MIT
