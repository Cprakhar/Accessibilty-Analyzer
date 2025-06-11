## Backend Technical Specifications

### API Endpoints

#### 1. Accessibility Analysis
- POST /api/analyze
  - Request Body: `{ "url": "String" }` or `{ "html": "String" }`
  - Response: `{ "reportId": "ObjectId", "analysisResults": {}, "suggestions": ["String"], "createdAt": "Date" }`

#### 2. Reports Management
- GET /api/reports
  - Returns: `[{ "_id": "ObjectId", "url": "String", "createdAt": "Date", "status": "String" }, ...]`
- GET /api/reports/:id
  - Returns: `{ "_id": "ObjectId", "userId": "ObjectId", "url": "String", "htmlSnapshot": "String", "analysisResults": {}, "createdAt": "Date", "updatedAt": "Date", "status": "String" }`
- DELETE /api/reports/:id
  - Request Param: `id` (report id)
  - Response: `{ "success": true, "message": "Report deleted." }`

#### 3. Dashboard Data
- GET /api/dashboard/summary
  - Returns: `{ "summaryData": {}, "createdAt": "Date", "updatedAt": "Date" }`

#### 4. Suggestions & Insights
- GET /api/reports/:id/suggestions
  - Returns: `{ "reportId": "ObjectId", "suggestions": ["String"], "createdAt": "Date" }`

#### 5. User Authentication (if needed)
- POST /api/auth/login
  - Request Body: `{ "email": "String", "password": "String" }`
  - Response: `{ "token": "String", "user": { "_id": "ObjectId", "email": "String", "name": "String" } }`
- POST /api/auth/register
  - Request Body: `{ "email": "String", "password": "String", "name": "String" }`
  - Response: `{ "token": "String", "user": { "_id": "ObjectId", "email": "String", "name": "String" } }`
- POST /api/auth/logout
  - Request: Auth token in header
  - Response: `{ "success": true, "message": "Logged out." }`
- GET /api/auth/me
  - Request: Auth token in header
  - Response: `{ "_id": "ObjectId", "email": "String", "name": "String", "createdAt": "Date" }`

#### 6. Health Check
- GET /api/health
  - Returns: `{ "status": "ok", "timestamp": "Date" }`

---

### Database Schemas

#### User
Stores user authentication and profile information.
```json
{
  "_id": "ObjectId",
  "email": "String",
  "passwordHash": "String",
  "name": "String",
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

#### Report
Stores results of each accessibility analysis.
```json
{
  "_id": "ObjectId",
  "userId": "ObjectId",
  "url": "String",
  "htmlSnapshot": "String",
  "analysisResults": {},
  "createdAt": "Date",
  "updatedAt": "Date",
  "status": "String"
}
```

#### Suggestion
Stores LLM-generated suggestions for each report.
```json
{
  "_id": "ObjectId",
  "reportId": "ObjectId",
  "suggestions": ["String"],
  "createdAt": "Date"
}
```

#### DashboardSummary (optional)
Stores cached dashboard/Power BI summary data.
```json
{
  "_id": "ObjectId",
  "userId": "ObjectId",
  "summaryData": {},
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

---

### API Response Format & Error Handling

All API responses should follow a standardized JSON structure to ensure consistency and ease of error handling on the frontend. Example response format:

```json
{
  "success": true,
  "message": "Operation completed successfully.",
  "data": { /* response data or null */ },
  "error": null
}
```

- `success`: Boolean indicating if the request was successful.
- `message`: Human-readable message for the client.
- `data`: The actual response payload (object, array, or null).
- `error`: Error details (object or null), e.g. `{ "code": "VALIDATION_ERROR", "details": "Email is required." }`

All error responses should set `success: false`, provide a relevant `message`, and include an `error` object with details.

---

### Security & Rate Limiting

#### Rate Limiting
- Implement per-user and per-IP rate limiting for all API endpoints to prevent abuse and denial-of-service attacks.
- Use Redis or in-memory store for distributed rate limiting.
- Example: Limit to 100 requests per minute per user/IP.

#### Input Validation
- Validate all incoming request data (body, query, params) for type, format, and length using Go validation libraries (e.g., go-playground/validator).
- Sanitize user input to prevent injection attacks.

#### CORS (Cross-Origin Resource Sharing)
- Restrict API access to trusted frontend origins only.
- Use CORS middleware in Gin to set allowed origins, methods, and headers.

#### CSRF (Cross-Site Request Forgery) Protection
- For cookie-based authentication, implement CSRF tokens for all state-changing requests.
- Use CSRF middleware if not using JWT.

#### XSS (Cross-Site Scripting) Protection
- Sanitize and escape all user-generated content before rendering on the frontend.
- Validate and sanitize input on the backend.

#### Authentication & Authorization
- Use JWT (JSON Web Tokens) for authentication and stateless session management.
- Issue a JWT on successful login or registration, and require it in the Authorization header for protected endpoints.
- Use short-lived access tokens and refresh tokens for improved security.

#### HTTPS Enforcement
- Serve all API endpoints over HTTPS in production to encrypt data in transit.

---

### Audit Logging

#### Logging Key Actions
- Log important user and system actions for security, compliance, and debugging purposes.
- Actions to log include: user login/logout, registration, failed/successful authentication attempts, report creation/deletion, analysis runs, and permission errors.
- Store logs in a centralized location (e.g., file, database, or log management service like ELK, Loki, or a cloud provider's logging solution).
- Include relevant metadata in each log entry: timestamp, user ID (if available), action type, status (success/failure), and request details (IP, endpoint, etc.).
- Ensure logs do not contain sensitive information (e.g., passwords, tokens).
- Regularly review and rotate logs for security and compliance.

---

### Background Jobs & Queue Processing

#### Asynchronous Processing
- Accessibility analysis and report generation can be time-consuming operations.
- Use background job processing to handle these tasks asynchronously, improving API responsiveness and user experience.

#### Job Queue
- Implement a job queue (e.g., using Redis) to manage and process background jobs such as:
  - Running accessibility audits
  - Generating analytics reports
  - Sending notifications or emails
- Each job should include relevant metadata (e.g., user ID, report ID, job type, status, timestamps).

#### Worker Service
- Run one or more worker processes/services that consume jobs from the queue and execute the required tasks.
- Workers should update job status and results in the database upon completion or failure.

#### Monitoring & Reliability
- Monitor the job queue and worker health to ensure timely processing and detect failures.
- Implement retry logic and dead-letter queues for failed jobs.

---

### Integration Details

#### Tableau Integration

- **Purpose:** Generate and embed analytics reports and dashboards for users.
- **Integration Method:** Use the Tableau REST API or Tableau JavaScript API.
- **Authentication:** Authenticate backend with Tableau using Personal Access Tokens or OAuth.
- **Workflow:**
  1. When a report is generated or updated, the backend sends relevant data to Tableau via API or triggers a data refresh.
  2. The backend retrieves the embed URL or view ID from Tableau.
  3. The frontend uses the embed URL to display the dashboard/report to the user.
- **Backend Responsibilities:**
  - Prepare and format data for Tableau (e.g., push to a data source or trigger extract refresh).
  - Call Tableau REST API endpoints to refresh data sources, manage workbooks, and retrieve embed URLs.
  - Store/report embed URLs or IDs in the database for frontend use.
  - Handle authentication token refresh for Tableau API access.

#### LLM (Large Language Model) Integration

- **Purpose:** Generate insights and suggestions based on accessibility analysis results.
- **Integration Method:** Use HTTP API calls to an LLM service e.g., Hugging Face Inference API.
- **Workflow:**
  1. After an accessibility analysis, the backend formats the results and sends them to the Hugging Face Inference API.
  2. The LLM returns suggestions or insights.
  3. The backend stores these suggestions in the database and includes them in the API response.
- **Backend Responsibilities:**
  - Format and sanitize analysis results before sending to the Hugging Face API.
  - Call the Hugging Face Inference API endpoint with the relevant prompt/data.
  - Handle API authentication (API keys, tokens).
  - Store LLM responses (suggestions) in the Suggestion collection.
  - Implement error handling and retries for failed API requests.

---

### Backend Directory Structure

Recommended directory structure for the backend:

```
backend/
├── api/                # Route handlers/controllers for each API endpoint
├── models/             # Database schema definitions (User, Report, Suggestion, etc.)
├── services/           # Business logic, integrations (Tableau, LLM, etc.)
├── jobs/               # Background job definitions and workers
├── utils/              # Utility functions (validation, logging, etc.)
├── middleware/         # Gin middleware (auth, rate limiting, CORS, etc.)
├── config/             # Configuration files (env, API keys, etc.)
├── main.go             # Application entry point
└── README.md           # Backend-specific documentation
```

- Place route handlers in `api/` (e.g., `api/report.go`, `api/auth.go`).
- Place MongoDB schema definitions in `models/` (e.g., `models/user.go`).
- Place Tableau and LLM integration logic in `services/`.
- Place background job logic and workers in `jobs/`.
- Place utility functions in `utils/`.
- Place custom middleware in `middleware/`.
- Place configuration and environment files in `config/`.
- The `main.go` file initializes the server and routes.