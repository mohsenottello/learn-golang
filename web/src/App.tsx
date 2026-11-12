import { ReactNode, useEffect, useRef } from "react";
import { Link, NavLink as RouterNavLink, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { Offcanvas } from "bootstrap";
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

type NavLink = { to: string; label: string; onClick?: () => void; cta?: boolean };

export default function App() {
  const { pathname, search } = useLocation();
  const { user, logout } = useAuth();
  const { site, logoVersion } = useSite();
  const { lang, setLang, t } = useLang();
  const navigate = useNavigate();
  const isAdminArea = pathname.startsWith("/admin");
  const print = new URLSearchParams(search).has("print");
  const drawer = useRef<HTMLDivElement>(null);

  // Close the mobile drawer whenever the route changes.
  useEffect(() => { if (drawer.current) Offcanvas.getInstance(drawer.current)?.hide(); }, [pathname]);

  const userLogout = () => logout().finally(() => navigate("/login"));
  const isStaff = user?.role === "admin" || user?.role === "owner";

  const links: NavLink[] = [];
  if (isAdminArea) {
    if (pathname === "/admin") links.push({ to: "#", label: t("nav_logout"), onClick: userLogout });
  } else {
    links.push({ to: "/", label: t("nav_home") });
    if (site.owner_slug) links.push({ to: "/resume", label: t("nav_resume") });
    if (user) links.push({ to: isStaff ? "/admin" : "/account", label: isStaff ? t("nav_admin") : t("nav_editor") });
    links.push({ to: "/contact", label: t("nav_contact") });
    if (user) {
      links.push({ to: "#", label: t("nav_logout"), onClick: userLogout });
    } else if (site.registration) {
      links.push({ to: "/login", label: t("nav_login") });
      links.push({ to: "/register", label: t("nav_register"), cta: true });
    }
  }

  const renderLink = (l: NavLink) => (
    <li key={l.label} className="nav-item">
      {l.onClick
        ? <a href={l.to} className="nav-link" onClick={(e) => { e.preventDefault(); l.onClick!(); }}>{l.label}</a>
        : l.cta
          ? <Link to={l.to} className="btn btn-primary btn-sm rounded-pill px-3 ms-lg-2 nav-cta">{l.label}</Link>
          : <RouterNavLink to={l.to} end className="nav-link">{l.label}</RouterNavLink>}
    </li>
  );

  const langToggle = (
    <li className="nav-item lang-toggle">
      {(["de", "en"] as Lang[]).map((l) => (
        <a key={l} href="#" className={lang === l ? "active" : ""} onClick={(e) => { e.preventDefault(); setLang(l); }}>
          {l.toUpperCase()}
        </a>
      ))}
    </li>
  );

  return (
    <>
      {!print && (
        <nav className="navbar navbar-expand-lg site-nav sticky-top" aria-label="Main navigation">
          <div className="container">
            <Link to={isAdminArea ? "/admin" : "/"} className="navbar-brand brand-logo">
              {site.has_logo
                ? <img src={`/api/site/logo?v=${logoVersion}`} alt={site.name} height={36} className="me-2 align-middle" />
                : <i className="material-icons align-middle me-2">spa</i>}
              {site.name}
            </Link>
            <button className="navbar-toggler" type="button" data-bs-toggle="offcanvas" data-bs-target="#mobile-menu"
              aria-controls="mobile-menu" aria-label="Open navigation menu">
              <span className="navbar-toggler-icon" />
            </button>

            <div ref={drawer} className="offcanvas offcanvas-end" tabIndex={-1} id="mobile-menu">
              <div className="offcanvas-header">
                <span className="offcanvas-title fw-bold">{site.name}</span>
                <button type="button" className="btn-close" data-bs-dismiss="offcanvas" aria-label="Close" />
              </div>
              <div className="offcanvas-body">
                <ul className="navbar-nav ms-auto align-items-lg-center gap-lg-1">
                  {links.map(renderLink)}
                  {langToggle}
                </ul>
              </div>
            </div>
          </div>
        </nav>
      )}

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
          <Route path="/resume" element={<PublicResume featured />} />
          <Route path="/resume/:lang" element={<PublicResume featured />} />
          <Route path="/:slug/:lang" element={<PublicResume />} />
          <Route path="/:slug" element={<PublicResume />} />
        </Routes>
      </main>
    </>
  );
}

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="card shadow-sm mx-auto my-4" style={{ maxWidth: 560 }}>
      <div className="card-body p-4">{children}</div>
    </div>
  );
}
