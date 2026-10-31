import { useState } from "react";
import { Lang, translate } from "../api";
import { useSite } from "../site";
import { useToast } from "../toast";

type Props = {
  from: Lang; to: Lang;
  texts: string[];                       // source values, in order
  onResult: (translated: string[]) => void;
  small?: boolean;
};

export default function TranslateButton({ from, to, texts, onResult, small }: Props) {
  const { site } = useSite();
  const toast = useToast();
  const [busy, setBusy] = useState(false);
  if (!site.translation_enabled) return null;

  async function run() {
    const nonEmpty = texts.map((t, i) => [t, i] as const).filter(([t]) => t.trim());
    if (nonEmpty.length === 0) { toast(`Nothing to translate from ${from.toUpperCase()}.`, "error"); return; }
    setBusy(true);
    try {
      const { texts: out } = await translate(nonEmpty.map(([t]) => t), from, to);
      const result = [...texts];
      nonEmpty.forEach(([, i], k) => { result[i] = out[k]; });
      onResult(result);
    } catch (e) { toast((e as Error).message, "error"); }
    finally { setBusy(false); }
  }

  return (
    <button type="button" className={`btn-flat teal-text${small ? " btn-small" : ""}`} onClick={run} disabled={busy}
      title={`Translate ${from.toUpperCase()} → ${to.toUpperCase()} with DeepL`}>
      <i className="material-icons left">translate</i>{from.toUpperCase()}→{to.toUpperCase()}
    </button>
  );
}
