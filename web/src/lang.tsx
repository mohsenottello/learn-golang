import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { Lang } from "./api";
import { STRINGS, StringKey } from "./i18n";
import { useSite } from "./site";

type LangState = { lang: Lang; setLang: (l: Lang) => void; t: (k: StringKey) => string };
const Ctx = createContext<LangState>({ lang: "de", setLang: () => {}, t: (k) => STRINGS.de[k] });

const KEY = "lang";
const isLang = (s: string | null): s is Lang => s === "de" || s === "en";

/** Must sit inside SiteProvider: the site's default language is the last fallback. */
export function LangProvider({ children }: { children: ReactNode }) {
  const { site } = useSite();
  const [chosen, setChosen] = useState<Lang | null>(() => {
    const saved = localStorage.getItem(KEY);
    return isLang(saved) ? saved : null;
  });

  // Priority: what the visitor chose → browser language → site default.
  const lang: Lang = chosen ?? (navigator.language.toLowerCase().startsWith("de") ? "de" : site.default_lang);

  const setLang = useCallback((l: Lang) => {
    setChosen(l);
    localStorage.setItem(KEY, l);
  }, []);

  useEffect(() => { document.documentElement.lang = lang; }, [lang]);

  const value = useMemo<LangState>(() => ({ lang, setLang, t: (k) => STRINGS[lang][k] }), [lang, setLang]);
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export const useLang = () => useContext(Ctx);
