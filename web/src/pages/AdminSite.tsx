import { ChangeEvent, useState } from "react";
import { Site, admin } from "../api";
import { useSite } from "../site";
import { useToast } from "../toast";
import { TextField } from "../ui";

export default function AdminSite() {
  const { site, refresh, logoVersion } = useSite();
  const [form, setForm] = useState<Site>(site);
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
            <TextField label="Tagline" type="text" value={form.tagline} onChange={set("tagline")} />
            <TextField label="Contact email" type="text" value={form.contact_email} onChange={set("contact_email")} />
            <div className="input-field">
              <textarea className="materialize-textarea" value={form.contact_text} style={{ minHeight: 120 }}
                onChange={(e) => set("contact_text")(e.target.value)} />
              <label className="active">Contact text</label>
            </div>
          </div>
        </div>
        <button className="btn teal" onClick={save}>Save</button>
      </div>
    </div>
  );
}
