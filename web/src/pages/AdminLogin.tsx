import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { admin } from "../api";
import { TextField, useSubmit } from "../ui";

export default function AdminLogin() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const { onSubmit, busy } = useSubmit(async () => {
    await admin.login(email, password);
    navigate("/admin");
  });

  return (
    <>
      <span className="card-title">Admin login</span>
      <form onSubmit={onSubmit}>
        <TextField label="Email" type="email" value={email} onChange={setEmail} />
        <TextField label="Password" type="password" value={password} onChange={setPassword} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          Log in <i className="material-icons right">admin_panel_settings</i>
        </button>
      </form>
    </>
  );
}
