import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import api from "../api/api";

function EmployeesPage() {
  const [employees, setEmployees] = useState([]);
  const [searchText, setSearch] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;

  useEffect(() => {
    fetchEmployees();
  }, []);

  async function fetchEmployees() {
    const res = await api.get("/employees");
    setEmployees(res.data);
  }

  async function deleteEmployee(id) {
    if (window.confirm("Удалить этого сотрудника?")) {
      await api.delete(`/employees/${id}`);
      fetchEmployees();
    }
  }

  const filtered = employees.filter((item) =>
    item.name.toLowerCase().includes(searchText.toLowerCase()) ||
    item.position.toLowerCase().includes(searchText.toLowerCase())
  );

  const totalPages = Math.ceil(filtered.length / itemsPerPage);
  const start = (currentPage - 1) * itemsPerPage;
  const paginated = filtered.slice(start, start + itemsPerPage);

  return (
    <div>
      <div className="page-header">
        <h1>Сотрудники</h1>

        <Link to="/employees/add" className="add-btn">
          + Новый сотрудник
        </Link>
      </div>

      <div className="search-box">
        <input
          type="text"
          placeholder="Поиск по имени или должности..."
          value={searchText}
          onChange={(e) => {
            setSearch(e.target.value);
            setCurrentPage(1);
          }}
        />
      </div>

      <div className="table-container">
        {paginated.length === 0 ? (
          <div className="empty-state">
            <p>Сотрудников пока нет</p>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Имя</th>
                <th>Должность</th>
                <th>Смена</th>
                <th>Телефон</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {paginated.map((item) => (
                <tr key={item.id}>
                  <td>{item.name}</td>
                  <td>{item.position}</td>
                  <td>{item.shift}</td>
                  <td>{item.phone}</td>
                  <td>
                  <Link to={`/employees/edit/${item.id}`}>Ред.</Link>
                  <a href="#" onClick={(e) => { e.preventDefault(); deleteEmployee(item.id); }}>Удал.</a>
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

export default EmployeesPage;
