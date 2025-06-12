"use client";
import React, { useEffect } from "react";
import { useAuth } from "../utils/AuthContext";
import { useRouter } from "next/navigation";

export default function Reports() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/login");
    }
  }, [loading, user, router]);

  if (loading) return <div className="flex-1 flex items-center justify-center">Loading...</div>;
  if (!user) return null; // Redirecting

  return (
    <div className="min-h-screen flex flex-col items-center justify-center">
      <h2 className="text-2xl font-semibold mb-4">Reports</h2>
      <p className="text-gray-700 mb-6">View your past accessibility analysis reports here.</p>
      {/* TODO: Implement reports listing and details */}
    </div>
  );
}
