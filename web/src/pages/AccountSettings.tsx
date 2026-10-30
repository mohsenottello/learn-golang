import { useState } from "react";
import { auth } from "../api";
import { useAuth } from "../auth";
import { useToast } from "../toast";
import { TextField, useSubmit } from "../ui";

export default function AccountSettings() {
  const { user, setUser } = useAuth();
  const toast = useToast();

  const [name, setName] = useState(user?.name ?? "");

  const nameForm = useSubmit(async () => {
    const u = await auth.updateName(name);
    setUser(u);          // sidebar/nav update from the auth context
    toast("Name updated.");
  });

  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [confirm, setConfirm] = useState("");
  const pwForm = useSubmit(async () => {
    if (next !== confirm) throw new Error("New passwords don't match.");
    const r = await auth.changePassword(current, next);
    toast(r.message);
    setCurrent(""); setNext(""); setConfirm("");
  });

  return (
    <>
      <div className="card">
        <div className="card-content">
          <span className="card-title">Account</span>
          <p className="grey-text">Signed in as <strong>{user?.email}</strong></p>
          <form onSubmit={nameForm.onSubmit}>
            <TextField label="Display name (shown on your resume)" type="text" value={name} onChange={setName} />
            <button className="btn teal" type="submit" disabled={nameForm.busy}>Save name</button>
          </form>
        </div>
      </div>

      <div className="card">
        <div className="card-content">
          <span className="card-title">Change password</span>
          <form onSubmit={pwForm.onSubmit}>
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
