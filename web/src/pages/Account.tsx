import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth";
import SideNav, { NavGroup } from "../components/SideNav";
import ResumeEditor, { RESUME_SECTIONS, ResumeSection } from "../components/ResumeEditor";
import AccountSettings from "./AccountSettings";

type Section = ResumeSection | "account";

const GROUPS: NavGroup<Section>[] = [
  { title: "My resume", items: RESUME_SECTIONS },
  { title: "Settings", items: [{ key: "account", label: "Account", icon: "settings" }] },
];

export default function Account() {
  const { user, loading, logout } = useAuth();
  const [section, setSection] = useState<Section>("profile");
  const navigate = useNavigate();

  // Admins/owners edit from the admin panel; visitors go to login.
  useEffect(() => {
    if (loading) return;
    if (!user) navigate("/login");
    else if (user.role !== "user") navigate("/admin");
  }, [loading, user, navigate]);

  if (loading || !user || user.role !== "user") return <p>Loading…</p>;

  return (
    <div className="editor">
      <SideNav
        header={<>{user.name}<br /><small className="grey-text">{user.email}</small></>}
        groups={GROUPS}
        active={section}
        onSelect={setSection}
        footer={<a href="#logout" onClick={(e) => { e.preventDefault(); logout().finally(() => navigate("/login")); }}>Log out</a>}
      />
      <section className="editor-main">
        {section === "account" ? <AccountSettings /> : <ResumeEditor section={section} />}
      </section>
    </div>
  );
}
