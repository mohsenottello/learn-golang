import { FormEvent, useEffect, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import { FullResume, Lang, resume as api } from "../api";
import { STRINGS } from "../i18n";
import { useLang } from "../lang";
import { useSite } from "../site";
import { useToast } from "../toast";

export default function PublicResume({ featured = false }: { featured?: boolean }) {
  const params = useParams();
  const [search] = useSearchParams();
  const print = search.has("print");
  const { site } = useSite();
  const { lang: siteLang, setLang: setSiteLang } = useLang();

  const slug = featured ? site.owner_slug : (params.slug ?? "");
  const basePath = featured ? "/resume" : `/${slug}`;

  const langParam = params.lang;
  const requested: Lang = langParam === "de" || langParam === "en" ? (langParam as Lang) : siteLang;
  const [data, setData] = useState<FullResume | null>(null);
  const [notFound, setNotFound] = useState(false);
    const [pw, setPw] = useState("");
  const toast = useToast();

  async function unlock(e: FormEvent) {
    e.preventDefault();
    try {
      await api.unlock(slug, pw);
      setPw("");
      setData(await api.public(slug, requested));
    } catch (err) { toast((err as Error).message, "error"); }
  }
  useEffect(() => {
    if (!slug) return;
    setData(null);
    setNotFound(false);
    api.public(slug, requested).then(setData).catch(() => setNotFound(true));
  }, [slug, requested]);

  // Tells the PDF renderer the page is fully drawn (it waits for this attribute).
  useEffect(() => {
    if (data) document.body.dataset.ready = "1";
    return () => { delete document.body.dataset.ready; };
  }, [data]);

  if (notFound) return <div className="container"><h5>404</h5><p>{STRINGS[siteLang].cv_notfound} {basePath}.</p></div>;
  if (!data) return <div className="container"><p>Loading…</p></div>;

  const lang: Lang = data.lang ?? requested;
  const S = STRINGS[lang];
  const L = { contact: S.cv_contact, about: S.cv_about, software: S.cv_software, languages: S.cv_languages,
    education: S.cv_education, experience: S.cv_experience, present: S.cv_present, link: S.cv_link, locale: S.locale };
  const r = data.resume;
  const fmt = (ym: string) => {
    if (!ym) return "";
    const [y, m] = ym.split("-").map(Number);
    return new Date(y, m - 1, 1).toLocaleDateString(L.locale, { month: "short", year: "numeric" });
  };
  const href = (u: string) => (/^https?:\/\//i.test(u) ? u : `https://${u}`);
  const dash = (from: string, to: string) =>
    !from && !to ? "" : `${fmt(from)}${from ? " – " : ""}${to ? fmt(to) : L.present}`;

    return (
    <>
      {!print && (
        <div className="cv-langbar">
          {!featured && (
            <>
              <Link to={`${basePath}/de`} className={lang === "de" ? "active" : ""} onClick={() => setSiteLang("de")}>DE</Link>
              <Link to={`${basePath}/en`} className={lang === "en" ? "active" : ""} onClick={() => setSiteLang("en")}>EN</Link>
            </>
          )}
          {data.hidden_experience === 0 && (
            <a className="btn btn-link btn-sm" href={`/api/resumes/${slug}/pdf?lang=${lang}`}>
              <i className="material-icons align-middle me-1">picture_as_pdf</i>{S.cv_pdf}
            </a>
          )}
        </div>
      )}
      <div className="cv" lang={lang} style={{
          "--side-bg": r.side_bg, "--side-text": r.side_text,
          "--main-bg": r.main_bg, "--main-text": r.main_text, "--accent": r.accent,
        } as React.CSSProperties}>
        <aside className="cv-side">
          {r.has_photo && <img className="cv-photo" src={`/api/resumes/${slug}/photo`} alt={data.name} />}
          <h4 className="cv-name">{data.name}</h4>
          {r.title && <div className="cv-title">{r.title}</div>}
          {r.subtitle && <div className="cv-subtitle">{r.subtitle}</div>}

          {(r.email || r.phone || r.location || r.linkedin) && (
            <>
              <h6 className="cv-section">{L.contact}</h6>
              <ul className="cv-contact">
                {r.email && <li><i className="material-icons cv-icon">mail</i><a href={`mailto:${r.email}`}>{r.email}</a></li>}
                {r.linkedin && <li><i className="material-icons cv-icon">link</i><a href={r.linkedin} target="_blank" rel="noreferrer">LinkedIn</a></li>}
                {r.phone && <li><i className="material-icons cv-icon">phone</i>{r.phone}</li>}
                {r.location && <li><i className="material-icons cv-icon">place</i>{r.location}</li>}
              </ul>
            </>
          )}

          {r.about && <><h6 className="cv-section">{L.about}</h6><p className="cv-about">{r.about}</p></>}

          {data.software.length > 0 && (
            <><h6 className="cv-section">{L.software}</h6>
              <ul className="cv-list">{data.software.map((s) => <li key={s.id}>{s.name}</li>)}</ul></>
          )}
          {data.languages.length > 0 && (
            <><h6 className="cv-section">{L.languages}</h6>
              <ul className="cv-list">{data.languages.map((l) => <li key={l.id}>{l.name}{l.level && ` (${l.level})`}</li>)}</ul></>
          )}
          {data.education.length > 0 && (
            <><h6 className="cv-section">{L.education}</h6>
              <ul className="cv-list">
                {data.education.map((e) => (
                  <li key={e.id}>
                    <strong>{e.degree}</strong>
                    {e.school && <div className="cv-muted">{e.school}</div>}
                    {(e.date_from || e.date_to) && <div className="cv-muted">{dash(e.date_from, e.date_to)}</div>}
                  </li>
                ))}
              </ul></>
          )}
        </aside>

        <main className="cv-main">
          {(data.experience.length > 0 || data.hidden_experience > 0) && <h5 className="cv-main-title">{L.experience}</h5>}
          {data.experience.map((x) => (
            <article key={x.id} className="cv-job">
              <div className="cv-job-head">
                <div>
                  {x.link
                    ? <a className="cv-company" href={href(x.link)} target="_blank" rel="noreferrer"><strong>{x.company}</strong></a>
                    : <strong>{x.company}</strong>}
                  {x.role && <div className="cv-role">{x.role}</div>}
                  {x.subtitle && <div className="cv-muted"><em>{x.subtitle}</em></div>}
                </div>
                <div className="cv-job-meta">
                  <strong>{dash(x.date_from, x.date_to)}</strong>
                  {x.location && <div className="cv-muted"><em>{x.location}</em></div>}
                </div>
              </div>
              {x.bullets.length > 0 && (
                <ul className="cv-bullets">{x.bullets.map((b) => <li key={b.id}>{b.text}</li>)}</ul>
              )}
            </article>
          ))}

          {data.hidden_experience > 0 && (
            <div className="cv-locked">
              <div className="cv-locked-blur">
                {Array.from({ length: Math.min(data.hidden_experience, 3) }).map((_, i) => (
                  <article key={i} className="cv-job" aria-hidden>
                    <div className="cv-job-head"><div><strong>████████ ████</strong><div className="cv-role">█████████</div></div><div className="cv-job-meta"><strong>████ – ████</strong></div></div>
                    <ul className="cv-bullets"><li>██████████ ████████ ███████████ ███</li><li>████████ ██████ ████████████</li></ul>
                  </article>
                ))}
              </div>
              <form className="cv-unlock" onSubmit={unlock}>
                <i className="material-icons">lock</i>
                <p>{S.cv_locked.replace("{n}", String(data.hidden_experience))}</p>
                <div className="input-group">
                  <input type="password" className="form-control" placeholder={S.cv_password} value={pw} onChange={(e) => setPw(e.target.value)} />
                  <button className="btn btn-primary" disabled={!pw}>{S.cv_unlock}</button>
                </div>
              </form>
            </div>
          )}
        </main>
      </div>
    </>
  );
}
