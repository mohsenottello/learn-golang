import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { auth } from "../api";
import { useAuth } from "../auth";
import { ErrorPanel, TextField, useSubmit } from "../ui";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const { setUser } = useAuth();
  const { onSubmit, busy, error } = useSubmit(async () => {
    setUser(await auth.login(email, password));
    navigate("/account");
  });

  return (
    <>
      <span className="card-title">Log in</span>
      <ErrorPanel message={error} />
      <form onSubmit={onSubmit}>
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <TextField label="Password" type="password" value={password} onChange={setPassword} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Log in <i className="material-icons right">login</i>
        </button>
        <Link className="btn-flat waves-effect" to="/forgot-password">Forgot password?</Link>
      </form>
      <p><Link to="/register">Need an account? Register</Link></p>
    </>
  );
}
