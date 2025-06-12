import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen gap-6 responsive-main">
      <h1 className="text-3xl font-bold mb-2">Welcome to Accessibility Analyser</h1>
      <p className="text-lg text-gray-700 mb-4">Analyze your web pages for accessibility and get improvement suggestions.</p>
      <div className="flex gap-4">
        <Link href="/login" className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">Login</Link>
        <Link href="/register" className="bg-gray-200 text-gray-800 px-6 py-2 rounded hover:bg-gray-300">Register</Link>
      </div>
    </div>
  );
}
