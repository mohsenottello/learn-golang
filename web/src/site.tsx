import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from "react";
import { HomeContent, Site, site as api } from "./api";

export const EMPTY_HOME: HomeContent = {
  greeting: "", headline: "", description: "", primary_label: "", secondary_label: "",
  note: "", floating_note: "",
  features_label: "", features_heading: "", features_intro: "", features: [],
  about_label: "", about_heading: "", about_text: "", about_button: "",
  cta_heading: "", cta_text: "", cta_button: "", footer_text: "",
};

const FALLBACK: Site = {
  name: "", tagline: "", tagline_en: "", contact_email: "", contact_text: "", contact_text_en: "",
  primary_color: "#4f46e5", registration: false, default_lang: "de", owner_slug: "",
  has_logo: false, has_portrait: false, home: { de: EMPTY_HOME, en: EMPTY_HOME }, translation_enabled: false,
};

type SiteState = { site: Site; loaded: boolean; refresh: () => Promise<void>; logoVersion: number };
const Ctx = createContext<SiteState>({ site: FALLBACK, loaded: false, refresh: async () => {}, logoVersion: 0 });

function darken(hex: string, f: number) {
  const n = parseInt(hex.slice(1), 16);
  const c = (v: number) => Math.round(v * f).toString(16).padStart(2, "0");
  return `#${c(n >> 16)}${c((n >> 8) & 255)}${c(n & 255)}`;
}

export function SiteProvider({ children }: { children: ReactNode }) {
  const [site, setSite] = useState<Site>(FALLBACK);
  const [loaded, setLoaded] = useState(false);
  const [logoVersion, setLogoVersion] = useState(0);

  const refresh = useCallback(async () => {
    setSite(await api.get());
    setLoaded(true);
    setLogoVersion((v) => v + 1); // busts the browser cache of /api/site/logo and /portrait
  }, []);

  useEffect(() => { refresh().catch(() => {}); }, [refresh]);

  // Theme colours only. The <title> and favicon come from the server-rendered HTML.
  useEffect(() => {
    if (!loaded) return;
    document.documentElement.style.setProperty("--primary", site.primary_color);
    document.documentElement.style.setProperty("--primary-dark", darken(site.primary_color, 0.75));
  }, [loaded, site.primary_color]);

  return <Ctx.Provider value={{ site, loaded, refresh, logoVersion }}>{children}</Ctx.Provider>;
}

export const useSite = () => useContext(Ctx);
