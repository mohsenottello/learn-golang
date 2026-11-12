import { useEffect, useState } from "react";
import { User, admin } from "../api";
import { useToast } from "../toast";
import { TextField, useSubmit } from "../ui";
import { useModal } from "../components/useModal";
import { useAuth } from "../auth";
export default function AdminUsers() {
  const [q, setQ] = useState("");
  const [users, setUsers] = useState<User[]>([]);
  const [editing, setEditing] = useState<User | null>(null);
  const [pwUser, setPwUser] = useState<User | null>(null);
  const toast = useToast();

  const load = () => admin.listUsers(q).then(setUsers).catch((e) => toast(e.message, "error"));
  const { user: me } = useAuth()

  // Debounced search: wait 300 ms after the last keystroke before querying.
  useEffect(() => {
    const t = setTimeout(load, 300);
    return () => clearTimeout(t);
  }, [q]); // eslint-disable-line react-hooks/exhaustive-deps

  async function saveEdit() {
    if (!editing) return;
    try {
      await admin.updateUser(editing);
      setEditing(null);
      toast("User updated.");
      load();
    } catch (e) { toast((e as Error).message, "error"); }
  }

  async function remove(u: User) {
    if (!confirm(`Delete ${u.email} and their resume?`)) return;
    try { await admin.deleteUser(u.id); toast("User deleted."); load(); }
    catch (e) { toast((e as Error).message, "error"); }
  }

  const date = (s: string) => new Date(s).toLocaleDateString();

  return (
    <div>
      <div className="card mb-3">
        <div className="card-body">
          <h5 className="card-title">Users</h5>
          <div className="input-group mb-3" style={{ maxWidth: 360 }}>
            <span className="input-group-text"><i className="material-icons">search</i></span>
            <input type="text" className="form-control" placeholder="Search name or email" value={q} onChange={(e) => setQ(e.target.value)} />
          </div>

          <table className="table table-hover align-middle">
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Verified</th>
                <th>Joined</th>
                <th style={{ width: 150 }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => editing?.id === u.id ? (
                <tr key={u.id} className="table-active">
                  <td><input type="text" className="form-control form-control-sm" value={editing.name} onChange={(e) => setEditing({ ...editing, name: e.target.value })} /></td>
                  <td><input type="email" className="form-control form-control-sm" value={editing.email} onChange={(e) => setEditing({ ...editing, email: e.target.value })} /></td>
                  <td>
                    <select className="form-select form-select-sm" value={editing.role} onChange={(e) => setEditing({ ...editing, role: e.target.value as User["role"] })}>
                      <option value="user">user</option>
                      <option value="admin">admin</option>
                      <option value="owner">owner</option>
                    </select>
                  </td>
                  <td>
                    <input type="checkbox" className="form-check-input" checked={editing.verified}
                      onChange={(e) => setEditing({ ...editing, verified: e.target.checked })} />
                  </td>
                  <td>{date(u.created_at)}</td>
                  <td>
                    <button className="btn btn-link btn-sm p-1" onClick={saveEdit}><i className="material-icons">check</i></button>
                    <button className="btn btn-link btn-sm p-1" onClick={() => setEditing(null)}><i className="material-icons">close</i></button>
                  </td>
                </tr>
              ) : (
                <tr key={u.id}>
                  <td>{u.name}</td>
                  <td>{u.email}</td>
                  <td>{u.role}</td>
                  <td>{u.verified
                    ? <i className="material-icons text-primary">check</i>
                    : <i className="material-icons text-secondary">hourglass_empty</i>}</td>
                  <td>{date(u.created_at)}</td>
                  <td className="text-end text-nowrap">
                    <button className="btn btn-link btn-sm p-1" title="Edit" onClick={() => setEditing({ ...u })}><i className="material-icons">edit</i></button>
                    <button className="btn btn-link btn-sm p-1" title="Set password" onClick={() => setPwUser(u)}><i className="material-icons">key</i></button>
                    {(u.role !== "owner" || u.id !== me?.id) && (
                      <button className="btn btn-link btn-sm p-1 text-danger" title="Delete" onClick={() => remove(u)}><i className="material-icons">delete</i></button>
                    )}
                  </td>
                </tr>
              ))}
              {users.length === 0 && <tr><td colSpan={5} className="text-secondary">No users found.</td></tr>}
            </tbody>
          </table>
        </div>

        {/* Mounted once; opens whenever pwUser is set */}
        <PasswordModal user={pwUser} onDone={() => setPwUser(null)} />
      </div>
    </div>
  );
}

function PasswordModal({ user, onDone }: { user: User | null; onDone: () => void }) {
  const [pw, setPw] = useState("");
  const [confirm, setConfirm] = useState("");
  const toast = useToast();
  const modal = useModal(onDone);

  useEffect(() => { if (user) { setPw(""); setConfirm(""); modal.open(); } }, [user]); // eslint-disable-line react-hooks/exhaustive-deps

  const form = useSubmit(async () => {
    if (!user) return;
    if (pw !== confirm) throw new Error("Passwords don't match.");
    await admin.setUserPassword(user.id, pw);
    toast(`Password changed for ${user.email}.`);
    modal.close();
  });

  return (
    <div ref={modal.ref} className="modal fade" tabIndex={-1}>
      <div className="modal-dialog">
        <form className="modal-content" onSubmit={form.onSubmit}>
          <div className="modal-header">
            <h5 className="modal-title">Set password{user && <> — {user.email}</>}</h5>
            <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close" />
          </div>
          <div className="modal-body">
            <TextField label="New password" type="password" value={pw} onChange={setPw} minLength={8} />
            <TextField label="Confirm" type="password" value={confirm} onChange={setConfirm} minLength={8} />
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-link" data-bs-dismiss="modal">Cancel</button>
            <button type="submit" className="btn btn-primary" disabled={form.busy || pw.length < 8}>Save</button>
          </div>
        </form>
      </div>
    </div>
  );
}
