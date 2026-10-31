import { useState } from "react";
import { Link } from "react-router-dom";
import { auth } from "../api";
import { useLang } from "../lang";
import { Message, TextField, useSubmit } from "../ui";

export default function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { t } = useLang();
  const { onSubmit, busy } = useSubmit(async () => {
    const r = await auth.forgotPassword(email);
    setDone(r.message);
  });

  if (done) return <Message title={t("check_email")}><p>{done}</p></Message>;

  return (
    <>
      <span className="card-title">{t("forgot_title")}</span>
      <p>{t("forgot_text")}</p>
      <form onSubmit={onSubmit}>
        <TextField label={t("f_email")} type="email" value={email} onChange={setEmail} />
        <button className="btn teal waves-effect waves-light" type="submit" disabled={busy}>
          {t("forgot_button")} <i className="material-icons right">send</i>
        </button>
        <Link className="btn-flat waves-effect" to="/login">{t("back_to_login")}</Link>
      </form>
    </>
  );
}
