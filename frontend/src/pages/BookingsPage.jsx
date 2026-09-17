import { useEffect, useMemo, useState } from "react";
import api from "../api/api";
import { BOOKING_STATUSES, bookingStatusActions, bookingStatusLabel } from "../constants/bookingStatus";

export default function BookingsPage() {
  const [items, setItems] = useState([]);
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState("all");
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState("");
  const [error, setError] = useState("");

  const load = async () => {
    setLoading(true);
    try {
      const response = await api.get("/admin/bookings");
      setItems(response.data.data || []);
      setError("");
    } catch (err) {
      setError(err.response?.status === 403 ? "У вас нет доступа к управлению бронированиями." : "Не удалось загрузить бронирования. Попробуйте обновить страницу.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const visibleItems = useMemo(() => items.filter((item) => {
    const searchable = `${item.guest_name} ${item.guest_email} ${item.id}`.toLowerCase();
    return (status === "all" || item.status === status) && searchable.includes(query.trim().toLowerCase());
  }), [items, query, status]);

  async function changeStatus(id, nextStatus) {
    const action = nextStatus === "cancelled" ? "отменить" : "изменить статус";
    if (nextStatus === "cancelled" && !window.confirm("Отменить это бронирование?")) return;
    setUpdating(id);
    try {
      await api.patch(`/admin/bookings/${id}/status`, { status: nextStatus });
      await load();
    } catch (err) {
      setError(err.response?.status === 400 ? "Этот переход статуса недоступен." : `Не удалось ${action} бронирование.`);
    } finally {
      setUpdating("");
    }
  }

  const clearFilters = () => { setQuery(""); setStatus("all"); };

  return <div>
    <div className="page-header">
      <div><span className="eyebrow">АДМИН-ПАНЕЛЬ</span><h1>Бронирования</h1><p className="muted">Поиск, обработка и контроль статусов брони.</p></div>
      <button className="secondary-btn" onClick={load} disabled={loading}>Обновить</button>
    </div>
    <section className="booking-toolbar" aria-label="Фильтры бронирований">
      <label className="sr-only" htmlFor="booking-search">Поиск бронирования</label>
      <input id="booking-search" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Имя, email или ID брони" />
      <label className="sr-only" htmlFor="booking-status">Статус бронирования</label>
      <select id="booking-status" value={status} onChange={(event) => setStatus(event.target.value)}>
        <option value="all">Все статусы</option>
        {Object.entries(BOOKING_STATUSES).map(([value, item]) => <option key={value} value={value}>{item.label}</option>)}
      </select>
      {(query || status !== "all") && <button className="text-btn" onClick={clearFilters}>Сбросить</button>}
    </section>
    {error && <p className="form-error" role="alert">{error}</p>}
    <div className="table-container">
      {loading ? <div className="empty-state">Загружаем бронирования…</div> : <table>
        <thead><tr><th>Гость</th><th>Даты</th><th>Гостей</th><th>Сумма</th><th>Статус</th><th><span className="sr-only">Действия</span></th></tr></thead>
        <tbody>{visibleItems.map((item) => <tr key={item.id}>
          <td><strong>{item.guest_name || "Гость"}</strong><br /><small>{item.guest_email || "—"}</small></td>
          <td>{item.check_in}<br /><small>до {item.check_out}</small></td>
          <td>{item.guests_count}</td><td>{Number(item.total_price).toLocaleString("ru-RU")} ₽</td>
          <td><span className={`status-badge ${BOOKING_STATUSES[item.status]?.className || ""}`}>{bookingStatusLabel(item.status)}</span></td>
          <td>{bookingStatusActions(item.status).length ? <select aria-label={`Действие для бронирования ${item.id}`} value="" disabled={updating === item.id} onChange={(event) => { if (event.target.value) changeStatus(item.id, event.target.value); }}><option value="">Действия…</option>{bookingStatusActions(item.status).map((action) => <option key={action.value} value={action.value}>{action.label}</option>)}</select> : <span className="muted">Нет действий</span>}</td>
        </tr>)}</tbody>
      </table>}
      {!loading && !visibleItems.length && <div className="empty-state">{items.length ? "По заданным фильтрам ничего не найдено." : "Бронирований пока нет."}</div>}
    </div>
  </div>;
}
