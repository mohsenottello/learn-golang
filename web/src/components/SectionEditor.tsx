import { DragEvent, useState } from "react";
import { ErrorPanel } from "../ui";

export type Field<T> = {
  key: keyof T & string;
  label: string;
  type?: "text" | "textarea" | "month";
  required?: boolean;
  half?: boolean;
  /** Render a second input for `${key}_en` next to the default-language one. */
  translatable?: boolean;
};

type Props<T extends { id: number }> = {
  title: string;
  items: T[];
  fields: Field<T>[];
  summary: (item: T) => string;
  onCreate: (draft: Omit<T, "id" | "position">) => Promise<void>;
  onUpdate: (item: T) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
  onReorder: (ids: number[]) => Promise<void>;
  /** Optional: clicking a row selects it (used to open sub-lists, e.g. bullets of an experience). */
  selectedId?: number | null;
  onSelect?: (id: number) => void;
  flat?: boolean; // render without the outer card (for nested use)
};

/** Generic list + form editor used by all resume sections. Rows are drag-sortable. */
export default function SectionEditor<T extends { id: number }>(p: Props<T>) {
  const [editing, setEditing] = useState<Partial<T> | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [order, setOrder] = useState<T[] | null>(null);
  const [dragId, setDragId] = useState<number | null>(null);

  const items = order ?? p.items;

  const empty = () => {
    const d: Record<string, unknown> = {};
    for (const f of p.fields) {
      d[f.key] = "";
      if (f.translatable) d[`${f.key}_en`] = "";
    }
    return d as Partial<T>;
  };

  async function save() {
    setBusy(true); setError(null);
    try {
      for (const f of p.fields) {
        const v = editing![f.key];
        if (f.required && (!v || (typeof v === "string" && !v.trim()))) throw new Error(`${f.label} is required.`);
      }
      if ((editing as T).id) await p.onUpdate(editing as T);
      else await p.onCreate(editing as Omit<T, "id" | "position">);
      setEditing(null);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setBusy(false);
    }
  }

  async function remove(id: number) {
    if (!confirm("Delete this entry?")) return;
    try { await p.onDelete(id); } catch (e) { setError((e as Error).message); }
  }

  function onDragStart(e: DragEvent, id: number) {
    setDragId(id); setOrder([...p.items]); e.dataTransfer.effectAllowed = "move";
  }
  function onDragOver(e: DragEvent, overId: number) {
    e.preventDefault();
    if (dragId === null || dragId === overId || !order) return;
    const from = order.findIndex((x) => x.id === dragId);
    const to = order.findIndex((x) => x.id === overId);
    if (from < 0 || to < 0) return;
    const next = [...order];
    const [moved] = next.splice(from, 1);
    next.splice(to, 0, moved);
    setOrder(next);
  }
  async function onDrop() {
    if (!order) return;
    const ids = order.map((x) => x.id);
    setDragId(null);
    try { await p.onReorder(ids); } catch (e) { setError((e as Error).message); } finally { setOrder(null); }
  }

  const set = (key: string, value: string) => setEditing({ ...editing!, [key]: value });

  const input = (key: string, label: string, f: Field<T>) => {
    const val = (editing![key as keyof T] as unknown as string) ?? "";
    if (f.type === "textarea") {
      return (<>
        <textarea className="materialize-textarea" value={val} onChange={(e) => set(key, e.target.value)} />
        <label className="active">{label}</label>
      </>);
    }
    return (<>
      <input type={f.type === "month" ? "month" : "text"} value={val} onChange={(e) => set(key, e.target.value)} />
      <label className="active">{label}</label>
    </>);
  };

  const body = (
    <>
      <ErrorPanel message={error} />
      {items.length === 0 && !editing && <p className="grey-text">Nothing here yet.</p>}
      <ul className="collection">
        {items.map((it) => {
          const selected = p.selectedId === it.id;
          return (
            <li key={it.id} className={`collection-item${selected ? " active teal lighten-5" : ""}`}
              draggable
              onDragStart={(e) => onDragStart(e, it.id)}
              onDragOver={(e) => onDragOver(e, it.id)}
              onDrop={onDrop}
              onDragEnd={() => { setDragId(null); setOrder(null); }}
              onClick={() => p.onSelect?.(it.id)}
              style={{ display: "flex", alignItems: "center", opacity: dragId === it.id ? 0.4 : 1,
                cursor: p.onSelect ? "pointer" : "grab", color: selected ? "#004d40" : undefined }}>
              <i className="material-icons grey-text" style={{ marginRight: 8, cursor: "grab" }}>drag_indicator</i>
              <span style={{ flex: 1 }}>{p.summary(it)}</span>
              <button className="btn-flat" onClick={(e) => { e.stopPropagation(); setEditing({ ...it }); }}><i className="material-icons">edit</i></button>
              <button className="btn-flat red-text" onClick={(e) => { e.stopPropagation(); remove(it.id); }}><i className="material-icons">delete</i></button>
            </li>
          );
        })}
      </ul>

      {editing ? (
        <div className="row" style={{ marginBottom: 0 }}>
          {p.fields.map((f) =>
            f.translatable ? (
              <div key={f.key} className="col s12">
                <div className="row" style={{ marginBottom: 0 }}>
                  <div className="input-field col s12 m6">{input(f.key, `${f.label} (DE)${f.required ? " *" : ""}`, f)}</div>
                  <div className="input-field col s12 m6">{input(`${f.key}_en`, `${f.label} (EN)`, f)}</div>
                </div>
              </div>
            ) : (
              <div key={f.key} className={`input-field col s12${f.half ? " m6" : ""}`}>
                {input(f.key, `${f.label}${f.required ? " *" : ""}`, f)}
              </div>
            ),
          )}
          <div className="col s12">
            <button className="btn teal" onClick={save} disabled={busy}>Save</button>{" "}
            <button className="btn-flat" onClick={() => setEditing(null)}>Cancel</button>
          </div>
        </div>
      ) : (
        <button className="btn-flat teal-text" onClick={() => setEditing(empty())}>
          <i className="material-icons left">add</i>Add
        </button>
      )}
    </>
  );

  if (p.flat) return <>{p.title && <h6>{p.title}</h6>}{body}</>;
  return (
    <div className="card">
      <div className="card-content">
        <span className="card-title">{p.title}</span>
        {body}
      </div>
    </div>
  );
}
