import { useState } from "react";
import { Link } from "react-router-dom";
import { auth } from "../api";
import { Message, TextField, useSubmit } from "../ui";
import { useSite } from "../site";

export default function Register() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { onSubmit, busy } = useSubmit(async () => {
    const r = await auth.register(email, name, password);
    setDone(r.message);
  });

  if (done) return <Message title="Check your email"><p>{done}</p></Message>;
  if (!useSite().site.registration) return <Message title="Registration disabled"><p>New user registration is currently disabled.</p></Message>;

  return (
    <>
      <span className="card-title">Create account</span>
      <form onSubmit={onSubmit}>
        <TextField label="Name" type="text" value={name} onChange={setName} />
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <TextField label="Password" type="password" value={password} onChange={setPassword} minLength={8} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Register <i className="material-icons right">person_add</i>
        </button>
        <Link className="btn-flat waves-effect" to="/login">Already have an account? Log in</Link>
      </form>
    </>
  );
}
