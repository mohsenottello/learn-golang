import { useLang } from "../lang";
import { useSite } from "../site";

export default function Contact() {
  const { site } = useSite();
  const { lang, t } = useLang();
  const text = (lang === "en" && site.contact_text_en) || site.contact_text;

  return (
    <>
      <div className="card-content">
        <span className="card-title">{t("contact_title")}</span>
        {text && <p style={{ whiteSpace: "pre-line" }}>{text}</p>}
      </div>
      <div className="card-action">
        {site.contact_email && (
          <a className="btn teal" href={`mailto:${site.contact_email}`}>
            <i className="material-icons left">mail</i>{site.contact_email}
          </a>
        )}
        {!text && !site.contact_email && <p className="grey-text">{t("contact_empty")}</p>}
      </div>
    </>
  );
}
