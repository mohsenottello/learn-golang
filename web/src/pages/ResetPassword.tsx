import { useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { auth } from "../api";
import { useLang } from "../lang";
import { Message, TextField, useSubmit } from "../ui";

export default function ResetPassword() {
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { t } = useLang();
  const { onSubmit, busy } = useSubmit(async () => {
    if (password !== confirm) throw new Error(t("reset_mismatch"));
    const r = await auth.resetPassword(token, password);
    setDone(r.message);
  });

  if (!token) return <Message title={t("invalid_link")}><p>{t("reset_missing")}</p></Message>;
  if (done) {
    return (
      <Message title={t("reset_done")}>
        <p>{done}</p>
        <Link className="btn btn-primary" to="/login">{t("login_button")}</Link>
      </Message>
    );
  }

  return (
    <>
      <h5 className="card-title">{t("reset_title")}</h5>
      <form onSubmit={onSubmit}>
        <TextField label={t("reset_new")} type="password" value={password} onChange={setPassword} minLength={8} />
        <TextField label={t("reset_confirm")} type="password" value={confirm} onChange={setConfirm} minLength={8} />
        <button className="btn btn-primary" type="submit" disabled={busy}>
          {t("reset_save")} <i className="material-icons align-middle ms-1">lock_reset</i>
        </button>
      </form>
    </>
  );
}
