import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { auth } from "../api";
import { useLang } from "../lang";
import { Message } from "../ui";

export default function Verify() {
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";
  const [state, setState] = useState<"loading" | "ok" | "error">("loading");
  const [message, setMessage] = useState("");
  const { t } = useLang();

  useEffect(() => {
    auth.verify(token)
      .then((r) => { setMessage(r.message); setState("ok"); })
      .catch((e: Error) => { setMessage(e.message); setState("error"); });
  }, [token]);

  if (state === "loading") return <p>{t("verify_wait")}</p>;
  return (
    <Message title={state === "ok" ? t("verify_ok") : t("verify_bad")}>
      <p>{message}</p>
      {state === "ok" && <Link className="btn teal" to="/login?next=/account">{t("login_button")}</Link>}
    </Message>
  );
}
