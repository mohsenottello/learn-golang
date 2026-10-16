import { ReactNode } from "react";
import { Link, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "./auth";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Verify from "./pages/Verify";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import Account from "./pages/Account";
import AdminLogin from "./pages/AdminLogin";
import AdminDashboard from "./pages/AdminDashboard";
import PublicResume from "./pages/PublicResume";

const appPaths = ["/", "/register", "/login", "/verify", "/forgot-password", "/reset-password", "/account", "/admin/login", "/admin"];

export default function App() {
  const { pathname } = useLocation();
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  // Anything that isn't an app route is a public resume, rendered without the app chrome.
  if (!appPaths.includes(pathname)) {
    return (
      <Routes>
        <Route path="/:slug/:lang" element={<PublicResume />} />
        <Route path="/:slug" element={<PublicResume />} />
      </Routes>
    );
  }

  return (
    <>
      <nav className="teal">
        <div className="nav-wrapper container">
          <Link to="/" className="brand-logo">MyApp</Link>
          <ul className="right">
            {user ? (
              <>
                <li><Link to="/account">My resume</Link></li>
                <li><a href="#logout" onClick={(e) => { e.preventDefault(); logout().then(() => navigate("/login")); }}>Log out</a></li>
              </>
            ) : (
              <>
                <li><Link to="/login">Login</Link></li>
                <li><Link to="/register">Register</Link></li>
              </>
            )}
          </ul>
        </div>
      </nav>
      <main className={pathname === "/account" ? "editor-shell" : "container"} style={pathname === "/account" ? undefined : { maxWidth: 820, marginTop: "2rem" }}>
        <Routes>
          <Route path="/" element={<Card><Home /></Card>} />
          <Route path="/register" element={<Card><Register /></Card>} />
          <Route path="/login" element={<Card><Login /></Card>} />
          <Route path="/verify" element={<Card><Verify /></Card>} />
          <Route path="/forgot-password" element={<Card><ForgotPassword /></Card>} />
          <Route path="/reset-password" element={<Card><ResetPassword /></Card>} />
          <Route path="/account" element={<Account />} />
          <Route path="/admin/login" element={<Card><AdminLogin /></Card>} />
          <Route path="/admin" element={<Card><AdminDashboard /></Card>} />
        </Routes>
      </main>
    </>
  );
}

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="card" style={{ maxWidth: 560, margin: "0 auto" }}>
      <div className="card-content">{children}</div>
    </div>
  );
}
