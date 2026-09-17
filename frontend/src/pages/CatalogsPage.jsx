import { useEffect, useState } from "react";
import api from "../api/api";

function Catalog({ title, endpoint, placeholder }) {
  const [items, setItems] = useState([]);
  const [value, setValue] = useState("");
  const [error, setError] = useState("");
  const load = () => api.get(endpoint).then((response) => setItems(response.data.data || [])).catch(() => setError("Не удалось загрузить справочник."));
  useEffect(() => { load(); }, [endpoint]);
  async function add(event) { event.preventDefault(); if (!value.trim()) return; try { await api.post(endpoint, endpoint === "/hotels" ? { name: value, address: "Адрес не указан" } : { name: value }); setValue(""); setError(""); await load(); } catch (requestError) { setError(requestError.response?.data?.error?.message || "Не удалось сохранить запись."); } }
  async function remove(id) { if (!window.confirm("Удалить запись?")) return; try { await api.delete(`${endpoint}/${id}`); await load(); } catch (requestError) { setError(requestError.response?.data?.error?.message || "Нельзя удалить используемую запись."); } }
  return <section className="catalog-card"><h2>{title}</h2><form className="inline-form" onSubmit={add}><label className="sr-only" htmlFor={endpoint}>{title}</label><input id={endpoint} value={value} onChange={(event) => setValue(event.target.value)} placeholder={placeholder} required /><button>Добавить</button></form>{error && <p className="form-error" role="alert">{error}</p>}<ul className="catalog-list">{items.map((item) => <li key={item.id}><span>{item.name}</span><button className="link-danger" onClick={() => remove(item.id)}>Удалить</button></li>)}{!items.length && <li className="muted">Пока пусто</li>}</ul></section>;
}

function LocationsCatalog() {
  const [hotels, setHotels] = useState([]);
  const [items, setItems] = useState([]);
  const [form, setForm] = useState({ hotel_id: "", area_name: "", floor: "", room: "" });
  const [error, setError] = useState("");
  const loadLocations = (hotelID) => { if (!hotelID) { setItems([]); return; } api.get("/locations", { params: { hotel_id: hotelID } }).then((response) => setItems(response.data.data || [])).catch(() => setError("Не удалось загрузить локации.")); };
  useEffect(() => { api.get("/hotels").then((response) => setHotels(response.data.data || [])).catch(() => setError("Не удалось загрузить гостиницы.")); }, []);
  useEffect(() => { loadLocations(form.hotel_id); }, [form.hotel_id]);
  const change = (event) => setForm((current) => ({ ...current, [event.target.name]: event.target.value }));
  async function add(event) { event.preventDefault(); setError(""); try { await api.post("/locations", form); setForm((current) => ({ ...current, area_name: "", floor: "", room: "" })); await loadLocations(form.hotel_id); } catch (requestError) { setError(requestError.response?.data?.error?.message || "Не удалось добавить локацию."); } }
  async function remove(id) { if (!window.confirm("Удалить локацию?")) return; try { await api.delete(`/locations/${id}`); await loadLocations(form.hotel_id); } catch (requestError) { setError(requestError.response?.data?.error?.message || "Нельзя удалить локацию, которая уже используется."); } }
  return <section className="catalog-card location-catalog"><h2>Локации</h2><p className="muted">Создайте место один раз — затем его можно выбрать в форме инцидента.</p><form onSubmit={add}><label>Гостиница<select name="hotel_id" value={form.hotel_id} onChange={change} required><option value="">Выберите гостиницу</option>{hotels.map((hotel) => <option key={hotel.id} value={hotel.id}>{hotel.name}</option>)}</select></label><div className="form-grid"><label>Зона или название<input name="area_name" value={form.area_name} onChange={change} placeholder="Например, лобби" /></label><label>Этаж<input name="floor" value={form.floor} onChange={change} placeholder="Например, 2" required /></label></div><label>Номер помещения (необязательно)<input name="room" value={form.room} onChange={change} placeholder="Например, 205" /></label><button disabled={!form.hotel_id}>Добавить локацию</button></form>{error && <p className="form-error" role="alert">{error}</p>}<ul className="catalog-list">{items.map((item) => <li key={item.id}><span>{[item.area_name, `этаж ${item.floor}`, item.room && `помещение ${item.room}`].filter(Boolean).join(" · ")}</span><button className="link-danger" onClick={() => remove(item.id)}>Удалить</button></li>)}{form.hotel_id && !items.length && <li className="muted">Для этой гостиницы локаций пока нет.</li>}</ul></section>;
}

export default function CatalogsPage() {
  return <div><div className="page-header"><div><span className="eyebrow">НАСТРОЙКИ</span><h1>Справочники</h1><p className="muted">Здесь настраиваются общие данные для форм и инцидентов.</p></div></div><div className="catalog-grid"><Catalog title="Гостиницы" endpoint="/hotels" placeholder="Название гостиницы" /><Catalog title="Категории" endpoint="/categories" placeholder="Например, техника" /><LocationsCatalog /></div></div>;
}
