import { useLang } from "../lang";
import { useSite } from "../site";

export default function Contact() {
  const { site } = useSite();
  const { lang, t } = useLang();
  const text = (lang === "en" && site.contact_text_en) || site.contact_text;

  return (
    <>
      <h5 className="card-title">{t("contact_title")}</h5>
      {text && <p style={{ whiteSpace: "pre-line" }}>{text}</p>}
      {site.contact_email && (
        <a className="btn btn-primary" href={`mailto:${site.contact_email}`}>
          <i className="material-icons align-middle me-1">mail</i>{site.contact_email}
        </a>
      )}
      {!text && !site.contact_email && <p className="text-secondary">{t("contact_empty")}</p>}
    </>
  );
}
