import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Admin, api } from "../api";

export default function Dashboard() {
  const [admin, setAdmin] = useState<Admin | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    api.me().then(setAdmin).catch(() => navigate("/admin/login"));
  }, [navigate]);

  async function logout() {
    await api.logout();
    navigate("/admin/login");
  }

  if (!admin) return <p>Loading…</p>;
  return (
    <>
      <span className="card-title">Admin dashboard</span>
      <p>Logged in as <strong>{admin.email}</strong></p>
      <button className="btn red lighten-1" onClick={logout}>Log out</button>
    </>
  );
}
