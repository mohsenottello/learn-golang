import { Lang } from "./api";

// Section labels on the public page, per language.
export const LABELS: Record<Lang, {
  contact: string; about: string; software: string; languages: string; education: string;
  experience: string; present: string; link: string; locale: string;
}> = {
  de: {
    contact: "Kontaktdaten", about: "Über mich", software: "Software", languages: "Sprachen",
    education: "Bildung", experience: "Berufserfahrung", present: "heute", link: "Link", locale: "de-DE",
  },
  en: {
    contact: "Contact", about: "About me", software: "Software", languages: "Languages",
    education: "Education", experience: "Experience", present: "present", link: "Link", locale: "en-GB",
  },
};
