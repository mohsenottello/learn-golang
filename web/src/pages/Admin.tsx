import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Admin as AdminT, admin } from "../api";
import AdminUsers from "./AdminUsers";
import AdminResumes from "./AdminResumes";
import AdminSite from "./AdminSite";
import AdminHome from "./AdminHome";

type Section = "home" | "users" | "resumes" | "site" | "password";
const SECTIONS: { key: Section; label: string; icon: string }[] = [
  { key: "home", label: "Home page", icon: "home" },
  { key: "users", label: "Users", icon: "people" },
  { key: "resumes", label: "Resumes", icon: "description" },
  { key: "site", label: "Website", icon: "language" },
  { key: "password", label: "Password", icon: "lock" },
];

export default function Admin() {
  const [me, setMe] = useState<AdminT | null>(null);
  const [section, setSection] = useState<Section>("users");
  const navigate = useNavigate();

  useEffect(() => {
    admin.me().then(setMe).catch(() => navigate("/admin/login"));
  }, [navigate]);

  if (!me) return <p>Loading…</p>;

  return (
    <div className="editor">
      <aside className="editor-side">
        <div className="editor-side-name">Admin<br /><small className="grey-text">{me.email}</small></div>
        {SECTIONS.map((s) => (
          <a key={s.key} href={`#${s.key}`} onClick={(e) => { e.preventDefault(); setSection(s.key); }}
            className={`editor-nav${section === s.key ? " active" : ""}`}>
            <i className="material-icons">{s.icon}</i>{s.label}
          </a>
        ))}
        <div className="editor-side-foot">
          <a href="#logout" onClick={(e) => { e.preventDefault(); admin.logout().then(() => navigate("/admin/login")); }}>
            Log out
          </a>
        </div>
      </aside>
      <section className="editor-main">
        {section === "home" && <AdminHome />}
        {section === "site" && <AdminSite />}
        {section === "users" && <AdminUsers />}
        {section === "resumes" && <AdminResumes />}
      </section>
    </div>
  );
}
