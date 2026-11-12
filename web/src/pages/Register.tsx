import { useState } from "react";
import { Link } from "react-router-dom";
import { auth } from "../api";
import { useLang } from "../lang";
import { useSite } from "../site";
import { Message, TextField, useSubmit } from "../ui";

export default function Register() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [done, setDone] = useState<string | null>(null);
  const { site } = useSite();
  const { t } = useLang();
  const { onSubmit, busy } = useSubmit(async () => {
    const r = await auth.register(email, name, password);
    setDone(r.message);
  });

  if (done) return <Message title={t("check_email")}><p>{done}</p></Message>;
  if (!site.registration) return <Message title={t("register_closed_title")}><p>{t("register_closed_text")}</p></Message>;

  return (
    <>
      <h5 className="card-title">{t("register_title")}</h5>
      <form onSubmit={onSubmit}>
        <TextField label={t("f_name")} type="text" value={name} onChange={setName} />
        <TextField label={t("f_email")} type="email" value={email} onChange={setEmail} />
        <TextField label={t("f_password")} type="password" value={password} onChange={setPassword} minLength={8} />
        <button className="btn btn-primary" type="submit" disabled={busy}>
          {t("register_button")} <i className="material-icons align-middle ms-1">person_add</i>
        </button>
        <Link className="btn btn-link" to="/login">{t("have_account")}</Link>
      </form>
    </>
  );
}
