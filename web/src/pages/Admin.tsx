import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth";
import SideNav, { NavGroup } from "../components/SideNav";
import ResumeEditor, { RESUME_SECTIONS, ResumeSection, isResumeSection } from "../components/ResumeEditor";
import AdminUsers from "./AdminUsers";
import AdminResumes from "./AdminResumes";
import AdminSite from "./AdminSite";
import AdminHome from "./AdminHome";


type AdminSection = "home" | "users" | "resumes" | "site";
type Section = ResumeSection | AdminSection;

const ADMIN_SECTIONS: { key: AdminSection; label: string; icon: string }[] = [
  { key: "home", label: "Home page", icon: "home" },
  { key: "users", label: "Users", icon: "people" },
  { key: "resumes", label: "Resumes", icon: "description" },
  { key: "site", label: "Website", icon: "language" },
];

export default function Admin() {
  const { user, loading, logout } = useAuth();
  const navigate = useNavigate();
  const isOwner = user?.role === "owner";
  const [section, setSection] = useState<Section>(isOwner ? "profile" : "users");
  const groups: NavGroup<Section>[] = isOwner
    ? [{ title: "My resume", items: RESUME_SECTIONS }, { title: "Site", items: ADMIN_SECTIONS }]
    : [{ items: ADMIN_SECTIONS }];
  // Only admins and owners belong here; plain users go to their editor, visitors to login.
  useEffect(() => {
    if (loading) return;
    if (!user || user.role === "user") navigate(user ? "/account" : "/login");
  }, [loading, user, navigate]);

  // Owners start on their resume; the role may arrive after the first render.
  useEffect(() => { if (isOwner) setSection("profile"); }, [isOwner]);

  if (loading || !user || user.role === "user") return <p>Loading…</p>;

  const navItem = (s: { key: Section; label: string; icon: string }) => (
    <a key={s.key} href={`#${s.key}`} onClick={(e) => { e.preventDefault(); setSection(s.key); }}
      className={`editor-nav${section === s.key ? " active" : ""}`}>
      <i className="material-icons">{s.icon}</i>{s.label}
    </a>
  );

  return (
    <div className="editor">
      <SideNav
        header={<>{isOwner ? "Owner" : "Admin"}<br /><small className="grey-text">{user.email}</small></>}
        groups={groups}
        active={section}
        onSelect={setSection}
        footer={<a href="#logout" onClick={(e) => { e.preventDefault(); logout().finally(() => navigate("/login")); }}>Log out</a>}
      />

      <section className="editor-main">
        {isResumeSection(section) ? (
          <ResumeEditor section={section} />
        ) : (
          <>
            {section === "home" && <AdminHome />}
            {section === "users" && <AdminUsers />}
            {section === "resumes" && <AdminResumes />}
            {section === "site" && <AdminSite />}
          </>
        )}
      </section>
    </div>
  );
}
