import { useSite } from "../site";

export default function Contact() {
  const { site } = useSite();
  return (
    <>
      <span className="card-title">Contact</span>
      {site.contact_text && <p style={{ whiteSpace: "pre-line" }}>{site.contact_text}</p>}
      {site.contact_email && (
        <a className="btn teal" href={`mailto:${site.contact_email}`}>
          <i className="material-icons left">mail</i>{site.contact_email}
        </a>
      )}
      {!site.contact_text && !site.contact_email && <p className="grey-text">No contact details yet.</p>}
    </>
  );
}
