import { ChangeEvent, useEffect, useState } from "react";
import { Lang, Site, admin } from "../api";
import { useSite } from "../site";
import { useToast } from "../toast";
import { TextField } from "../ui";

export default function AdminSite() {
  const { site, refresh, logoVersion } = useSite();
  const [form, setForm] = useState<Site>(site);
  useEffect(() => { setForm(site); }, [site]);
  const toast = useToast();
  const set = (k: keyof Site) => (v: string) => setForm({ ...form, [k]: v });

  async function save() {
    try {
      await admin.updateSite(form);
      await refresh();
      toast("Website settings saved.");
    } catch (e) { toast((e as Error).message, "error"); }
  }

  async function uploadLogo(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    try {
      await admin.uploadLogo(file);
      await refresh();
      toast("Logo updated.");
    } catch (err) { toast((err as Error).message, "error"); }
  }

  return (
    <div className="card">
      <div className="card-content">
        <span className="card-title">Website</span>
        <div className="row" style={{ marginBottom: 0 }}>
          <div className="col s12 m4 center-align">
            {site.has_logo
              ? <img src={`/api/site/logo?v=${logoVersion}`} alt="" style={{ maxHeight: 96, maxWidth: "100%" }} />
              : <div className="grey lighten-3" style={{ height: 96, lineHeight: "96px" }}>No logo</div>}
            <div className="file-field input-field">
              <div className="btn-flat teal-text"><span>Upload logo</span><input type="file" accept="image/*" onChange={uploadLogo} /></div>
              <div className="file-path-wrapper" style={{ display: "none" }}><input className="file-path" type="text" /></div>
            </div>
          </div>
          <div className="col s12 m8">
            <TextField label="Site name" type="text" value={form.name} onChange={set("name")} />
            <div className="row" style={{ marginBottom: 0 }}>
              <div className="col s12 m6"><TextField label="Tagline (DE)" type="text" value={form.tagline} onChange={set("tagline")} /></div>
              <div className="col s12 m6"><TextField label="Tagline (EN)" type="text" value={form.tagline_en} onChange={set("tagline_en")} /></div>
            </div>
            <TextField label="Contact email" type="text" value={form.contact_email} onChange={set("contact_email")} />
            <div className="row" style={{ marginBottom: 0 }}>
              <div className="input-field col s12 m6">
                <textarea className="materialize-textarea" value={form.contact_text} style={{ minHeight: 120 }}
                  onChange={(e) => set("contact_text")(e.target.value)} />
                <label className="active">Contact text (DE)</label>
              </div>
              <div className="input-field col s12 m6">
                <textarea className="materialize-textarea" value={form.contact_text_en} style={{ minHeight: 120 }}
                  onChange={(e) => set("contact_text_en")(e.target.value)} />
                <label className="active">Contact text (EN)</label>
              </div>
            </div>

            <label>Default language (for visitors whose browser is neither DE nor a saved choice)</label>
            <select className="browser-default" value={form.default_lang}
              onChange={(e) => setForm({ ...form, default_lang: e.target.value as Lang })}>
              <option value="de">Deutsch</option>
              <option value="en">English</option>
            </select>

            <p>
              <label>
                <input type="checkbox" className="filled-in" checked={form.registration}
                  onChange={(e) => setForm({ ...form, registration: e.target.checked })} />
                <span>Allow user registration</span>
              </label>
            </p>
          </div>
        </div>
        <button className="btn teal" onClick={save}>Save</button>
      </div>
    </div>
  );
}
