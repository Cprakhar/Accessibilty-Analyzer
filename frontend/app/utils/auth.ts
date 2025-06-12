// Utility to get the current user from the backend using the JWT token
export async function fetchCurrentUser() {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
  if (!token) return null;
  const res = await fetch("/api/auth/me", {
    headers: { "Authorization": token },
  });
  if (!res.ok) return null;
  const data = await res.json();
  if (data && data._id) return data;
  return null;
}
