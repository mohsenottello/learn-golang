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
    <div className="card mb-3">
      <div className="card-body">
        <h5 className="card-title">Website</h5>
        <div className="row">
          <div className="col-12 col-md-4 text-center">
            {site.has_logo
              ? <img src={`/api/site/logo?v=${logoVersion}`} alt={site.name} style={{ maxHeight: 96, maxWidth: "100%" }} />
              : <div className="bg-secondary-subtle rounded" style={{ height: 96, lineHeight: "96px" }}>No logo</div>}
            <div className="mt-2"><label className="form-label small text-secondary">Upload logo</label><input type="file" className="form-control form-control-sm" accept="image/*" onChange={uploadLogo} /></div>
          </div>
          <div className="col-12 col-md-8">
            <TextField label="Site name" type="text" value={form.name} onChange={set("name")} />
            <div className="row">
              <div className="col-12 col-md-6"><TextField label="Tagline (DE)" type="text" value={form.tagline} onChange={set("tagline")} /></div>
              <div className="col-12 col-md-6"><TextField label="Tagline (EN)" type="text" value={form.tagline_en} onChange={set("tagline_en")} /></div>
            </div>
            <TextField label="Contact email" type="text" value={form.contact_email} onChange={set("contact_email")} />
            <div className="row">
              <div className="col-12 col-md-6">
                <div className="form-floating mb-3">
                  <textarea className="form-control" value={form.contact_text} style={{ height: 130 }} placeholder="Contact text (DE)"
                    onChange={(e) => set("contact_text")(e.target.value)} />
                  <label>Contact text (DE)</label>
                </div>
              </div>
              <div className="col-12 col-md-6">
                <div className="form-floating mb-3">
                  <textarea className="form-control" value={form.contact_text_en} style={{ height: 130 }} placeholder="Contact text (EN)"
                    onChange={(e) => set("contact_text_en")(e.target.value)} />
                  <label>Contact text (EN)</label>
                </div>
              </div>
            </div>

            <label className="form-label">Default language (for visitors whose browser is neither DE nor a saved choice)</label>
            <select className="form-select mb-3" value={form.default_lang}
              onChange={(e) => setForm({ ...form, default_lang: e.target.value as Lang })}>
              <option value="de">Deutsch</option>
              <option value="en">English</option>
            </select>

            <div className="form-check"><input type="checkbox" className="form-check-input" id="registration" checked={form.registration}
                  onChange={(e) => setForm({ ...form, registration: e.target.checked })} /><label className="form-check-label" htmlFor="registration">Allow user registration</label></div>
          </div>
        </div>
        <button className="btn btn-primary" onClick={save}>Save</button>
      </div>
    </div>
  );
}
