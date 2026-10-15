import { Link, Route, Routes } from "react-router-dom";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Verify from "./pages/Verify";
import Dashboard from "./pages/Dashboard";

export default function App() {
  return (
    <>
      <nav className="teal">
        <div className="nav-wrapper container">
          <Link to="/" className="brand-logo">MyApp</Link>
          <ul className="right">
            <li><Link to="/admin/login">Login</Link></li>
            <li><Link to="/admin/register">Register</Link></li>
          </ul>
        </div>
      </nav>
      <main className="container" style={{ maxWidth: 480, marginTop: "3rem" }}>
        <div className="card">
          <div className="card-content">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/admin/login" element={<Login />} />
              <Route path="/admin/register" element={<Register />} />
              <Route path="/admin/verify" element={<Verify />} />
              <Route path="/admin" element={<Dashboard />} />
            </Routes>
          </div>
        </div>
      </main>
    </>
  );
}
