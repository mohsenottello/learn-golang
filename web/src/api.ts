export type User = { id: number; email: string; name: string; verified: boolean; created_at: string };
export type Admin = { id: number; email: string; created_at: string };
type Msg = { message: string };

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
  }
}

async function request<T>(path: string, method: "GET" | "POST" | "DELETE", body?: unknown): Promise<T> {
  const res = await fetch(path, {
    method,
    headers: body === undefined ? {} : { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
    credentials: "same-origin",
  });
  if (res.status === 204) return undefined as T;
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new ApiError(res.status, data.error ?? "Request failed");
  return data as T;
}

const get = <T,>(p: string) => request<T>(p, "GET");
const post = <T,>(p: string, b: unknown) => request<T>(p, "POST", b);
const del = <T,>(p: string) => request<T>(p, "DELETE");

export const auth = {
  register: (email: string, name: string, password: string) => post<Msg>("/api/auth/register", { email, name, password }),
  verify: (token: string) => post<Msg>("/api/auth/verify", { token }),
  login: (email: string, password: string) => post<User>("/api/auth/login", { email, password }),
  logout: () => post<void>("/api/auth/logout", {}),
  me: () => get<User>("/api/auth/me"),
  forgotPassword: (email: string) => post<Msg>("/api/auth/forgot-password", { email }),
  resetPassword: (token: string, password: string) => post<Msg>("/api/auth/reset-password", { token, password }),
};

export const admin = {
  login: (email: string, password: string) => post<Admin>("/api/admin/login", { email, password }),
  logout: () => post<void>("/api/admin/logout", {}),
  me: () => get<Admin>("/api/admin/me"),
  listUsers: () => get<User[]>("/api/admin/users"),
  deleteUser: (id: number) => del<void>(`/api/admin/users/${id}`),
};
