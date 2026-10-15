import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api";
import { ErrorPanel, TextField } from "../components";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const navigate = useNavigate();

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await api.login(email, password);
      navigate("/admin");
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <span className="card-title">Admin login</span>
      <ErrorPanel message={error} />
      <form onSubmit={submit}>
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <TextField label="Password" type="password" value={password} onChange={setPassword} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Log in <i className="material-icons right">login</i>
        </button>
        <Link className="btn-flat waves-effect" to="/admin/register">Need an account? Register</Link>
      </form>
    </>
  );
}
