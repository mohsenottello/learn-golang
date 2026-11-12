import { useEffect, useRef } from "react";
import { Modal } from "bootstrap";

/** Initialises a Bootstrap modal on the ref'd element; returns open/close. */
export function useModal(onClose?: () => void) {
  const ref = useRef<HTMLDivElement>(null);
  const inst = useRef<Modal | null>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    inst.current = new Modal(el);
    const handler = () => onClose?.();
    el.addEventListener("hidden.bs.modal", handler);
    return () => { el.removeEventListener("hidden.bs.modal", handler); inst.current?.dispose(); };
  }, [onClose]);

  return { ref, open: () => inst.current?.show(), close: () => inst.current?.hide() };
}
