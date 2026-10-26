import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from "react";
import { Site, site as api } from "./api";

const FALLBACK: Site = { name: "MyApp", tagline: "", contact_email: "", contact_text: "", has_logo: false };

type SiteState = { site: Site; refresh: () => Promise<void>; logoVersion: number };
const Ctx = createContext<SiteState>({ site: FALLBACK, refresh: async () => {}, logoVersion: 0 });

export function SiteProvider({ children }: { children: ReactNode }) {
  const [site, setSite] = useState<Site>(FALLBACK);
  const [logoVersion, setLogoVersion] = useState(0);

  const refresh = useCallback(async () => {
    setSite(await api.get());
    setLogoVersion((v) => v + 1); // busts the browser cache of /api/site/logo
  }, []);

  useEffect(() => { refresh().catch(() => {}); }, [refresh]);
  useEffect(() => { document.title = site.name; }, [site.name]);
  useEffect(() => {
  document.title = site.name;
  const link = document.getElementById("favicon") as HTMLLinkElement | null;
  if (link) link.href = site.has_logo ? `/api/site/logo?v=${logoVersion}` : "/favicon.svg";
}, [site.name, site.has_logo, logoVersion]);

  return <Ctx.Provider value={{ site, refresh, logoVersion }}>{children}</Ctx.Provider>;
}

export const useSite = () => useContext(Ctx);
