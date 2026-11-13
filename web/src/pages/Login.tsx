import { useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { auth } from "../api";
import { useAuth } from "../auth";
import { useSite } from "../site";
import { useLang } from "../lang";
import { TextField, useSubmit } from "../ui";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { setUser } = useAuth();
  const { t } = useLang();
  const { site } = useSite();
  const { onSubmit, busy } = useSubmit(async () => {
    const u = await auth.login(email, password);
    setUser(u);
    navigate(params.get("next") ?? (u.role === "user" ? "/account" : "/admin"));
  });

  return (
    <>
      <h5 className="card-title">{t("login_title")}</h5>
      <form onSubmit={onSubmit} className="section">
        <TextField label={t("f_email")} type="email" value={email} onChange={setEmail} />
        <TextField label={t("f_password")} type="password" value={password} onChange={setPassword} />
        <div className="d-flex justify-content-between align-items-center">
          <button className="btn btn-primary" type="submit" disabled={busy}>
            {t("login_button")} <i className="material-icons align-middle ms-1">login</i>
          </button>

          <Link className="btn btn-danger" to="/forgot-password">
            {t("forgot_link")}
          </Link>
        </div>
      </form>

      {site.registration && (
        <p>
          <Link to="/register">{t("need_account")}</Link>
        </p>
      )}
    </>
  );
}
