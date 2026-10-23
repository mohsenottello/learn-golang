import { useEffect, useRef } from "react";

/** Initialises a Materialize modal on the ref'd element; returns open/close. */
export function useModal(onClose?: () => void) {
  const ref = useRef<HTMLDivElement>(null);
  const inst = useRef<M.Modal | null>(null);

  useEffect(() => {
    if (!ref.current) return;
    inst.current = M.Modal.init(ref.current, { onCloseEnd: onClose });
    return () => inst.current?.destroy();
  }, [onClose]);

  return { ref, open: () => inst.current?.open(), close: () => inst.current?.close() };
}
