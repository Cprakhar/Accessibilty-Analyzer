"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    const res = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password, name }),
    });
    const data = await res.json();
    if (data.success) {
      localStorage.setItem("token", data.token);
      router.push("/dashboard");
    } else {
      setError(data.message || "Registration failed");
    }
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <form onSubmit={handleSubmit} className="bg-white p-8 rounded shadow w-80 flex flex-col gap-4 responsive-form">
        <h2 className="text-2xl font-bold mb-2">Register</h2>
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={e => setName(e.target.value)}
          className="border p-2 rounded"
          required
        />
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          className="border p-2 rounded"
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          className="border p-2 rounded"
          required
        />
        {error && <div className="text-red-500 text-sm">{error}</div>}
        <button type="submit" className="bg-blue-600 text-white py-2 rounded hover:bg-blue-700">Register</button>
        <a href="/login" className="text-blue-600 text-sm hover:underline text-center">Already have an account? Login</a>
      </form>
    </div>
  );
}
