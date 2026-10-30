import { ReactNode, useEffect } from "react";
import { Link, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "./auth";
import { useSite } from "./site";
import { admin } from "./api";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Verify from "./pages/Verify";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import Account from "./pages/Account";
import AdminLogin from "./pages/AdminLogin";
import PublicResume from "./pages/PublicResume";
import Admin from "./pages/Admin";
import Contact from "./pages/Contact";

const appPaths = ["/", "/register", "/login", "/verify", "/forgot-password", "/reset-password", "/account", "/contact", "/admin/login", "/admin"];

type NavLink = { to: string; label: string; onClick?: () => void; cta?: boolean };

export default function App() {
  const { pathname } = useLocation();
  const { user, logout } = useAuth();
  const { site, logoVersion } = useSite();
  const registration = site.registration;
  const navigate = useNavigate();
  const isAdminArea = pathname.startsWith("/admin");
  const resumeHref = `/${site.home.featured_slug}`;

  // Materialize's mobile drawer needs a one-time init after the DOM exists.
  useEffect(() => { M.Sidenav.init(document.querySelectorAll(".sidenav")); }, []);

  const userLogout = () => logout().finally(() => navigate("/login"));

  const links: NavLink[] = [];

  if (isAdminArea) {
    if (pathname === "/admin") {
      links.push({ to: "#", label: "Log out", onClick: userLogout });
    }
  } else if (!registration) {
    links.push(
      { to: "/", label: "Home" },
      { to: "/contact", label: "Contact" },
    );
  } else if (user) {
    links.push(
      { to: "/", label: "Home" },
      { to: resumeHref, label: "My Resume" },
      { to: "/account", label: "My editor" },
      { to: "/contact", label: "Contact" },
      { to: "#", label: "Log out", onClick: userLogout },
    );
  } else {
    links.push(
      { to: "/", label: "Home" },
      { to: resumeHref, label: "My Resume" },
      { to: "/contact", label: "Contact" },
      { to: "/login", label: "Login" },
      { to: "/register", label: "Get Started", cta: true },
    );
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
          <ul className="right hide-on-med-and-down">{links.map(renderLink)}</ul>
        </div>
      </nav>
      <ul className="sidenav" id="mobile-menu">{links.map(renderLink)}</ul>

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
          <Route path="/admin/login" element={<Card><AdminLogin /></Card>} />
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
