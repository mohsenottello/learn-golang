import { ChangeEvent, CSSProperties, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Bullet, Education, Experience, FullResume, Lang, Language, Resume, Software, resume as api } from "../api";
import { useToast } from "../toast";
import { TextField } from "../ui";
import SectionEditor from "./SectionEditor";
import { StringKey } from "../i18n";

export type ResumeSection = "profile" | "software" | "languages" | "education" | "experience";
export const RESUME_SECTIONS: { key: ResumeSection; label: StringKey; icon: string }[] = [
  { key: "profile", label: "sec_profile", icon: "person" },
  { key: "software", label: "sec_software", icon: "build" },
  { key: "languages", label: "sec_languages", icon: "translate" },
  { key: "education", label: "sec_education", icon: "school" },
  { key: "experience", label: "sec_experience", icon: "work" },
];
export const isResumeSection = (s: string): s is ResumeSection => RESUME_SECTIONS.some((x) => x.key === s);

export default function ResumeEditor({ section }: { section: ResumeSection }) {
  const [data, setData] = useState<FullResume | null>(null);
  const [selectedExp, setSelectedExp] = useState<number | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    api.mine().then(setData).catch(() => navigate("/login"));
  }, [navigate]);

  useEffect(() => {
    if (selectedExp === null && data?.experience.length) setSelectedExp(data.experience[0].id);
  }, [data, selectedExp]);

  if (!data) return <p>Loading…</p>;
  const reload = async () => setData(await api.mine());

  switch (section) {
    case "profile":
      return <ProfileEditor data={data} onChange={setData} />;

    case "software":
      return (
        <SectionEditor<Software>
          title="Software" items={data.software}
          fields={[{ key: "name", label: "Name", required: true }]}
          summary={(x) => x.name}
          onCreate={async (d) => { await api.createItem("software", d); await reload(); }}
          onUpdate={async (x) => { await api.updateItem("software", x); await reload(); }}
          onDelete={async (id) => { await api.deleteItem("software", id); await reload(); }}
          onReorder={async (ids) => { await api.reorder("software", ids); await reload(); }} />
      );

    case "languages":
      return (
        <SectionEditor<Language>
          title="Languages" items={data.languages}
          fields={[
            { key: "name", label: "Language", required: true, translatable: true },
            { key: "level", label: "Level", translatable: true },
          ]}
          summary={(x) => `${x.name}${x.level ? ` (${x.level})` : ""}`}
          onCreate={async (d) => { await api.createItem("languages", d); await reload(); }}
          onUpdate={async (x) => { await api.updateItem("languages", x); await reload(); }}
          onDelete={async (id) => { await api.deleteItem("languages", id); await reload(); }}
          onReorder={async (ids) => { await api.reorder("languages", ids); await reload(); }} />
      );

    case "education":
      return (
        <SectionEditor<Education>
          title="Education" items={data.education}
          fields={[
            { key: "degree", label: "Degree", required: true, translatable: true },
            { key: "school", label: "School", translatable: true },
            { key: "date_from", label: "From", type: "month", half: true },
            { key: "date_to", label: "To", type: "month", half: true },
          ]}
          summary={(x) => `${x.degree} — ${x.school}`}
          onCreate={async (d) => { await api.createItem("education", d); await reload(); }}
          onUpdate={async (x) => { await api.updateItem("education", x); await reload(); }}
          onDelete={async (id) => { await api.deleteItem("education", id); await reload(); }}
          onReorder={async (ids) => { await api.reorder("education", ids); await reload(); }} />
      );

    case "experience": {
      const exp = data.experience.find((x) => x.id === selectedExp);
      return (
        <>
          <SectionEditor<Experience>
            title="Experience" items={data.experience}
            selectedId={selectedExp} onSelect={setSelectedExp}
            fields={[
              { key: "company", label: "Company", required: true, half: true },
              { key: "link", label: "Company website", half: true },
              { key: "role", label: "Role", translatable: true },
              { key: "subtitle", label: "Subtitle", translatable: true },
              { key: "location", label: "Location", translatable: true },
              { key: "date_from", label: "From", type: "month", half: true },
              { key: "date_to", label: "To (empty = present)", type: "month", half: true },
            ]}
            summary={(x) => `${x.company} — ${x.role}  ·  ${x.bullets.length} responsibilit${x.bullets.length === 1 ? "y" : "ies"}`}
            onCreate={async (d) => { await api.createItem("experience", d); await reload(); }}
            onUpdate={async (x) => { await api.updateItem("experience", x); await reload(); }}
            onDelete={async (id) => { await api.deleteItem("experience", id); if (selectedExp === id) setSelectedExp(null); await reload(); }}
            onReorder={async (ids) => { await api.reorder("experience", ids); await reload(); }} />

          {exp && (
            <div className="card">
              <div className="card-content">
                <span className="card-title">Responsibilities — {exp.company}</span>
                <SectionEditor<Bullet>
                  flat title="" items={exp.bullets}
                  fields={[{ key: "text", label: "Text", required: true, translatable: true, type: "textarea" }]}
                  summary={(b) => b.text}
                  onCreate={async (d) => { await api.createBullet(exp.id, d); await reload(); }}
                  onUpdate={async (b) => { await api.updateBullet(exp.id, b); await reload(); }}
                  onDelete={async (id) => { await api.deleteBullet(exp.id, id); await reload(); }}
                  onReorder={async (ids) => { await api.reorderBullets(exp.id, ids); await reload(); }} />
              </div>
            </div>
          )}
          {data.experience.length > 1 && (
            <p className="grey-text" style={{ fontSize: ".85rem" }}>Click another experience to edit its responsibilities.</p>
          )}
        </>
      );
    }
  }
}

/* ---------------- Profile (header, contact, colours) ---------------- */

const DEFAULT_COLORS = { side_bg: "#7f7a72", side_text: "#ffffff", main_bg: "#ffffff", main_text: "#222222", accent: "#555555" } as const;
const COLOR_FIELDS: { key: keyof typeof DEFAULT_COLORS; label: string }[] = [
  { key: "side_bg", label: "Sidebar background" },
  { key: "side_text", label: "Sidebar text" },
  { key: "main_bg", label: "Main background" },
  { key: "main_text", label: "Main text" },
  { key: "accent", label: "Lines & headings" },
];

function ProfileEditor({ data, onChange }: { data: FullResume; onChange: (d: FullResume) => void }) {
  const [header, setHeader] = useState<Resume>({ ...DEFAULT_COLORS, ...data.resume });
  const [photoVersion, setPhotoVersion] = useState(0);
  const toast = useToast();
  const set = (k: keyof Resume) => (v: string) => setHeader({ ...header, [k]: v });

  async function save() {
    try {
      const r = await api.save(header);
      setHeader({ ...DEFAULT_COLORS, ...r });
      onChange({ ...data, resume: r });
      toast("Profile saved.");
    } catch (e) { toast((e as Error).message, "error"); }
  }

  async function uploadPhoto(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    try {
      await api.uploadPhoto(file);
      setHeader({ ...header, has_photo: true });
      setPhotoVersion((v) => v + 1);
      toast("Photo uploaded.");
    } catch (err) { toast((err as Error).message, "error"); }
  }

  const pair = (label: string, key: "title" | "subtitle", hint: string) => (
    <div className="row" style={{ marginBottom: 0 }}>
      <div className="col s12 m6"><TextField label={`${label} (DE)`} type="text" value={header[key]} onChange={set(key)} /></div>
      <div className="col s12 m6"><TextField label={`${label} (EN)`} type="text" value={header[`${key}_en`]} onChange={set(`${key}_en`)} /></div>
      <div className="col s12 grey-text" style={{ fontSize: ".8rem", marginTop: "-.75rem" }}>{hint}</div>
    </div>
  );

  return (
    <div className="card">
      <div className="card-content">
        <span className="card-title">Profile</span>
        <p style={{ marginTop: 0 }}>
          Public page:{" "}
          {header.published
            ? <Link to={`/${header.slug}`} target="_blank">/{header.slug}</Link>
            : <span className="grey-text">/{header.slug} (not published)</span>}
        </p>

        <div className="row" style={{ marginBottom: 0 }}>
          <div className="col s12 m4 center-align">
            {header.has_photo
              ? <img src={`/api/me/resume/photo?v=${photoVersion}`} alt="" className="circle responsive-img" style={{ width: 160, height: 160, objectFit: "cover" }} />
              : <div className="circle grey lighten-2" style={{ width: 160, height: 160, margin: "0 auto" }} />}
            <div className="file-field input-field">
              <div className="btn-flat teal-text"><span>Upload photo</span><input type="file" accept="image/*" onChange={uploadPhoto} /></div>
              <div className="file-path-wrapper" style={{ display: "none" }}><input className="file-path" type="text" /></div>
            </div>
          </div>
          <div className="col s12 m8">
            <TextField label="URL slug" type="text" value={header.slug} onChange={(v) => set("slug")(v.toLowerCase())} />
            <label>Default language of the public page</label>
            <select className="browser-default" value={header.default_lang}
              onChange={(e) => setHeader({ ...header, default_lang: e.target.value as Lang })}>
              <option value="de">Deutsch</option>
              <option value="en">English</option>
            </select>
            <p>
              <label>
                <input type="checkbox" className="filled-in" checked={header.published}
                  onChange={(e) => setHeader({ ...header, published: e.target.checked })} />
                <span>Published (visible to anyone with the link)</span>
              </label>
            </p>
          </div>
        </div>

        {pair("Title", "title", "e.g. Architektin")}
        {pair("Subtitle", "subtitle", "e.g. Bauvorlageberechtigt · Kammer Nr. 157104")}

        <div className="row" style={{ marginBottom: 0 }}>
          <div className="col s12 m6"><TextField label="Email" type="text" value={header.email} onChange={set("email")} /></div>
          <div className="col s12 m6"><TextField label="Phone" type="text" value={header.phone} onChange={set("phone")} /></div>
          <div className="col s12 m6"><TextField label="Location" type="text" value={header.location} onChange={set("location")} /></div>
          <div className="col s12 m6"><TextField label="LinkedIn URL" type="text" value={header.linkedin} onChange={set("linkedin")} /></div>
          <div className="input-field col s12 m6">
            <textarea className="materialize-textarea" value={header.about} style={{ minHeight: 140 }} onChange={(e) => set("about")(e.target.value)} />
            <label className="active">About me (DE)</label>
          </div>
          <div className="input-field col s12 m6">
            <textarea className="materialize-textarea" value={header.about_en} style={{ minHeight: 140 }} onChange={(e) => set("about_en")(e.target.value)} />
            <label className="active">About me (EN)</label>
          </div>
        </div>

        <h6 style={{ marginTop: "1.5rem" }}>Colours</h6>
        <div className="row" style={{ marginBottom: 0 }}>
          <div className="col s12 m7">
            <div className="row" style={{ marginBottom: 0 }}>
              {COLOR_FIELDS.map((c) => (
                <div key={c.key} className="col s6" style={{ marginBottom: 12 }}>
                  <label style={{ display: "block" }}>{c.label}</label>
                  <input type="color" value={header[c.key] ?? DEFAULT_COLORS[c.key]} onChange={(e) => set(c.key)(e.target.value)} />
                </div>
              ))}
            </div>
            <button type="button" className="btn-flat" onClick={() => setHeader({ ...header, ...DEFAULT_COLORS })}>
              <i className="material-icons left">restart_alt</i>Reset to defaults
            </button>
          </div>
          <div className="col s12 m5"><ThemePreview r={header} /></div>
        </div>

        <button className="btn teal" onClick={save}>Save</button>
      </div>
    </div>
  );
}

function ThemePreview({ r }: { r: Resume }) {
  const vars = {
    "--side-bg": r.side_bg, "--side-text": r.side_text,
    "--main-bg": r.main_bg, "--main-text": r.main_text, "--accent": r.accent,
  } as CSSProperties;
  return (
    <div className="theme-preview" style={vars}>
      <div className="tp-side">
        <div className="tp-photo" /><div className="tp-name">Name</div><div className="tp-rule" />
        <div className="tp-line" /><div className="tp-line short" />
      </div>
      <div className="tp-main">
        <div className="tp-title">Berufserfahrung</div>
        <div className="tp-job"><div className="tp-line dark" /><div className="tp-line short" /></div>
        <div className="tp-job"><div className="tp-line dark" /><div className="tp-line short" /></div>
      </div>
    </div>
  );
}
