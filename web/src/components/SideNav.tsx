import { ReactNode, useState } from "react";
import { STRINGS, StringKey } from "../i18n";
import { useLang } from "../lang";

export type NavItem<K extends string> = { key: K; label: string; icon: string };
export type NavGroup<K extends string> = { title?: string; items: NavItem<K>[] };

type Props<K extends string> = {
  header: ReactNode;
  groups: NavGroup<K>[];
  active: K;
  onSelect: (key: K) => void;
  footer?: ReactNode;
};

export default function SideNav<K extends string>({ header, groups, active, onSelect, footer }: Props<K>) {
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem("sidenav.collapsed") === "1");
  const { lang } = useLang();
  const label = (l: string) => (l in STRINGS[lang] ? STRINGS[lang][l as StringKey] : l);
  const [closed, setClosed] = useState<Record<string, boolean>>({});

  const toggleCollapsed = () => {
    const next = !collapsed;
    setCollapsed(next);
    localStorage.setItem("sidenav.collapsed", next ? "1" : "0");
  };

  return (
    <aside className={`editor-side${collapsed ? " collapsed" : ""}`}>
      <div className="editor-side-head">
        {!collapsed && <div className="editor-side-name">{header}</div>}
        <button className="btn-flat editor-collapse" onClick={toggleCollapsed} title={collapsed ? "Expand" : "Collapse"}>
          <i className="material-icons">{collapsed ? "chevron_right" : "chevron_left"}</i>
        </button>
      </div>

      {groups.map((g, i) => {
        const id = g.title ?? `g${i}`;
        const isClosed = !!closed[id] && !collapsed;
        return (
          <div key={id} className="editor-group-wrap">
            {g.title && !collapsed && (
              <button className="editor-group" onClick={() => setClosed({ ...closed, [id]: !closed[id] })}>
                <span>{g.title}</span>
                <i className="material-icons">{isClosed ? "expand_more" : "expand_less"}</i>
              </button>
            )}
            {!isClosed && g.items.map((it) => (
              <a key={it.key} href={`#${it.key}`} title={collapsed ? label(it.label) : undefined}
                onClick={(e) => { e.preventDefault(); onSelect(it.key); }}
                className={`editor-nav${active === it.key ? " active" : ""}`}>
                <i className="material-icons">{it.icon}</i>
                {!collapsed && <span>{label(it.label)}</span>}
              </a>
            ))}
          </div>
        );
      })}

      {footer && !collapsed && <div className="editor-side-foot">{footer}</div>}
    </aside>
  );
}
