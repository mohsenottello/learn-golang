import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { User, auth } from "../api";

export default function Account() {
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    auth.me().then(setUser).catch(() => navigate("/login"));
  }, [navigate]);

  if (!user) return <p>Loading…</p>;
  return (
    <>
      <span className="card-title">Hello, {user.name}</span>
      <p>Signed in as <strong>{user.email}</strong></p>
      <button className="btn red lighten-1" onClick={() => auth.logout().then(() => navigate("/login"))}>
        Log out
      </button>
    </>
  );
}
