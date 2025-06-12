'use client';

// Dashboard page for authenticated users
// TODO: Implement data fetching and integration with backend

import React, { useEffect } from "react";
import { useAuth } from "../utils/AuthContext";
import { useRouter } from "next/navigation";

export default function Dashboard() {
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
    <div className="min-h-screen flex flex-col">
      {/* Navbar is now rendered globally in layout.tsx */}
      <main className="flex-1 flex flex-col items-center justify-center p-8 responsive-main">
        <h2 className="text-2xl font-semibold mb-4">Dashboard</h2>
        <p className="text-gray-700 mb-6">Welcome, {user.name || user.email}! Here you can view your accessibility reports, run new analyses, and get improvement suggestions.</p>
        {/* TODO: Add dashboard widgets and analytics here */}
      </main>
    </div>
  );
}
