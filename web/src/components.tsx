import { FormEvent, ReactNode, useId, useState } from "react";

export function TextField(props: {
  label: string;
  type: string;
  value: string;
  onChange: (v: string) => void;
  minLength?: number;
}) {
  const id = useId();
  return (
    <div className="input-field">
      <input id={id} type={props.type} value={props.value} minLength={props.minLength} required
        onChange={(e) => props.onChange(e.target.value)} />
      <label htmlFor={id} className={props.value ? "active" : ""}>{props.label}</label>
    </div>
  );
}

export function ErrorPanel({ message }: { message: string | null }) {
  if (!message) return null;
  return <div className="card-panel red lighten-4 red-text text-darken-4">{message}</div>;
}

export function Message({ title, children }: { title: string; children: ReactNode }) {
  return (
    <>
      <span className="card-title">{title}</span>
      {children}
    </>
  );
}

/** Wraps a submit handler with busy/error state so every form looks the same. */
export function useSubmit(action: () => Promise<void>) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await action();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }
  return { onSubmit, busy, error };
}
