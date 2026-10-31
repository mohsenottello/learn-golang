import { ChangeEvent, useEffect, useState } from "react";
import { Feature, HomeContent, Lang, admin } from "../api";
import { useSite } from "../site";
import { useToast } from "../toast";
import { TextField } from "../ui";

type F = { key: keyof HomeContent; label: string; multiline?: boolean };
const GROUPS: { title: string; fields: F[] }[] = [
  { title: "Hero", fields: [
    { key: "greeting", label: "Greeting" }, { key: "headline", label: "Headline" },
    { key: "description", label: "Description", multiline: true },
    { key: "primary_label", label: "Primary button" }, { key: "secondary_label", label: "Secondary button" },
    { key: "note", label: "Note under buttons" }, { key: "floating_note", label: "Note on portrait" },
    { key: "featured_slug", label: "Resume slug for 'My Resume'" },
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
  async function copyFromOther() {
    const other = site.home[lang === "de" ? "en" : "de"];
    if (!confirm(`Replace the ${lang.toUpperCase()} form with the ${lang === "de" ? "EN" : "DE"} texts? (not saved until you click Save)`)) return;
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
      <div className="card"><div className="card-content">
        <span className="card-title">Portrait</span>
        {site.has_portrait && <img src={`/api/site/portrait?v=${logoVersion}`} alt="" style={{ maxHeight: 160, borderRadius: 12 }} />}
        <div className="file-field input-field">
          <div className="btn-flat teal-text"><span>Upload portrait</span><input type="file" accept="image/*" onChange={uploadPortrait} /></div>
          <div className="file-path-wrapper" style={{ display: "none" }}><input className="file-path" type="text" /></div>
        </div>
        <p className="grey-text" style={{ fontSize: ".8rem" }}>The portrait is shared by both languages.</p>
      </div></div>

      <div className="admin-lang-tabs">
        {LANGS.map((l) => (
          <button key={l.key} className={`btn-flat${lang === l.key ? " active" : ""}`} onClick={() => setLang(l.key)}>{l.label}</button>
        ))}
        <button className="btn-flat right" onClick={copyFromOther} title="Start from the other language's texts">
          <i className="material-icons left">content_copy</i>Copy from {lang === "de" ? "EN" : "DE"}
        </button>
      </div>

      {GROUPS.map((g) => (
        <div key={g.title} className="card"><div className="card-content">
          <span className="card-title">{g.title}</span>
          {g.fields.map((f) => f.multiline ? (
            <div key={f.key} className="input-field">
              <textarea className="materialize-textarea" value={form[f.key] as string} style={{ minHeight: 100 }}
                onChange={(e) => set(f.key)(e.target.value)} />
              <label className="active">{f.label}</label>
            </div>
          ) : (
            <TextField key={f.key} label={f.label} type="text" value={form[f.key] as string} onChange={set(f.key)} />
          ))}
        </div></div>
      ))}

      <div className="card"><div className="card-content">
        <span className="card-title">Feature cards</span>
        {form.features.map((f, i) => (
          <div key={i} className="row" style={{ marginBottom: 0, borderBottom: "1px solid #eee", paddingBottom: 8 }}>
            <div className="col s6 m2"><TextField label="Icon" type="text" value={f.icon} onChange={(v) => setFeature(i, { icon: v })} /></div>
            <div className="col s6 m2">
              <label>Colour</label>
              <select className="browser-default" value={f.color} onChange={(e) => setFeature(i, { color: e.target.value as Feature["color"] })}>
                {COLORS.map((c) => <option key={c} value={c}>{c}</option>)}
              </select>
            </div>
            <div className="col s12 m3"><TextField label="Title" type="text" value={f.title} onChange={(v) => setFeature(i, { title: v })} /></div>
            <div className="col s10 m4"><TextField label="Text" type="text" value={f.text} onChange={(v) => setFeature(i, { text: v })} /></div>
            <div className="col s2 m1">
              <button className="btn-flat red-text" onClick={() => setForm({ ...form, features: form.features.filter((_, j) => j !== i) })}>
                <i className="material-icons">delete</i>
              </button>
            </div>
          </div>
        ))}
        {form.features.length < 6 && (
          <button className="btn-flat teal-text" onClick={() => setForm({ ...form, features: [...form.features, { icon: "star_outline", color: "teal", title: "", text: "" }] })}>
            <i className="material-icons left">add</i>Add card
          </button>
        )}
        <p className="grey-text" style={{ fontSize: ".8rem" }}>Icon names: <a href="https://fonts.google.com/icons?icon.set=Material+Icons" target="_blank" rel="noreferrer">Material Icons</a></p>
      </div></div>

      <button className="btn teal" onClick={save}>Save {lang.toUpperCase()}</button>{" "}
      <button className="btn-flat" onClick={reset}>Reset {lang.toUpperCase()} to defaults</button>
    </>
  );
}
