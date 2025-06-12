"use client";
import { createContext, useContext, useEffect, useState, ReactNode, useMemo, useCallback } from "react";
import { fetchCurrentUser } from "../utils/auth";

export type User = { _id: string; email: string; name?: string } | null;

interface AuthContextType {
  user: User;
  loading: boolean;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  refreshUser: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User>(null);
  const [loading, setLoading] = useState(true);

  const storageKey = "auth_user";
  const getCachedUser = useCallback(() => {
    if (typeof window === "undefined") return null;
    const u = window.localStorage.getItem(storageKey);
    return u ? (JSON.parse(u) as User) : null;
  }, []);
  const setCachedUser = useCallback((u: User) => {
    if (typeof window === "undefined") return;
    if (u) window.localStorage.setItem(storageKey, JSON.stringify(u));
    else window.localStorage.removeItem(storageKey);
  }, []);

  const refreshUser = useCallback(async () => {
    setLoading(true);
    let u = getCachedUser();
    if (!u) {
      u = await fetchCurrentUser();
      setCachedUser(u);
    }
    setUser(u);
    setLoading(false);
  }, [getCachedUser, setCachedUser]);

  useEffect(() => {
    // On mount, try cache first
    const cached = getCachedUser();
    if (cached) {
      setUser(cached);
      setLoading(false);
    } else {
      refreshUser();
    }
    function syncUser() {
      const u = getCachedUser();
      setUser(u);
    }
    window.addEventListener("storage", syncUser);
    return () => window.removeEventListener("storage", syncUser);
  }, [getCachedUser, refreshUser]);

  // When user changes, update cache
  useEffect(() => {
    setCachedUser(user);
  }, [user, setCachedUser]);

  // Memoize context value
  const contextValue = useMemo(() => ({ user, loading, refreshUser }), [user, loading, refreshUser]);

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
