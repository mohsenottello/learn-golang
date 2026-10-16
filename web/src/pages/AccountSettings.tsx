import { useState } from "react";
import { auth } from "../api";
import { useAuth } from "../auth";
import { ErrorPanel, TextField, useSubmit } from "../ui";

export default function AccountSettings({ onNameChange }: { onNameChange: (name: string) => void }) {
  const { user, setUser } = useAuth();
  const [name, setName] = useState(user?.name ?? "");
  const [nameSaved, setNameSaved] = useState(false);
  const nameForm = useSubmit(async () => {
    const u = await auth.updateName(name);
    setUser(u);
    onNameChange(u.name);
    setNameSaved(true);
  });

  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [confirm, setConfirm] = useState("");
  const [pwMsg, setPwMsg] = useState<string | null>(null);
  const pwForm = useSubmit(async () => {
    if (next !== confirm) throw new Error("New passwords don't match.");
    const r = await auth.changePassword(current, next);
    setPwMsg(r.message);
    setCurrent(""); setNext(""); setConfirm("");
  });

  return (
    <>
      <div className="card">
        <div className="card-content">
          <span className="card-title">Account</span>
          <p className="grey-text">Signed in as <strong>{user?.email}</strong></p>
          <ErrorPanel message={nameForm.error} />
          {nameSaved && <div className="card-panel teal lighten-5 teal-text text-darken-3">Name updated.</div>}
          <form onSubmit={(e) => { setNameSaved(false); nameForm.onSubmit(e); }}>
            <TextField label="Display name (shown on your resume)" type="text" value={name} onChange={setName} />
            <button className="btn teal" type="submit" disabled={nameForm.busy}>Save name</button>
          </form>
        </div>
      </div>

      <div className="card">
        <div className="card-content">
          <span className="card-title">Change password</span>
          <ErrorPanel message={pwForm.error} />
          {pwMsg && <div className="card-panel teal lighten-5 teal-text text-darken-3">{pwMsg}</div>}
          <form onSubmit={(e) => { setPwMsg(null); pwForm.onSubmit(e); }}>
            <TextField label="Current password" type="password" value={current} onChange={setCurrent} />
            <TextField label="New password" type="password" value={next} onChange={setNext} minLength={8} />
            <TextField label="Confirm new password" type="password" value={confirm} onChange={setConfirm} minLength={8} />
            <button className="btn teal" type="submit" disabled={pwForm.busy}>Change password</button>
          </form>
        </div>
      </div>
    </>
  );
}
