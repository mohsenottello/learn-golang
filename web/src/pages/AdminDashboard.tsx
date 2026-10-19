import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Admin, User, admin } from "../api";
import { ErrorPanel } from "../ui";

export default function AdminDashboard() {
  const [me, setMe] = useState<Admin | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    admin.me()
      .then((a) => { setMe(a); return admin.listUsers(); })
      .then(setUsers)
      .catch(() => navigate("/admin/login"));
  }, [navigate]);

  async function remove(u: User) {
    if (!confirm(`Delete ${u.email}?`)) return;
    try {
      await admin.deleteUser(u.id);
      setUsers((list) => list.filter((x) => x.id !== u.id));
    } catch (e) {
      setError((e as Error).message);
    }
  }

  if (!me) return <p>Loading…</p>;
  return (
    <>
      <span className="card-title">Users</span>
      <p>Admin: <strong>{me.email}</strong></p>
      <table className="striped">
        <thead>
          <tr><th>Name</th><th>Email</th><th>Verified</th><th>Joined</th><th></th></tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id}>
              <td>{u.name}</td>
              <td>{u.email}</td>
              <td>{u.verified ? <i className="material-icons teal-text">check</i> : <i className="material-icons grey-text">hourglass_empty</i>}</td>
              <td>{new Date(u.created_at).toLocaleDateString()}</td>
              <td><button className="btn-flat red-text" onClick={() => remove(u)}><i className="material-icons">delete</i></button></td>
            </tr>
          ))}
          {users.length === 0 && <tr><td colSpan={5}>No users yet.</td></tr>}
        </tbody>
      </table>
      <br />
      <button className="btn red lighten-1" onClick={() => admin.logout().then(() => navigate("/admin/login"))}>
        Log out
      </button>
    </>
  );
}
