import { NavLink } from "react-router-dom";

function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="logo">
        <h2>Гостиница Космос</h2>
      </div>

      <nav className="menu">
        <NavLink to="/">Инциденты</NavLink>

        <NavLink to="/employees">Сотрудники</NavLink>
      </nav>
    </aside>
  );
}

export default Sidebar;
