import { useEffect, useState } from "react";
import api from "../api/api";

const roles = [
  { value: "guest", label: "Гость" },
  { value: "booking_manager", label: "Менеджер бронирований" },
  { value: "security_manager", label: "Менеджер безопасности" },
  { value: "hotel_manager", label: "Управляющий гостиницей" },
  { value: "room_manager", label: "Менеджер номерного фонда" },
  { value: "super_admin", label: "Главный администратор" },
];

export default function AccessPage() {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState("");

  const load = async () => {
    setLoading(true);
    try { const response = await api.get("/admin/users"); setUsers(response.data.data || []); setError(""); }
    catch { setError("Не удалось загрузить пользователей. Раздел доступен только главному администратору."); }
    finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  async function change(id, role) {
    setUpdating(id);
    try { await api.patch(`/admin/users/${id}/role`, { role }); await load(); }
    catch { setError("Не удалось изменить роль пользователя. Попробуйте ещё раз."); }
    finally { setUpdating(""); }
  }

  return <div><div className="page-header"><div><span className="eyebrow">ДОСТУП</span><h1>Сотрудники и роли</h1><p className="muted">Роль определяет доступ к разделам и действиям в системе.</p></div></div>
    {error && <p className="form-error" role="alert">{error}</p>}
    <div className="table-container">{loading ? <div className="empty-state">Загружаем пользователей…</div> : <table><thead><tr><th>Пользователь</th><th>Email</th><th>Роль</th></tr></thead><tbody>{users.map((user) => <tr key={user.id}><td>{user.username}</td><td>{user.email}</td><td><select value={user.role} disabled={updating === user.id} aria-label={`Роль пользователя ${user.username}`} onChange={(event) => change(user.id, event.target.value)}>{roles.map((role) => <option key={role.value} value={role.value}>{role.label}</option>)}</select></td></tr>)}</tbody></table>}{!loading && !users.length && !error && <div className="empty-state">Пользователей пока нет.</div>}</div>
  </div>;
}
