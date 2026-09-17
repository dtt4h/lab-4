import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import api, { can } from "../api/api";
import { bookingStatusLabel } from "../constants/bookingStatus";
import { incidentStatusLabel } from "../constants/incidentStatus";

const dateTime = (value) => value ? new Date(value).toLocaleString("ru-RU") : "—";

export default function DashboardPage() {
  const incidentsAllowed = can("incident.read");
  const bookingsAllowed = can("booking.read");
  const [state, setState] = useState({ incidents: [], bookings: [], loading: true, error: "" });

  const load = () => {
    setState((current) => ({ ...current, loading: true, error: "" }));
    Promise.all([
      incidentsAllowed ? api.get("/incidents?limit=100") : Promise.resolve({ data: { data: [] } }),
      bookingsAllowed ? api.get("/admin/bookings") : Promise.resolve({ data: { data: [] } }),
    ]).then(([incidents, bookings]) => setState({ incidents: incidents.data.data || [], bookings: bookings.data.data || [], loading: false, error: "" }))
      .catch(() => setState((current) => ({ ...current, loading: false, error: "Не удалось загрузить данные дашборда." })));
  };

  useEffect(() => { load(); }, [incidentsAllowed, bookingsAllowed]);
  if (state.loading) return <div className="empty-state">Загружаем данные дашборда…</div>;
  if (state.error) return <div><p className="form-error" role="alert">{state.error}</p><button className="secondary-btn" onClick={load}>Повторить</button></div>;

  const activeIncidents = state.incidents.filter((item) => ["open", "in_progress"].includes(item.status)).length;
  const pendingBookings = state.bookings.filter((item) => item.status === "pending").length;
  const recentIncidents = [...state.incidents].sort((a, b) => new Date(b.occurred_at) - new Date(a.occurred_at)).slice(0, 4);
  const recentBookings = [...state.bookings].sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(0, 4);

  return <div className="dashboard"><div className="page-header"><div><span className="eyebrow">ОБЗОР</span><h1>Дашборд</h1><p className="muted">Главное за текущий момент.</p></div><button className="secondary-btn" onClick={load}>Обновить</button></div>
    <section className="stats-grid">{incidentsAllowed && <><article><span>Активные инциденты</span><b>{activeIncidents}</b><small>из {state.incidents.length} всего</small></article><article><span>Новые инциденты</span><b>{state.incidents.filter((item) => item.status === "open").length}</b><small>требуют реакции</small></article></>}{bookingsAllowed && <><article><span>Ожидают оплаты</span><b>{pendingBookings}</b><small>из {state.bookings.length} бронирований</small></article><article><span>Подтверждённые брони</span><b>{state.bookings.filter((item) => item.status === "confirmed").length}</b><small>актуальные заезды</small></article></>}</section>
    <div className="dashboard-grid">{incidentsAllowed && <section className="dashboard-panel"><div className="panel-heading"><h2>Последние инциденты</h2><Link to="/incidents">Все</Link></div>{recentIncidents.length ? recentIncidents.map((item) => <Link className="activity-row" key={item.id} to={`/incidents/${item.id}`}><div><strong>{item.title}</strong><small>{dateTime(item.occurred_at)}</small></div><span className={`status-badge incident-status-${item.status}`}>{incidentStatusLabel(item.status)}</span></Link>) : <p className="muted">Инцидентов пока нет.</p>}</section>}{bookingsAllowed && <section className="dashboard-panel"><div className="panel-heading"><h2>Последние бронирования</h2><Link to="/bookings">Все</Link></div>{recentBookings.length ? recentBookings.map((item) => <Link className="activity-row" key={item.id} to="/bookings"><div><strong>{item.guest_name || "Гость"}</strong><small>{item.check_in} — {item.check_out}</small></div><span className={`status-badge status-${item.status}`}>{bookingStatusLabel(item.status)}</span></Link>) : <p className="muted">Бронирований пока нет.</p>}</section>}</div>
  </div>;
}
