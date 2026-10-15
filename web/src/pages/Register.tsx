import { FormEvent, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";
import { ErrorPanel, TextField } from "../components";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      const res = await api.register(email, password);
      setDone(res.message);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  if (done) {
    return (
      <>
        <span className="card-title">Check your email</span>
        <p>{done}</p>
      </>
    );
  }

  return (
    <>
      <span className="card-title">Register admin</span>
      <ErrorPanel message={error} />
      <form onSubmit={submit}>
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <TextField label="Password" type="password" value={password} onChange={setPassword} minLength={8} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Register <i className="material-icons right">person_add</i>
        </button>
        <Link className="btn-flat waves-effect" to="/admin/login">Already have an account? Log in</Link>
      </form>
    </>
  );
}
