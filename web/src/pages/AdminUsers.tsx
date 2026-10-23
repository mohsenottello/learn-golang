import { useEffect, useState } from "react";
import { User, admin } from "../api";
import { useToast } from "../toast";
import { TextField, useSubmit } from "../ui";
import { useModal } from "../components/useModal";

export default function AdminUsers() {
  const [q, setQ] = useState("");
  const [users, setUsers] = useState<User[]>([]);
  const [editing, setEditing] = useState<User | null>(null);
  const [pwUser, setPwUser] = useState<User | null>(null);
  const toast = useToast();

  const load = () => admin.listUsers(q).then(setUsers).catch((e) => toast(e.message, "error"));

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
    <div className="row">
      <div className="card">
        <div className="card-content">
          <span className="card-title">Users</span>
          <div className="input-field" style={{ maxWidth: 360 }}>
            <i className="material-icons prefix">search</i>
            <input id="q" type="text" value={q} onChange={(e) => setQ(e.target.value)} />
            <label htmlFor="q" className={q ? "active" : ""}>Search name or email</label>
          </div>

          <table className="highlight">
            <thead>
              <tr><th>Name</th><th>Email</th><th>Verified</th><th>Joined</th><th style={{ width: 150 }}>Actions</th></tr>
            </thead>
            <tbody>
              {users.map((u) => editing?.id === u.id ? (
                <tr key={u.id} className="teal lighten-5">
                  <td><input type="text" value={editing.name} onChange={(e) => setEditing({ ...editing, name: e.target.value })} /></td>
                  <td><input type="email" value={editing.email} onChange={(e) => setEditing({ ...editing, email: e.target.value })} /></td>
                  <td>
                    <label>
                      <input type="checkbox" className="filled-in" checked={editing.verified}
                        onChange={(e) => setEditing({ ...editing, verified: e.target.checked })} />
                      <span></span>
                    </label>
                  </td>
                  <td>{date(u.created_at)}</td>
                  <td>
                    <button className="btn-flat teal-text" onClick={saveEdit}><i className="material-icons">check</i></button>
                    <button className="btn-flat" onClick={() => setEditing(null)}><i className="material-icons">close</i></button>
                  </td>
                </tr>
              ) : (
                <tr key={u.id}>
                  <td>{u.name}</td>
                  <td>{u.email}</td>
                  <td>{u.verified
                    ? <i className="material-icons teal-text">check</i>
                    : <i className="material-icons grey-text">hourglass_empty</i>}</td>
                  <td>{date(u.created_at)}</td>
                  <td className="right-align" style={{ whiteSpace: "nowrap" }}>
                    <button className="btn-flat" title="Edit" onClick={() => setEditing({ ...u })}><i className="material-icons">edit</i></button>
                    <button className="btn-flat" title="Set password" onClick={() => setPwUser(u)}><i className="material-icons">key</i></button>
                    <button className="btn-flat red-text" title="Delete" onClick={() => remove(u)}><i className="material-icons">delete</i></button>
                  </td>
                </tr>
              ))}
              {users.length === 0 && <tr><td colSpan={5} className="grey-text">No users found.</td></tr>}
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
    <div ref={modal.ref} className="modal">
      <form onSubmit={form.onSubmit}>
        <div className="modal-content">
          <h5>Set password{user && <> — {user.email}</>}</h5>
          <TextField label="New password" type="password" value={pw} onChange={setPw} minLength={8} />
          <TextField label="Confirm" type="password" value={confirm} onChange={setConfirm} minLength={8} />
        </div>
        <div className="modal-footer">
          <button type="button" className="btn-flat modal-close">Cancel</button>
          <button type="submit" className="btn teal" disabled={form.busy || pw.length < 8}>Save</button>
        </div>
      </form>
    </div>
  );
}
