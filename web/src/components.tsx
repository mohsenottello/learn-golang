import { useId } from "react";

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
      <input
        id={id}
        type={props.type}
        value={props.value}
        minLength={props.minLength}
        required
        onChange={(e) => props.onChange(e.target.value)}
      />
      {/* Materialize floats the label via the "active" class; React manages it here */}
      <label htmlFor={id} className={props.value ? "active" : ""}>{props.label}</label>
    </div>
  );
}

export function ErrorPanel({ message }: { message: string | null }) {
  if (!message) return null;
  return <div className="card-panel red lighten-4 red-text text-darken-4">{message}</div>;
}
