import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { FullResume, Lang, resume as api } from "../api";
import { LABELS } from "../labels";

export default function PublicResume() {
  const { slug = "", lang: langParam } = useParams();
  const requested = langParam === "de" || langParam === "en" ? (langParam as Lang) : undefined;
  const [data, setData] = useState<FullResume | null>(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    setData(null);
    api.public(slug, requested).then(setData).catch(() => setNotFound(true));
  }, [slug, requested]);

  useEffect(() => { if (data) document.title = data.name; }, [data]);

  if (notFound) return <div className="container"><h5>Not found</h5><p>No published resume at /{slug}.</p></div>;
  if (!data) return <div className="container"><p>Loading…</p></div>;

  const lang: Lang = data.lang ?? "de";
  const L = LABELS[lang];
  const r = data.resume;
  const fmt = (ym: string) => {
    if (!ym) return "";
    const [y, m] = ym.split("-").map(Number);
    return new Date(y, m - 1, 1).toLocaleDateString(L.locale, { month: "short", year: "numeric" });
  };
  const dash = (from: string, to: string) =>
    !from && !to ? "" : `${fmt(from)}${from ? " – " : ""}${to ? fmt(to) : L.present}`;

  return (
    <>
      <div className="cv-langbar">
        <Link to={`/${slug}/de`} className={lang === "de" ? "active" : ""}>DE</Link>
        <Link to={`/${slug}/en`} className={lang === "en" ? "active" : ""}>EN</Link>
      </div>
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
              <ul className="browser-default cv-contact">
                {r.email && <li><i className="material-icons tiny">mail</i><a href={`mailto:${r.email}`}>{r.email}</a></li>}
                {r.linkedin && <li><i className="material-icons tiny">link</i><a href={r.linkedin} target="_blank" rel="noreferrer">LinkedIn</a></li>}
                {r.phone && <li><i className="material-icons tiny">phone</i>{r.phone}</li>}
                {r.location && <li><i className="material-icons tiny">place</i>{r.location}</li>}
              </ul>
            </>
          )}

          {r.about && <><h6 className="cv-section">{L.about}</h6><p className="cv-about">{r.about}</p></>}

          {data.software.length > 0 && (
            <><h6 className="cv-section">{L.software}</h6>
              <ul className="browser-default cv-list">{data.software.map((s) => <li key={s.id}>{s.name}</li>)}</ul></>
          )}
          {data.languages.length > 0 && (
            <><h6 className="cv-section">{L.languages}</h6>
              <ul className="browser-default cv-list">{data.languages.map((l) => <li key={l.id}>{l.name}{l.level && ` (${l.level})`}</li>)}</ul></>
          )}
          {data.education.length > 0 && (
            <><h6 className="cv-section">{L.education}</h6>
              <ul className="browser-default cv-list">
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
          {data.experience.length > 0 && <h5 className="cv-main-title">{L.experience}</h5>}
          {data.experience.map((x) => (
            <article key={x.id} className="cv-job">
              <div className="cv-job-head">
                <div>
                  <strong>{x.company}</strong>
                  {x.link && <> <a href={x.link} target="_blank" rel="noreferrer">({L.link})</a></>}
                  {x.role && <div className="cv-role">{x.role}</div>}
                  {x.subtitle && <div className="cv-muted"><em>{x.subtitle}</em></div>}
                </div>
                <div className="cv-job-meta">
                  <strong>{dash(x.date_from, x.date_to)}</strong>
                  {x.location && <div className="cv-muted"><em>{x.location}</em></div>}
                </div>
              </div>
              {x.bullets.length > 0 && (
                <ul className="browser-default cv-bullets">{x.bullets.map((b) => <li key={b.id}>{b.text}</li>)}</ul>
              )}
            </article>
          ))}
        </main>
      </div>
    </>
  );
}
