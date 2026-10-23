import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { AdminResume, admin } from "../api";
import { useToast } from "../toast";

export default function AdminResumes() {
  const [list, setList] = useState<AdminResume[]>([]);
  const toast = useToast();

  const load = () => admin.listResumes().then(setList).catch((e) => toast(e.message, "error"));
  useEffect(() => { load(); }, []); // eslint-disable-line react-hooks/exhaustive-deps

  async function toggle(r: AdminResume) {
    try {
      await admin.setPublished(r.user_id, !r.published);
      toast(r.published ? "Unpublished." : "Published.");
      load();
    } catch (e) { toast((e as Error).message, "error"); }
  }

  return (
    <div className="card">
      <div className="card-content">
        <span className="card-title">Resumes</span>
        <table className="striped">
          <thead><tr><th>Owner</th><th>Slug</th><th>Status</th><th>Updated</th><th></th></tr></thead>
          <tbody>
            {list.map((r) => (
              <tr key={r.user_id}>
                <td>{r.name}<br /><small className="grey-text">{r.email}</small></td>
                <td>{r.published ? <Link to={`/${r.slug}`} target="_blank">/{r.slug}</Link> : `/${r.slug}`}</td>
                <td>{r.published ? <span className="teal-text">published</span> : <span className="grey-text">draft</span>}</td>
                <td>{new Date(r.updated_at).toLocaleDateString()}</td>
                <td>
                  <button className={`btn-flat ${r.published ? "red-text" : "teal-text"}`} onClick={() => toggle(r)}>
                    {r.published ? "Unpublish" : "Publish"}
                  </button>
                </td>
              </tr>
            ))}
            {list.length === 0 && <tr><td colSpan={5} className="grey-text">No resumes yet.</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}
