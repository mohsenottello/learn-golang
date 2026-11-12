import { ChangeEvent, useEffect, useState } from "react";
import { Feature, HomeContent, Lang, admin } from "../api";
import { useSite } from "../site";
import { useToast } from "../toast";
import { TextField } from "../ui";
import TranslateButton from "../components/TranslateButton";

type F = { key: keyof HomeContent; label: string; multiline?: boolean };
const GROUPS: { title: string; fields: F[] }[] = [
  { title: "Hero", fields: [
    { key: "greeting", label: "Greeting" }, { key: "headline", label: "Headline" },
    { key: "description", label: "Description", multiline: true },
    { key: "primary_label", label: "Primary button" }, { key: "secondary_label", label: "Secondary button" },
    { key: "note", label: "Note under buttons" }, { key: "floating_note", label: "Note on portrait" },
  ]},
  { title: "Features", fields: [
    { key: "features_label", label: "Label" }, { key: "features_heading", label: "Heading" },
    { key: "features_intro", label: "Intro", multiline: true },
  ]},
  { title: "About", fields: [
    { key: "about_label", label: "Label" }, { key: "about_heading", label: "Heading" },
    { key: "about_text", label: "Text", multiline: true }, { key: "about_button", label: "Button" },
  ]},
  { title: "Call to action", fields: [
    { key: "cta_heading", label: "Heading" }, { key: "cta_text", label: "Text", multiline: true }, { key: "cta_button", label: "Button" },
  ]},
  { title: "Footer", fields: [{ key: "footer_text", label: "Footer text", multiline: true }] },
];
const COLORS: Feature["color"][] = ["teal", "purple", "pink", "blue"];
const LANGS: { key: Lang; label: string }[] = [{ key: "de", label: "Deutsch" }, { key: "en", label: "English" }];

export default function AdminHome() {
  const { site, refresh, logoVersion } = useSite();
  const [lang, setLang] = useState<Lang>(site.default_lang);
  const [form, setForm] = useState<HomeContent>(site.home[lang]);
  const toast = useToast();
  const TEXT_KEYS = GROUPS.flatMap((g) => g.fields.map((f) => f.key));
  const other = site.home[lang === "de" ? "en" : "de"];
  const otherLang: Lang = lang === "de" ? "en" : "de";

  // Re-sync whenever the server copy or the selected language changes.
  useEffect(() => { setForm(site.home[lang]); }, [site.home, lang]);

  const set = (k: keyof HomeContent) => (v: string) => setForm({ ...form, [k]: v });
  const setFeature = (i: number, patch: Partial<Feature>) =>
    setForm({ ...form, features: form.features.map((f, j) => (j === i ? { ...f, ...patch } : f)) });

  async function save() {
    try { await admin.updateHome(lang, form); await refresh(); toast(`Home page (${lang.toUpperCase()}) saved.`); }
    catch (e) { toast((e as Error).message, "error"); }
  }
  async function reset() {
    if (!confirm(`Reset the ${lang.toUpperCase()} home page texts to the defaults?`)) return;
    try { await admin.resetHome(lang); await refresh(); toast("Reset to defaults."); }
    catch (e) { toast((e as Error).message, "error"); }
  }
  function copyFromOther() {
    if (!confirm(`Replace the ${lang.toUpperCase()} form with the ${otherLang.toUpperCase()} texts? (not saved until you click Save)`)) return;
    setForm(other);
  }
  async function uploadPortrait(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    try { await admin.uploadPortrait(file); await refresh(); toast("Portrait updated."); }
    catch (err) { toast((err as Error).message, "error"); }
  }

  return (
    <>
      <div className="card mb-3"><div className="card-body">
        <h5 className="card-title">Portrait</h5>
        {site.has_portrait && <img src={`/api/site/portrait?v=${logoVersion}`} alt={site.name} style={{ maxHeight: 160, borderRadius: 12 }} />}
        <div className="mt-2"><label className="form-label small text-secondary">Upload portrait</label><input type="file" className="form-control form-control-sm" accept="image/*" onChange={uploadPortrait} /></div>
        <p className="text-secondary" style={{ fontSize: ".8rem" }}>The portrait is shared by both languages.</p>
      </div></div>

      <div className="admin-lang-tabs">
        <ul className="nav nav-tabs border-0 mb-0">
          {LANGS.map((l) => (
            <li key={l.key} className="nav-item">
              <button className={`nav-link${lang === l.key ? " active" : ""}`} onClick={() => setLang(l.key)}>{l.label}</button>
            </li>
          ))}
        </ul>
        <span className="ms-auto d-flex gap-1">
          <button className="btn btn-link" onClick={copyFromOther} title="Start from the other language's texts">
            <i className="material-icons align-middle me-1">content_copy</i>Copy from {otherLang.toUpperCase()}
          </button>
          <TranslateButton from={otherLang} to={lang}
            texts={[
              ...TEXT_KEYS.map((k) => other[k] as string),
              ...other.features.flatMap((f) => [f.title, f.text]),
            ]}
            onResult={(out) => {
              const next: HomeContent = { ...form };
              TEXT_KEYS.forEach((k, i) => { (next as any)[k] = out[i]; });
              const base = TEXT_KEYS.length;
              next.features = other.features.map((f, i) => ({ ...f, title: out[base + i * 2], text: out[base + i * 2 + 1] }));
              setForm(next);
            }} />
        </span>
      </div>

      {GROUPS.map((g) => (
        <div key={g.title} className="card mb-3"><div className="card-body">
          <h5 className="card-title">{g.title}</h5>
          {g.fields.map((f) => f.multiline ? (
            <div key={f.key} className="form-floating mb-3">
              <textarea className="form-control" value={form[f.key] as string} style={{ height: 110 }} placeholder={f.label}
                onChange={(e) => set(f.key)(e.target.value)} />
              <label>{f.label}</label>
            </div>
          ) : (
            <TextField key={f.key} label={f.label} type="text" value={form[f.key] as string} onChange={set(f.key)} />
          ))}
        </div></div>
      ))}

      <div className="card mb-3"><div className="card-body">
        <h5 className="card-title">Feature cards</h5>
        {form.features.map((f, i) => (
          <div key={i} className="row" style={{ marginBottom: 0, borderBottom: "1px solid #eee", paddingBottom: 8 }}>
            <div className="col-6 col-md-2"><TextField label="Icon" type="text" value={f.icon} onChange={(v) => setFeature(i, { icon: v })} /></div>
            <div className="col-6 col-md-2">
              <label className="form-label small">Colour</label>
              <select className="form-select" value={f.color} onChange={(e) => setFeature(i, { color: e.target.value as Feature["color"] })}>
                {COLORS.map((c) => <option key={c} value={c}>{c}</option>)}
              </select>
            </div>
            <div className="col-12 col-md-3"><TextField label="Title" type="text" value={f.title} onChange={(v) => setFeature(i, { title: v })} /></div>
            <div className="col-10 col-md-4"><TextField label="Text" type="text" value={f.text} onChange={(v) => setFeature(i, { text: v })} /></div>
            <div className="col-2 col-md-1">
              <button className="btn btn-link text-danger" onClick={() => setForm({ ...form, features: form.features.filter((_, j) => j !== i) })}>
                <i className="material-icons">delete</i>
              </button>
            </div>
          </div>
        ))}
        {form.features.length < 6 && (
          <button className="btn btn-link" onClick={() => setForm({ ...form, features: [...form.features, { icon: "star_outline", color: "teal", title: "", text: "" }] })}>
            <i className="material-icons align-middle me-1">add</i>Add card
          </button>
        )}
        <p className="text-secondary" style={{ fontSize: ".8rem" }}>Icon names: <a href="https://fonts.google.com/icons?icon.set=Material+Icons" target="_blank" rel="noreferrer">Material Icons</a></p>
      </div></div>

      <button className="btn btn-primary" onClick={save}>Save {lang.toUpperCase()}</button>{" "}
      <button className="btn btn-link" onClick={reset}>Reset {lang.toUpperCase()} to defaults</button>
    </>
  );
}
