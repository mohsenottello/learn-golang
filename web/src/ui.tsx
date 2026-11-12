import { FormEvent, ReactNode, useId, useState } from "react";
import { useToast } from "./toast";

/** Floating-label text input (Bootstrap form-floating). */
export function TextField(props: {
  label: string;
  type: string;
  value: string;
  onChange: (v: string) => void;
  minLength?: number;
}) {
  const id = useId();
  return (
    <div className="form-floating mb-3">
      <input id={id} className="form-control" type={props.type} value={props.value} minLength={props.minLength} required
        placeholder={props.label} onChange={(e) => props.onChange(e.target.value)} />
      <label htmlFor={id}>{props.label}</label>
    </div>
  );
}

export function Message({ title, children }: { title: string; children: ReactNode }) {
  return (
    <>
      <h5 className="card-title">{title}</h5>
      {children}
    </>
  );
}

/** Wraps a submit handler with busy/error state so every form looks the same. */
export function useSubmit(action: () => Promise<void>) {
  const [busy, setBusy] = useState(false);
  const toast = useToast();
  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    try {
      await action();
    } catch (err) {
      toast((err as Error).message, "error");
    } finally {
      setBusy(false);
    }
  }
  return { onSubmit, busy };
}
