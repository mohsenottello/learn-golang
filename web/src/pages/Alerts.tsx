import { FormEvent, useEffect, useState } from "react";
import { Alerts as AlertsT, alerts as api } from "../api";
import { useToast } from "../toast";

export default function Alerts() {
  const [a, setA] = useState<AlertsT | null>(null);
  const [symbol, setSymbol] = useState("BTC");
  const [pct, setPct] = useState("5");
  const toast = useToast();
  const load = () => api.get().then(setA).catch((e) => toast(e.message, "error"));

  useEffect(() => { load(); }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // While not linked, poll so the page flips to "connected" right after /start.
  useEffect(() => {
    if (!a || a.telegram_linked) return;
    const t = setInterval(load, 5000);
    return () => clearInterval(t);
  }, [a?.telegram_linked]); // eslint-disable-line react-hooks/exhaustive-deps

  async function add(e: FormEvent) {
    e.preventDefault();
    try { await api.addWatch(symbol, Number(pct)); toast(`Watching ${symbol.toUpperCase()}.`); load(); }
    catch (err) { toast((err as Error).message, "error"); }
  }
  async function remove(id: number) {
    try { await api.deleteWatch(id); load(); } catch (err) { toast((err as Error).message, "error"); }
  }
  async function unlink() {
    if (!confirm("Disconnect Telegram? You'll stop receiving alerts.")) return;
    try { await api.unlink(); load(); } catch (err) { toast((err as Error).message, "error"); }
  }

  if (!a) return <p>Loading…</p>;
  if (!a.enabled) return <div className="alert alert-secondary">Price alerts aren't configured on this server.</div>;

  const money = (p: number) => p >= 1000 ? `$${p.toFixed(0)}` : p >= 1 ? `$${p.toFixed(2)}` : `$${p.toPrecision(4)}`;
  const start = `/start ${a.link_code}`;

  return (
    <>
      <div className="card mb-3"><div className="card-body">
        <h5 className="card-title">Telegram</h5>
        {a.telegram_linked ? (
          <p className="mb-2">Connected <i className="material-icons align-middle text-success">check_circle</i>{" "}
            <button className="btn btn-link btn-sm" onClick={unlink}>Disconnect</button></p>
        ) : (
          <>
            <p>Open <a href={`https://t.me/${a.bot_username}`} target="_blank" rel="noreferrer">@{a.bot_username}</a> and send this message:</p>
            <div className="input-group" style={{ maxWidth: 360 }}>
              <input className="form-control font-monospace" readOnly value={start} />
              <button className="btn btn-outline-secondary" onClick={() => navigator.clipboard.writeText(start)}>Copy</button>
            </div>
            <p className="text-secondary small mt-2">This page updates automatically once the bot has linked your chat.</p>
          </>
        )}
      </div></div>

      <div className="card mb-3"><div className="card-body">
        <h5 className="card-title">Price alerts</h5>
        <form className="row g-2 align-items-end mb-3" onSubmit={add}>
          <div className="col-5 col-md-3">
            <label className="form-label small">Coin</label>
            <input className="form-control" value={symbol} onChange={(e) => setSymbol(e.target.value.toUpperCase())} placeholder="BTC" />
          </div>
          <div className="col-4 col-md-2">
            <label className="form-label small">Move (%)</label>
            <input className="form-control" type="number" step="0.1" min="0.1" max="100" value={pct} onChange={(e) => setPct(e.target.value)} />
          </div>
          <div className="col-3 col-md-2"><button className="btn btn-primary w-100" disabled={!symbol || !pct}>Add</button></div>
        </form>
        {a.watches.length === 0 ? <p className="text-secondary">No alerts yet.</p> : (
          <table className="table table-hover align-middle">
            <thead><tr><th>Coin</th><th>Alert on</th><th>Reference</th><th>Since</th><th></th></tr></thead>
            <tbody>
              {a.watches.map((w) => (
                <tr key={w.id}>
                  <td><strong>{w.symbol}</strong></td>
                  <td>±{w.pct}%</td>
                  <td>{money(w.ref_price)}</td>
                  <td className="text-secondary small">{new Date(w.ref_time).toLocaleString()}</td>
                  <td className="text-end"><button className="btn btn-link btn-sm p-1 text-danger" onClick={() => remove(w.id)}><i className="material-icons">delete</i></button></td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <p className="text-secondary small mb-0">An alert fires when the price moves by the chosen percent from the reference; the reference then resets to the new price.</p>
      </div></div>
    </>
  );
}
