export type Admin = { id: number; email: string };

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
  }
}

async function request<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(path, {
    method: body === undefined ? "GET" : "POST",
    headers: body === undefined ? {} : { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
    credentials: "same-origin",
  });
  if (res.status === 204) return undefined as T;
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new ApiError(res.status, data.error ?? "Request failed");
  return data as T;
}

export const api = {
  register: (email: string, password: string) =>
    request<{ message: string }>("/api/admin/register", { email, password }),
  verify: (token: string) =>
    request<{ message: string }>("/api/admin/verify", { token }),
  login: (email: string, password: string) =>
    request<Admin>("/api/admin/login", { email, password }),
  logout: () => request<void>("/api/admin/logout", {}),
  me: () => request<Admin>("/api/admin/me"),
};
