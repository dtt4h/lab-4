import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import api from "../api/api";

function IncidentsPage() {
  const [incidents, setIncidents] = useState([]);
  const [search, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;

  useEffect(() => {
    fetchIncidents();
  }, []);

  async function fetchIncidents() {
    const res = await api.get("/incidents");
    setIncidents(res.data);
  }

  async function deleteIncident(id) {
    if (window.confirm("Удалить этот инцидент?")) {
      setIncidents((prev) => prev.filter((item) => item.id !== id));
      
      try {
        await api.delete(`/incidents/${id}`);
      } catch (err) {
        console.error("Ошибка при удалении:", err);
        setIncidents((prev) => [...prev, { id }]);
        alert("Не удалось удалить инцидент");
      }
    }
  }

  const filtered = incidents.filter(
    (item) =>
      item.title.toLowerCase().includes(search.toLowerCase()) ||
      item.hotel.toLowerCase().includes(search.toLowerCase()) ||
      item.status.toLowerCase().includes(search.toLowerCase())
  ).sort((a, b) => new Date(b.date) - new Date(a.date));

  const totalPages = Math.ceil(filtered.length / itemsPerPage);
  const start = (currentPage - 1) * itemsPerPage;
  const paginated = filtered.slice(start, start + itemsPerPage);

  return (
    <div>
      <div className="page-header">
        <h1>Инциденты</h1>

        <Link to="/incidents/add" className="add-btn">
          + Новый инцидент
        </Link>
      </div>

      <div className="search-box">
        <input
          type="text"
          placeholder="Поиск по заголовку, гостинице или статусу..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setCurrentPage(1);
          }}
        />
      </div>

      <div className="table-container">
        {paginated.length === 0 ? (
          <div className="empty-state">
            <p>Инцидентов пока нет</p>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Заголовок</th>
                <th>Дата</th>
                <th>Этаж</th>
                <th>Уровень</th>
                <th>Статус</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {paginated.map((item) => (
                <tr key={item.id}>
                  <td>{item.title}</td>
                  <td>{new Date(item.date).toLocaleString("ru-RU")}</td>
                  <td>{item.floor}</td>
                  <td>{item.level}</td>
                  <td>{item.status}</td>
                  <td>
                  <Link to={`/incidents/${item.id}`}>Детали</Link>
                  <Link to={`/incidents/edit/${item.id}`}>Ред.</Link>
                  <a href="#" onClick={(e) => { e.preventDefault(); deleteIncident(item.id); }}>Удал.</a>
                </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {totalPages > 1 && (
        <div className="pagination">
          <button
            disabled={currentPage === 1}
            onClick={() => setCurrentPage(currentPage - 1)}
          >
            Назад
          </button>
          <span>
            Стр. {currentPage} из {totalPages}
          </span>
          <button
            disabled={currentPage === totalPages}
            onClick={() => setCurrentPage(currentPage + 1)}
          >
            Вперёд
          </button>
        </div>
      )}
    </div>
  );
}

export default IncidentsPage;
