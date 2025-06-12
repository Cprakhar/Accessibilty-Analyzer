// Utility to get the API base URL from env or fallback
export function getApiBaseUrl() {
  return process.env.NEXT_PUBLIC_API_BASE_URL || "/api";
}
