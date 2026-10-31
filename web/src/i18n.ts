import { Lang } from "./api";

// Every UI string the public site shows, in both languages. Keys are shared,
// so a key missing from one language is a compile error.
const de = {
  // nav
  nav_home: "Start", nav_resume: "Mein Lebenslauf", nav_editor: "Mein Editor", nav_contact: "Kontakt",
  nav_login: "Anmelden", nav_register: "Loslegen", nav_logout: "Abmelden", nav_admin: "Admin",
  // auth pages
  login_title: "Anmelden", login_button: "Anmelden", forgot_link: "Passwort vergessen?", need_account: "Noch kein Konto? Registrieren",
  register_title: "Konto erstellen", register_button: "Registrieren", have_account: "Schon ein Konto? Anmelden",
  register_closed_title: "Registrierung geschlossen", register_closed_text: "Die Registrierung neuer Nutzer ist derzeit deaktiviert.",
  check_email: "Bitte E-Mail prüfen",
  forgot_title: "Passwort zurücksetzen", forgot_text: "Gib deine E-Mail ein und wir senden dir einen Link für ein neues Passwort.",
  forgot_button: "Link senden", back_to_login: "Zurück zur Anmeldung",
  reset_title: "Neues Passwort wählen", reset_new: "Neues Passwort", reset_confirm: "Passwort bestätigen", reset_save: "Speichern",
  reset_done: "Passwort aktualisiert", reset_missing: "Diesem Link fehlt das Token.", reset_mismatch: "Die Passwörter stimmen nicht überein.",
  verify_wait: "Wird geprüft…", verify_ok: "Bestätigt", verify_bad: "Ungültiger Link",
  invalid_link: "Ungültiger Link",
  // fields
  f_name: "Name", f_email: "E-Mail", f_password: "Passwort",
  // contact
  contact_title: "Kontakt", contact_empty: "Noch keine Kontaktdaten.",
  // home extras
  home_professional_profile: "Berufliches Profil", home_create: "Lebenslauf erstellen", rights: "Alle Rechte vorbehalten.",
  // editor sidebar
  sec_profile: "Profil", sec_software: "Software", sec_languages: "Sprachen", sec_education: "Bildung", sec_experience: "Berufserfahrung", sec_account: "Konto",
  // resume page
  cv_contact: "Kontaktdaten", cv_about: "Über mich", cv_software: "Software", cv_languages: "Sprachen", cv_education: "Bildung",
  cv_experience: "Berufserfahrung", cv_present: "heute", cv_link: "Link", cv_pdf: "PDF", cv_notfound: "Kein veröffentlichter Lebenslauf unter",
  locale: "de-DE",
};

const en: typeof de = {
  nav_home: "Home", nav_resume: "My Resume", nav_editor: "My editor", nav_contact: "Contact",
  nav_login: "Login", nav_register: "Get Started", nav_logout: "Log out", nav_admin: "Admin",
  login_title: "Log in", login_button: "Log in", forgot_link: "Forgot password?", need_account: "Need an account? Register",
  register_title: "Create account", register_button: "Register", have_account: "Already have an account? Log in",
  register_closed_title: "Registration closed", register_closed_text: "New user registration is currently disabled.",
  check_email: "Check your email",
  forgot_title: "Reset your password", forgot_text: "Enter your email and we'll send you a link to choose a new password.",
  forgot_button: "Send link", back_to_login: "Back to login",
  reset_title: "Choose a new password", reset_new: "New password", reset_confirm: "Confirm password", reset_save: "Save",
  reset_done: "Password updated", reset_missing: "This reset link is missing its token.", reset_mismatch: "Passwords don't match.",
  verify_wait: "Verifying…", verify_ok: "Verified", verify_bad: "Invalid link",
  invalid_link: "Invalid link",
  f_name: "Name", f_email: "Email", f_password: "Password",
  contact_title: "Contact", contact_empty: "No contact details yet.",
  home_professional_profile: "Professional Profile", home_create: "Create Resume", rights: "All rights reserved.",
  sec_profile: "Profile", sec_software: "Software", sec_languages: "Languages", sec_education: "Education", sec_experience: "Experience", sec_account: "Account",
  cv_contact: "Contact", cv_about: "About me", cv_software: "Software", cv_languages: "Languages", cv_education: "Education",
  cv_experience: "Experience", cv_present: "present", cv_link: "Link", cv_pdf: "PDF", cv_notfound: "No published resume at",
  locale: "en-GB",
};

export type StringKey = keyof typeof de;
export const STRINGS: Record<Lang, typeof de> = { de, en };
