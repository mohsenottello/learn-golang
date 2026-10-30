import { Link } from "react-router-dom";
import { useSite } from "../site";
import { useAuth } from "../auth";

export default function Home() {
  const { site, logoVersion } = useSite();
  const registration = site.registration;
  const { user } = useAuth();
  const h = site.home;
  const resumeHref = `/${h.featured_slug}`;
  const createHref = user ? "/account" : "/register";
  const initials = h.headline.split(/\s+/).filter(Boolean).map((w) => w[0]).join("").slice(0, 2).toUpperCase();

  return (
    <>
      {/* Hero */}
      <section className="hero">
        <div className="container">
          <div className="row" style={{ marginBottom: 0 }}>
            <div className="col s12 m6 hero-copy">
              <p className="script-text">{h.greeting}</p>
              <h1>{h.headline}</h1>
              <p className="hero-description">{h.description}</p>
              <div className="hero-buttons">
                <Link to={resumeHref} className="btn btn-main">
                  {h.primary_label}<i className="material-icons right">arrow_forward</i>
                </Link>
                {h.secondary_label && (
                  <Link to={createHref} className="btn-outline">
                    {h.secondary_label}
                  </Link>
                )}
              </div>
              {h.note && <p className="hero-note"><i className="material-icons">check_circle_outline</i> {h.note}</p>}
            </div>

            <div className="col s12 m6">
              <div className="portrait-wrap">
                <div className="portrait-shape"></div>
                {site.has_portrait ? (
                  <img src={`/api/site/portrait?v=${logoVersion}`} alt={h.headline} className="portrait" />
                ) : (
                  <div className="portrait" style={{ display: "flex", alignItems: "center", justifyContent: "center",
                    fontSize: 80, fontWeight: 800, color: "var(--primary)" }}>{initials}</div>
                )}
                {h.floating_note && (
                  <div className="floating-note"><i className="material-icons">auto_awesome</i>{h.floating_note}</div>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      {registration && (
        <section className="section-padding features-section" id="features">
          <div className="container center">
            <div className="section-label">{h.features_label}</div>
            <h2 className="section-heading">{h.features_heading}</h2>
            <p className="section-intro">{h.features_intro}</p>
            <div className="row">
              {h.features.map((f, i) => (
                <div key={i} className="col s12 m6 l3">
                  <article className="feature-card">
                    <div className={`feature-icon icon-${f.color}`}><i className="material-icons">{f.icon}</i></div>
                    <h5>{f.title}</h5>
                    <p>{f.text}</p>
                  </article>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* About the resume */}
      <section className="section-padding resume-section" id="about-resume">
        <div className="container">
          <div className="row" style={{ display: "flex", flexWrap: "wrap", marginBottom: 0 }}>
            <div className="col s12 m6">
              <div className="resume-illustration">
                <div className="resume-paper">
                  <div className="resume-person">
                    <div className="resume-avatar"><i className="material-icons">person</i></div>
                    <div><h5>{h.headline}</h5><p>Professional Profile</p></div>
                  </div>
                  <div className="resume-line-title">EXPERIENCE</div>
                  <div className="resume-line long"></div><div className="resume-line medium"></div><div className="resume-line short"></div>
                  <div className="resume-line-title">EDUCATION</div>
                  <div className="resume-line long"></div><div className="resume-line medium"></div>
                  <div className="resume-line-title">SKILLS</div>
                  <div className="resume-skills"><span>Skills</span><span>Experience</span><span>Education</span></div>
                </div>
              </div>
            </div>
            <div className="col s12 m6 resume-copy">
              <div className="section-label">{h.about_label}</div>
              <h2 className="section-heading">{h.about_heading}</h2>
              <p>{h.about_text}</p>
              <Link to={resumeHref} className="btn btn-main">
                {h.about_button}<i className="material-icons right">arrow_forward</i>
              </Link>

              <Link to={resumeHref} className="resume-url">
                <i className="material-icons">link</i>{window.location.host}{resumeHref}
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Call to action */}
      {registration && (
         <section className="cta-section center">
        <div className="container">
          <i className="material-icons cta-icon">description</i>
          <h2>{h.cta_heading}</h2>
          <p>{h.cta_text}</p>
          <Link to={createHref} className="btn btn-main" style={{ background: "#fff", color: "var(--primary)" }}>
            {h.cta_button}<i className="material-icons right">arrow_forward</i>
          </Link>
        </div>
      </section>
      )}

      {/* Footer */}
      <footer>
        <div className="container">
          <div className="row" style={{ marginBottom: 0 }}>
            <div className="col s12 m5">
              <Link to="/" className="footer-brand"><i className="material-icons">spa</i>{site.name}</Link>
              <p>{h.footer_text}</p>
            </div>
            <div className="col s12 m7">
              <div className="footer-links">
                <Link to="/">Home</Link>
                <Link to={resumeHref}>My Resume</Link>
                {registration && (
                  <Link to={createHref}>Create Resume</Link>
                )}
                <Link to="/contact">Contact</Link>
              </div>
            </div>
          </div>
          <div className="footer-bottom">
            <p>© {new Date().getFullYear()} {site.name}. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </>
  );
}
