import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { api } from "../api";

export default function Verify() {
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";
  const [state, setState] = useState<"loading" | "ok" | "error">("loading");
  const [message, setMessage] = useState("");

  useEffect(() => {
    api.verify(token)
      .then((r) => { setMessage(r.message); setState("ok"); })
      .catch((e: Error) => { setMessage(e.message); setState("error"); });
  }, [token]);

  if (state === "loading") return <p>Verifying…</p>;
  return (
    <>
      <span className="card-title">{state === "ok" ? "Verified" : "Invalid link"}</span>
      <p>{message}</p>
      {state === "ok" && <Link className="btn teal" to="/admin/login">Log in</Link>}
    </>
  );
}
