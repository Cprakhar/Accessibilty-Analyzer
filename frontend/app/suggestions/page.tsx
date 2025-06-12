"use client";

import { useEffect, useState } from "react";
import { getApiBaseUrl } from "../utils/api";
import { useAuth } from "../utils/AuthContext";
import { useRouter } from "next/navigation";

export default function Suggestions() {
  const { user, loading } = useAuth();
  const [suggestions, setSuggestions] = useState<string[] | null>(null);
  const [error, setError] = useState("");

  const router = useRouter();

  useEffect(() => {
    if (!user && !loading) {
      router.replace("/login");
    }
  }, [user, loading, router]);

  useEffect(() => {
    if (!user) return;
    const reportId = "demo-report-id";
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    fetch(`${getApiBaseUrl()}/reports/${reportId}/suggestions`, {
      headers: {
        'Authorization': token || "",
      },
    })
      .then(async (res) => {
        const data = await res.json();
        if (res.ok && data.success) {
          setSuggestions(data.suggestions || data.data?.suggestions || []);
        } else {
          setError(data.message || "Failed to fetch suggestions");
        }
      })
      .catch(() => setError("Network error"));
  }, [user]);

  if (loading || !user) return null; // Loading or redirecting

  return (
    <div className="min-h-screen flex flex-col items-center justify-center">
      <h2 className="text-2xl font-semibold mb-4">Suggestions</h2>
      <p className="text-gray-700 mb-6">Get AI-powered suggestions to improve your site&#39;s accessibility.</p>
      {error ? (
        <p className="text-red-500">{error}</p>
      ) : suggestions && suggestions.length > 0 ? (
        <ul className="list-disc pl-6 text-left">
          {suggestions.map((s, i) => (
            <li key={i}>{s}</li>
          ))}
        </ul>
      ) : (
        <p className="text-gray-700">No suggestions found.</p>
      )}
    </div>
  );
}
