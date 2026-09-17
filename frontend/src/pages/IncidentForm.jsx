import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import api from "../api/api";

const emptyForm = { title: "", description: "", priority: "medium", status: "open", hotel_id: "", location_id: "", category_id: "", occurred_at: "" };
const localNow = () => new Date(Date.now() - new Date().getTimezoneOffset() * 60000).toISOString().slice(0, 16);

export default function IncidentForm() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const [form, setForm] = useState(emptyForm);
  const [hotels, setHotels] = useState([]);
  const [categories, setCategories] = useState([]);
  const [locations, setLocations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    Promise.all([api.get("/hotels"), api.get("/categories"), isEdit ? api.get(`/incidents/${id}`) : Promise.resolve(null)])
      .then(([hotelResponse, categoryResponse, incidentResponse]) => {
        setHotels(hotelResponse.data.data || []);
        setCategories(categoryResponse.data.data || []);
        if (incidentResponse) {
          const item = incidentResponse.data.data;
          setForm({ ...emptyForm, ...item, location_id: item.location_id || "", category_id: item.category_id || "", occurred_at: item.occurred_at?.slice(0, 16) || localNow() });
        } else setForm((current) => ({ ...current, occurred_at: localNow() }));
      })
      .catch(() => setError("Не удалось загрузить данные формы. Обновите страницу и попробуйте снова."))
      .finally(() => setLoading(false));
  }, [id, isEdit]);

  useEffect(() => {
    if (!form.hotel_id) { setLocations([]); return; }
    api.get("/locations", { params: { hotel_id: form.hotel_id } }).then((response) => setLocations(response.data.data || [])).catch(() => setError("Не удалось загрузить локации выбранной гостиницы."));
  }, [form.hotel_id]);

  function change(event) {
    const { name, value } = event.target;
    setError("");
    setForm((current) => ({ ...current, [name]: value, ...(name === "hotel_id" ? { location_id: "" } : {}) }));
  }

  async function submit(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    try {
      const body = { ...form, location_id: form.location_id || null, category_id: form.category_id || null, occurred_at: new Date(form.occurred_at).toISOString() };
      if (isEdit) await api.put(`/incidents/${id}`, body); else await api.post("/incidents", body);
      navigate("/incidents");
    } catch (requestError) {
      setError(requestError.response?.data?.error?.message || "Не удалось сохранить инцидент. Проверьте заполненные поля.");
    } finally { setSaving(false); }
  }

  if (loading) return <div className="empty-state">Подготавливаем форму…</div>;
  return <section className="form-card incident-form"><span className="eyebrow">ИНЦИДЕНТ</span><h1>{isEdit ? "Редактирование инцидента" : "Новый инцидент"}</h1><p className="muted">Укажите место и детали — это поможет быстрее передать задачу нужному сотруднику.</p>{error && <p className="form-error" role="alert">{error}</p>}
    <form onSubmit={submit}><label>Краткое описание<input name="title" value={form.title} onChange={change} placeholder="Например: не работает кондиционер" required /></label><div className="form-grid"><label>Гостиница<select name="hotel_id" value={form.hotel_id} onChange={change} required><option value="">Выберите гостиницу</option>{hotels.map((hotel) => <option key={hotel.id} value={hotel.id}>{hotel.name}</option>)}</select></label><label>Место<select name="location_id" value={form.location_id} onChange={change} disabled={!form.hotel_id}><option value="">{form.hotel_id ? "Не указано" : "Сначала выберите гостиницу"}</option>{locations.map((location) => <option key={location.id} value={location.id}>{[location.area_name, location.floor && `этаж ${location.floor}`, location.room && `помещение ${location.room}`].filter(Boolean).join(" · ")}</option>)}</select></label></div><div className="form-grid"><label>Категория<select name="category_id" value={form.category_id} onChange={change}><option value="">Без категории</option>{categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}</select></label><label>Приоритет<select name="priority" value={form.priority} onChange={change}><option value="low">Низкий</option><option value="medium">Средний</option><option value="high">Высокий</option><option value="critical">Критичный</option></select></label></div><div className="form-grid"><label>Статус<select name="status" value={form.status} onChange={change}><option value="open">Открыт</option><option value="in_progress">В работе</option><option value="resolved">Решён</option><option value="closed">Закрыт</option></select></label><label>Когда обнаружено<input type="datetime-local" name="occurred_at" value={form.occurred_at} onChange={change} required /></label></div><label>Подробности<textarea name="description" placeholder="Что произошло, как это влияет на гостя или работу и что уже проверили" value={form.description} onChange={change} /></label><button disabled={saving}>{saving ? "Сохраняем…" : "Сохранить инцидент"}</button></form>
  </section>;
}
