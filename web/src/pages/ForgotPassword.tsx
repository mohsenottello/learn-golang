import { useState } from "react";
import { Link } from "react-router-dom";
import { auth } from "../api";
import { ErrorPanel, Message, TextField, useSubmit } from "../ui";

export default function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { onSubmit, busy, error } = useSubmit(async () => {
    const r = await auth.forgotPassword(email);
    setDone(r.message);
  });

  if (done) return <Message title="Check your email"><p>{done}</p></Message>;

  return (
    <>
      <span className="card-title">Reset your password</span>
      <p>Enter your email and we'll send you a link to choose a new password.</p>
      <ErrorPanel message={error} />
      <form onSubmit={onSubmit}>
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Send link <i className="material-icons right">send</i>
        </button>
        <Link className="btn-flat waves-effect" to="/login">Back to login</Link>
      </form>
    </>
  );
}
