import { useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { auth } from "../api";
import { ErrorPanel, Message, TextField, useSubmit } from "../components";

export default function ResetPassword() {
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { onSubmit, busy, error } = useSubmit(async () => {
    if (password !== confirm) throw new Error("Passwords don't match.");
    const r = await auth.resetPassword(token, password);
    setDone(r.message);
  });

  if (!token) return <Message title="Invalid link"><p>This reset link is missing its token.</p></Message>;
  if (done) {
    return (
      <Message title="Password updated">
        <p>{done}</p>
        <Link className="btn teal" to="/login">Log in</Link>
      </Message>
    );
  }

  return (
    <>
      <span className="card-title">Choose a new password</span>
      <ErrorPanel message={error} />
      <form onSubmit={onSubmit}>
        <TextField label="New password" type="password" value={password} onChange={setPassword} minLength={8} />
        <TextField label="Confirm password" type="password" value={confirm} onChange={setConfirm} minLength={8} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Save <i className="material-icons right">lock_reset</i>
        </button>
      </form>
    </>
  );
}
