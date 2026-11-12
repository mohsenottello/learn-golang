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
    <div className="card mb-3">
      <div className="card-body">
        <h5 className="card-title">Resumes</h5>
        <table className="table table-striped align-middle">
          <thead><tr><th>Owner</th><th>Slug</th><th>Status</th><th>Updated</th><th></th></tr></thead>
          <tbody>
            {list.map((r) => (
              <tr key={r.user_id}>
                <td>{r.name}<br /><small className="text-secondary">{r.email}</small></td>
                <td>{r.published ? <Link to={`/${r.slug}`} target="_blank">/{r.slug}</Link> : `/${r.slug}`}</td>
                <td>{r.published ? <span className="text-primary">published</span> : <span className="text-secondary">draft</span>}</td>
                <td>{new Date(r.updated_at).toLocaleDateString()}</td>
                <td>
                  <button className={`btn btn-link ${r.published ? "text-danger" : "text-primary"}`} onClick={() => toggle(r)}>
                    {r.published ? "Unpublish" : "Publish"}
                  </button>
                </td>
              </tr>
            ))}
            {list.length === 0 && <tr><td colSpan={5} className="text-secondary">No resumes yet.</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}
