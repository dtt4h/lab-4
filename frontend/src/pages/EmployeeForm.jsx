import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../api/api";

function EmployeeForm() {
  const navigate = useNavigate();

  const { id } = useParams();

  const isEdit = !!id;

  const [form, setForm] = useState({
    name: "",
    position: "",
    phone: "",
    email: "",
    shift: "Утро",
  });

  useEffect(() => {
    if (isEdit) {
      fetchEmployee();
    }
  }, []);

  async function fetchEmployee() {
    const res = await api.get(`/employees/${id}`);
    setForm(res.data);
  }

  function handleChange(e) {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  }

  async function handleSubmit(e) {
    e.preventDefault();

    if (!form.name || !form.position || !form.phone) {
      alert("заполни обязательные поля");
      return;
    }

    if (isEdit) {
      await api.put(`/employees/${id}`, form);
    } else {
      await api.post("/employees", form);
    }

    navigate("/employees");
  }

  return (
    <div className="form-card">
      <h1>
        {isEdit
          ? "Редактировать сотрудника"
          : "Новый сотрудник"}
      </h1>

      <form onSubmit={handleSubmit}>
        <input
          name="name"
          placeholder="Имя"
          value={form.name}
          onChange={handleChange}
          required
        />

        <input
          name="position"
          placeholder="Должность"
          value={form.position}
          onChange={handleChange}
          required
        />

        <input
          name="phone"
          placeholder="Телефон"
          value={form.phone}
          onChange={handleChange}
          required
        />

        <input
          name="email"
          placeholder="Email"
          value={form.email}
          onChange={handleChange}
        />

        <select name="shift" value={form.shift} onChange={handleChange} required>
          <option value="Утро">Утро</option>
          <option value="День">День</option>
          <option value="Вечер">Вечер</option>
          <option value="Ночь">Ночь</option>
        </select>

        <button type="submit">Сохранить</button>
      </form>
    </div>
  );
}

export default EmployeeForm;
