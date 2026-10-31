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
        <Link className="btn teal" to="/login">{t("login_button")}</Link>
      </Message>
    );
  }

  return (
    <>
      <span className="card-title">{t("reset_title")}</span>
      <form onSubmit={onSubmit}>
        <TextField label={t("reset_new")} type="password" value={password} onChange={setPassword} minLength={8} />
        <TextField label={t("reset_confirm")} type="password" value={confirm} onChange={setConfirm} minLength={8} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          {t("reset_save")} <i className="material-icons right">lock_reset</i>
        </button>
      </form>
    </>
  );
}
