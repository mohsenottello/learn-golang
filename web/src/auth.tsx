import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import { User, auth as api } from "./api";

type AuthState = {
  user: User | null;
  loading: boolean;
  setUser: (u: User | null) => void;
  logout: () => Promise<void>;
};

const Ctx = createContext<AuthState>({ user: null, loading: true, setUser: () => {}, logout: async () => {} });

/** Loads the current user once; pages and the nav read it via useAuth(). */
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.me().then(setUser).catch(() => setUser(null)).finally(() => setLoading(false));
  }, []);

  const logout = async () => {
    await api.logout();
    setUser(null);
  };

  return <Ctx.Provider value={{ user, loading, setUser, logout }}>{children}</Ctx.Provider>;
}

export const useAuth = () => useContext(Ctx);
