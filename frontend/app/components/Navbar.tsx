// Simple navigation bar for the app
"use client";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "../utils/AuthContext";
import { useRouter } from "next/navigation";

function getInitials(name?: string, email?: string) {
  if (!name && email) return email[0]?.toUpperCase() || "U";
  if (!name) return "U";
  const parts = name.trim().split(" ");
  if (parts.length === 1) return parts[0][0]?.toUpperCase() || "U";
  return (
    (parts[0][0] || "").toUpperCase() +
    (parts[parts.length - 1][0] || "").toUpperCase()
  );
}

export default function Navbar() {
  const { user } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const router = useRouter();
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setDropdownOpen(false);
      }
    }
    if (dropdownOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [dropdownOpen]);

  function handleLogout() {
    localStorage.removeItem("token");
    window.dispatchEvent(new Event("storage")); // trigger sync in all tabs
    router.push("/login");
  }

  return (
    <nav className="bg-white border-b shadow px-6 py-4 flex items-center justify-between responsive-nav">
      <div className="flex items-center gap-2">
        <span className="font-bold text-xl text-blue-700">Accessibility Analyser</span>
      </div>
      <div className="flex gap-4 responsive-nav-links items-center">
        <Link href="/dashboard" className="text-blue-700 font-medium hover:underline">Dashboard</Link>
        <Link href="/reports" className="text-gray-700 hover:underline">Reports</Link>
        <Link href="/suggestions" className="text-gray-700 hover:underline">Suggestions</Link>
        {user ? (
          <div className="relative" ref={dropdownRef}>
            <button
              className="ml-4 w-10 h-10 bg-blue-600 text-white rounded-full flex items-center justify-center text-lg font-bold focus:outline-none"
              onClick={e => {
                e.preventDefault();
                e.stopPropagation();
                setDropdownOpen(true);
              }}
              aria-label="User menu"
              tabIndex={0}
            >
              {getInitials(user.name, user.email)}
            </button>
            {dropdownOpen && (
              <div className="absolute right-0 mt-2 w-40 bg-white border rounded shadow z-10">
                <button
                  className="block w-full text-left px-4 py-2 hover:bg-gray-100"
                  onClick={() => { setDropdownOpen(false); router.push("/settings"); }}
                >
                  Settings
                </button>
                <button
                  className="block w-full text-left px-4 py-2 hover:bg-gray-100 text-red-600"
                  onClick={() => { setDropdownOpen(false); handleLogout(); }}
                >
                  Logout
                </button>
              </div>
            )}
          </div>
        ) : (
          <Link href="/login" className="ml-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Login</Link>
        )}
      </div>
    </nav>
  );
}
