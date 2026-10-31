import { ReactNode, useEffect } from "react";
import { Link, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "./auth";
import { useSite } from "./site";
import { useLang } from "./lang";
import { Lang } from "./api";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Verify from "./pages/Verify";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import Account from "./pages/Account";
import PublicResume from "./pages/PublicResume";
import Admin from "./pages/Admin";
import Contact from "./pages/Contact";

const appPaths = ["/", "/register", "/login", "/verify", "/forgot-password", "/reset-password", "/account", "/contact", "/admin"];

type NavLink = { to: string; label: string; onClick?: () => void; cta?: boolean };

export default function App() {
  const { pathname } = useLocation();
  const { user, logout } = useAuth();
  const { site, logoVersion } = useSite();
  const { lang, setLang, t } = useLang();
  const navigate = useNavigate();
  const isAdminArea = pathname.startsWith("/admin");
  const resumeHref = `/${site.home[lang].featured_slug || site.home.de.featured_slug}`;

  useEffect(() => { M.Sidenav.init(document.querySelectorAll(".sidenav")); }, []);

  const userLogout = () => logout().finally(() => navigate("/login"));
  const isStaff = user?.role === "admin" || user?.role === "owner";

  const links: NavLink[] = [];
  if (isAdminArea) {
    if (pathname === "/admin") links.push({ to: "#", label: t("nav_logout"), onClick: userLogout });
  } else {
    links.push({ to: "/", label: t("nav_home") }, { to: resumeHref, label: t("nav_resume") });
    if (user) links.push({ to: isStaff ? "/admin" : "/account", label: isStaff ? t("nav_admin") : t("nav_editor") });
    links.push({ to: "/contact", label: t("nav_contact") });
    if (user) {
      links.push({ to: "#", label: t("nav_logout"), onClick: userLogout });
    } else {
      links.push({ to: "/login", label: t("nav_login") });
      if (site.registration) links.push({ to: "/register", label: t("nav_register"), cta: true });
    }
  }

  if (!appPaths.includes(pathname)) {
    return (
      <Routes>
        <Route path="/:slug/:lang" element={<PublicResume />} />
        <Route path="/:slug" element={<PublicResume />} />
      </Routes>
    );
  }

  const renderLink = (l: NavLink) => (
    <li key={l.label}>
      {l.onClick
        ? <a href={l.to} onClick={(e) => { e.preventDefault(); l.onClick!(); }}>{l.label}</a>
        : <Link to={l.to} className={l.cta ? "nav-cta" : undefined}>{l.label}</Link>}
    </li>
  );

  const langToggle = (
    <li key="lang" className="lang-toggle">
      {(["de", "en"] as Lang[]).map((l) => (
        <a key={l} href="#" className={lang === l ? "active" : ""} onClick={(e) => { e.preventDefault(); setLang(l); }}>
          {l.toUpperCase()}
        </a>
      ))}
    </li>
  );

  return (
    <>
      <nav aria-label="Main navigation">
        <div className="nav-wrapper container">
          <Link to={isAdminArea ? "/admin" : "/"} className="brand-logo">
            {site.has_logo
              ? <img src={`/api/site/logo?v=${logoVersion}`} alt="" style={{ height: 36, verticalAlign: "middle", marginRight: 8 }} />
              : <i className="material-icons">spa</i>}
            {site.name}
          </Link>
          <a href="#" data-target="mobile-menu" className="sidenav-trigger right" aria-label="Open navigation menu">
            <i className="material-icons" style={{ color: "#18223c" }}>menu</i>
          </a>
          <ul className="right hide-on-med-and-down">{links.map(renderLink)}{langToggle}</ul>
        </div>
      </nav>
      <ul className="sidenav" id="mobile-menu">{links.map(renderLink)}{langToggle}</ul>

      <main className={pathname === "/" ? undefined : "editor-shell"}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Card><Register /></Card>} />
          <Route path="/login" element={<Card><Login /></Card>} />
          <Route path="/verify" element={<Card><Verify /></Card>} />
          <Route path="/forgot-password" element={<Card><ForgotPassword /></Card>} />
          <Route path="/reset-password" element={<Card><ResetPassword /></Card>} />
          <Route path="/account" element={<Account />} />
          <Route path="/contact" element={<Card><Contact /></Card>} />
          <Route path="/admin" element={<Admin />} />
        </Routes>
      </main>
    </>
  );
}

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="card" style={{ maxWidth: 560, margin: "2rem auto" }}>
      <div className="card-content">{children}</div>
    </div>
  );
}
