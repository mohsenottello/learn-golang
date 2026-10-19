export type User = { id: number; email: string; name: string; verified: boolean; created_at: string };
export type Admin = { id: number; email: string; created_at: string };
type Msg = { message: string };

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
  }
}

async function request<T>(path: string, method: "GET" | "POST" | "PUT" | "DELETE", body?: unknown): Promise<T> {
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
  updateName: (name: string) => request<User>("/api/auth/me", "PUT", { name }),
  changePassword: (current_password: string, new_password: string) =>
    request<Msg>("/api/auth/password", "PUT", { current_password, new_password }),
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

// ---------- resume ----------

export type Lang = "de" | "en";

export type Resume = {
  user_id?: number; slug: string; published: boolean; default_lang: Lang;
  title: string; title_en: string; subtitle: string; subtitle_en: string; about: string; about_en: string;
  email: string; phone: string; location: string; linkedin: string; has_photo: boolean;
  side_bg?: string; side_text?: string; main_bg?: string; main_text?: string; accent?: string;
  updated_at?: string;
};
export type Software = { id: number; name: string; position: number };
export type Language = { id: number; name: string; name_en: string; level: string; level_en: string; position: number };
export type Education = {
  id: number; degree: string; degree_en: string; school: string; school_en: string;
  date_from: string; date_to: string; position: number;
};
export type Bullet = { id: number; text: string; text_en: string; position: number };
export type Experience = {
  id: number; company: string; link: string; role: string; role_en: string; subtitle: string; subtitle_en: string;
  date_from: string; date_to: string; location: string; location_en: string; bullets: Bullet[]; position: number;
};
export type FullResume = {
  resume: Resume; name: string; lang?: Lang;
  software: Software[]; languages: Language[]; education: Education[]; experience: Experience[];
};

export type SectionKey = "software" | "languages" | "education" | "experience";
type ItemOf<K extends SectionKey> =
  K extends "software" ? Software : K extends "languages" ? Language : K extends "education" ? Education : Experience;

const put = <T,>(p: string, b: unknown) => request<T>(p, "PUT", b);

export const resume = {
  mine: () => get<FullResume>("/api/me/resume"),
  save: (r: Resume) => put<Resume>("/api/me/resume", r),
  uploadPhoto: async (file: File) => {
    const fd = new FormData();
    fd.append("photo", file);
    const res = await fetch("/api/me/resume/photo", { method: "PUT", body: fd, credentials: "same-origin" });
    if (!res.ok) throw new ApiError(res.status, (await res.json().catch(() => ({}))).error ?? "Upload failed");
  },
  createItem: <K extends SectionKey>(section: K, item: Omit<ItemOf<K>, "id" | "position">) =>
    post<ItemOf<K>>(`/api/me/resume/${section}`, item),
  updateItem: <K extends SectionKey>(section: K, item: ItemOf<K>) =>
    put<ItemOf<K>>(`/api/me/resume/${section}/${item.id}`, item),
  deleteItem: (section: SectionKey, id: number) => del<void>(`/api/me/resume/${section}/${id}`),
  reorder: (section: SectionKey, ids: number[]) => put<void>(`/api/me/resume/order/${section}`, { ids }),

  createBullet: (expId: number, b: Omit<Bullet, "id" | "position">) => post<Bullet>(`/api/me/resume/experience/${expId}/bullets`, b),
  updateBullet: (expId: number, b: Bullet) => put<Bullet>(`/api/me/resume/experience/${expId}/bullets/${b.id}`, b),
  deleteBullet: (expId: number, id: number) => del<void>(`/api/me/resume/experience/${expId}/bullets/${id}`),
  reorderBullets: (expId: number, ids: number[]) => put<void>(`/api/me/resume/experience/${expId}/bullets-order`, { ids }),

  public: (slug: string, lang?: Lang) => get<FullResume>(`/api/resumes/${slug}${lang ? `?lang=${lang}` : ""}`),
};
