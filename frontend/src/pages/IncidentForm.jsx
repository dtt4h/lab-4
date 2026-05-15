import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../api/api";

function IncidentForm() {
  const navigate = useNavigate();

  const { id } = useParams();

  const isEdit = !!id;

  const [form, setForm] = useState({
    title: "",
    hotel: "Гостиница Космос",
    level: "Средний",
    status: "Открыт",
    description: "",
    date: "",
    floor: "",
  });

  useEffect(() => {
    if (isEdit) {
      fetchIncident();
    } else {
      const now = new Date();
      const offset = now.getTimezoneOffset() * 60000;
      const localISOTime = new Date(now.getTime() - offset).toISOString().slice(0, 16);
      setForm((prev) => ({ ...prev, date: localISOTime }));
    }
  }, []);

  async function fetchIncident() {
    const res = await api.get(`/incidents/${id}`);
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

    if (!form.title || !form.level || !form.status) {
      alert("нужно заполнить название и статус");
      return;
    }

    if (isEdit) {
      await api.put(`/incidents/${id}`, form);
    } else {
      await api.post("/incidents", form);
    }

    navigate("/");
  }

  return (
    <div className="form-card">
      <h1>
        {isEdit
          ? "Редактировать инцидент"
          : "Новый инцидент"}
      </h1>

      <form onSubmit={handleSubmit}>
        <input
          name="title"
          placeholder="Название"
          value={form.title}
          onChange={handleChange}
          required
        />

        <input
          name="hotel"
          placeholder="Отель"
          value={form.hotel}
          onChange={handleChange}
          disabled
          style={{ background: "#f5f5f5" }}
        />

        <select name="level" value={form.level} onChange={handleChange} required>
          <option value="Низкий">Низкий</option>
          <option value="Средний">Средний</option>
          <option value="Высокий">Высокий</option>
        </select>

        <select name="status" value={form.status} onChange={handleChange} required>
          <option value="Открыт">Открыт</option>
          <option value="Проверяется">Проверяется</option>
          <option value="Решён">Решён</option>
        </select>

        <input
          type="datetime-local"
          name="date"
          placeholder="Дата и время"
          value={form.date}
          onChange={handleChange}
          required
        />

        <input
          name="floor"
          placeholder="Этаж"
          value={form.floor}
          onChange={handleChange}
        />

        <textarea
          name="description"
          placeholder="Описание"
          value={form.description}
          onChange={handleChange}
        />

        <button type="submit">Сохранить</button>
      </form>
    </div>
  );
}

export default IncidentForm;
