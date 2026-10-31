import { useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { auth } from "../api";
import { useAuth } from "../auth";
import { useLang } from "../lang";
import { TextField, useSubmit } from "../ui";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { setUser } = useAuth();
  const { t } = useLang();
  const { onSubmit, busy } = useSubmit(async () => {
    const u = await auth.login(email, password);
    setUser(u);
    navigate(params.get("next") ?? (u.role === "user" ? "/account" : "/admin"));
  });

  return (
    <>
      <span className="card-title">{t("login_title")}</span>
      <form onSubmit={onSubmit} className="section">
        <TextField label={t("f_email")} type="email" value={email} onChange={setEmail} />
        <TextField label={t("f_password")} type="password" value={password} onChange={setPassword} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          {t("login_button")} <i className="material-icons right">login</i>
        </button>
        <Link className="btn-flat waves-effect" to="/forgot-password">{t("forgot_link")}</Link>
      </form>
      <p><Link to="/register">{t("need_account")}</Link></p>
    </>
  );
}
