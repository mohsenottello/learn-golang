import { createContext, ReactNode, useContext } from "react";

type Kind = "success" | "error";
const show = (text: string, kind: Kind = "success") =>
  M.toast({ html: text, classes: kind === "error" ? "red darken-1" : "teal", displayLength: 3000 });

const Ctx = createContext(show);
export function ToastProvider({ children }: { children: ReactNode }) {
  return <Ctx.Provider value={show}>{children}</Ctx.Provider>;
}
export const useToast = () => useContext(Ctx);
