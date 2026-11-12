import { createContext, ReactNode, useContext } from "react";
import { Toast } from "bootstrap";

type Kind = "success" | "error";

/** Shows a Bootstrap toast in a fixed bottom-right stack. No React state needed. */
function show(text: string, kind: Kind = "success") {
  let stack = document.getElementById("toast-stack");
  if (!stack) {
    stack = document.createElement("div");
    stack.id = "toast-stack";
    stack.className = "toast-container position-fixed bottom-0 end-0 p-3";
    document.body.appendChild(stack);
  }
  const el = document.createElement("div");
  el.className = `toast align-items-center text-bg-${kind === "error" ? "danger" : "success"} border-0`;
  el.setAttribute("role", "alert");
  el.innerHTML = `<div class="d-flex"><div class="toast-body"></div>
    <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button></div>`;
  el.querySelector(".toast-body")!.textContent = text; // textContent: never inject HTML
  stack.appendChild(el);
  el.addEventListener("hidden.bs.toast", () => el.remove());
  new Toast(el, { delay: 3000 }).show();
}

const Ctx = createContext(show);
export function ToastProvider({ children }: { children: ReactNode }) {
  return <Ctx.Provider value={show}>{children}</Ctx.Provider>;
}
export const useToast = () => useContext(Ctx);
