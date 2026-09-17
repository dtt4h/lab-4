import { NavLink, useNavigate } from "react-router-dom";
import api, { can, setAccessToken, setCurrentUser } from "../api/api";

export default function Sidebar() {
  const navigate = useNavigate();
  const logout = async () => {
    try { await api.post("/auth/logout"); } finally {
      setAccessToken(null);
      setCurrentUser(null);
      navigate("/login");
    }
  };

  return <aside className="sidebar"><div className="logo"><h2>Админ-панель<br />отелей</h2></div>
    <nav className="menu" aria-label="Основная навигация">
      {(can("hotel.read") || can("booking.read") || can("incident.read")) && <NavLink to="/dashboard">Дашборд</NavLink>}
      {can("booking.read") && <NavLink to="/bookings">Бронирования</NavLink>}
      {can("room.read") && <NavLink to="/rooms">Номера</NavLink>}
      {can("incident.read") && <NavLink to="/incidents">Инциденты</NavLink>}
      {can("catalog.manage") && <NavLink to="/catalogs">Справочники</NavLink>}
      {can("roles.manage") && <NavLink to="/access">Сотрудники и роли</NavLink>}
    </nav>
    <button className="link-danger" onClick={logout}>Выйти</button>
  </aside>;
}
