import { ChangeEvent, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Bullet, Education, Experience, FullResume, Lang, Language, Resume, Software, resume as api } from "../api";
import { useAuth } from "../auth";
import { ErrorPanel, TextField } from "../ui";
import SectionEditor from "../components/SectionEditor";
import AccountSettings from "./AccountSettings";

type Section = "profile" | "software" | "languages" | "education" | "experience" | "account";
const SECTIONS: { key: Section; label: string; icon: string }[] = [
  { key: "profile", label: "Profile", icon: "person" },
  { key: "software", label: "Software", icon: "build" },
  { key: "languages", label: "Languages", icon: "translate" },
  { key: "education", label: "Education", icon: "school" },
  { key: "experience", label: "Experience", icon: "work" },
  { key: "account", label: "Account", icon: "settings" },
];

export default function Account() {
  const [data, setData] = useState<FullResume | null>(null);
  const [section, setSection] = useState<Section>("profile");
  const [selectedExp, setSelectedExp] = useState<number | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    api.mine().then(setData).catch(() => navigate("/login"));
  }, [navigate]);

  if (!data) return <p>Loading…</p>;
  const reload = async () => setData(await api.mine());
  if (selectedExp === null && data.experience.length > 0) setSelectedExp(data.experience[0].id);
  const r = data.resume;

  return (
    <div className="editor">
      <aside className="editor-side">
        <div className="editor-side-name">{data.name}</div>
        {SECTIONS.map((s) => (
          <a key={s.key} href={`#${s.key}`} onClick={(e) => { e.preventDefault(); setSection(s.key); }}
            className={`editor-nav${section === s.key ? " active" : ""}`}>
            <i className="material-icons">{s.icon}</i>{s.label}
          </a>
        ))}
        <div className="editor-side-foot">
          Public page<br />
          {r.published
            ? <Link to={`/${r.slug}`} target="_blank">/{r.slug}</Link>
            : <span className="grey-text">/{r.slug} (not published)</span>}
        </div>
      </aside>

      <section className="editor-main">
        {section === "profile" && <ProfileEditor data={data} onChange={setData} />}
        {section === "account" && <AccountSettings onNameChange={(name) => setData({ ...data, name })} />}

        {section === "software" && (
          <SectionEditor<Software>
            title="Software" items={data.software}
            fields={[{ key: "name", label: "Name", required: true }]}
            summary={(x) => x.name}
            onCreate={async (d) => { await api.createItem("software", d); await reload(); }}
            onUpdate={async (x) => { await api.updateItem("software", x); await reload(); }}
            onDelete={async (id) => { await api.deleteItem("software", id); await reload(); }}
            onReorder={async (ids) => { await api.reorder("software", ids); await reload(); }} />
        )}

        {section === "languages" && (
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
        )}

        {section === "education" && (
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
        )}

        {section === "experience" && (
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

            {selectedExp !== null && (() => {
              const exp = data.experience.find((x) => x.id === selectedExp);
              if (!exp) return null;
              return (
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
              );
            })()}
            {data.experience.length > 1 && (
              <p className="grey-text" style={{ fontSize: ".85rem" }}>Click another experience to edit its responsibilities.</p>
            )}
          </>
        )}
      </section>
    </div>
  );
}

function ProfileEditor({ data, onChange }: { data: FullResume; onChange: (d: FullResume) => void }) {
  const [header, setHeader] = useState<Resume>(data.resume);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);
  const [photoVersion, setPhotoVersion] = useState(0);
  const { logout } = useAuth();
  const navigate = useNavigate();

  const set = (k: keyof Resume) => (v: string) => setHeader({ ...header, [k]: v });

  async function save() {
    setError(null); setSaved(false);
    try {
      const r = await api.save(header);
      setHeader(r);
      onChange({ ...data, resume: r });
      setSaved(true);
    } catch (e) { setError((e as Error).message); }
  }

  async function uploadPhoto(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    try {
      await api.uploadPhoto(file);
      setHeader({ ...header, has_photo: true });
      setPhotoVersion((v) => v + 1);
    } catch (err) { setError((err as Error).message); }
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
        <ErrorPanel message={error} />
        {saved && <div className="card-panel teal lighten-5 teal-text text-darken-3">Saved.</div>}

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
        <button className="btn teal" onClick={save}>Save</button>{" "}
        <button className="btn-flat red-text right" onClick={() => logout().then(() => navigate("/login"))}>Log out</button>
      </div>
    </div>
  );
}
