import { useEffect, useState } from "react";
import api from "../api/api";

const emptyType = { name: "", description: "", capacity: 1, price_per_night: "" };
const emptyRoom = { room_type_id: "", number: "", floor: "" };

export default function RoomsPage() {
  const [hotels, setHotels] = useState([]);
  const [hotelID, setHotelID] = useState("");
  const [types, setTypes] = useState([]);
  const [typeForm, setTypeForm] = useState(emptyType);
  const [roomForm, setRoomForm] = useState(emptyRoom);
  const [notice, setNotice] = useState("");
  const [savingType, setSavingType] = useState(false);
  const [savingRoom, setSavingRoom] = useState(false);

  useEffect(() => { api.get("/hotels").then((response) => setHotels(response.data.data || [])).catch(() => setNotice("Не удалось загрузить гостиницы.")); }, []);
  const loadTypes = () => { if (!hotelID) { setTypes([]); return; } api.get("/admin/room-types", { params: { hotel_id: hotelID } }).then((response) => setTypes(response.data.data || [])).catch(() => setNotice("Не удалось загрузить типы номеров.")); };
  useEffect(() => { setRoomForm(emptyRoom); loadTypes(); }, [hotelID]);
  const changeType = (event) => setTypeForm((current) => ({ ...current, [event.target.name]: event.target.value }));
  const changeRoom = (event) => setRoomForm((current) => ({ ...current, [event.target.name]: event.target.value }));

  async function createType(event) {
    event.preventDefault(); setSavingType(true); setNotice("");
    try { await api.post("/admin/room-types", { ...typeForm, hotel_id: hotelID, capacity: Number(typeForm.capacity), price_per_night: Number(typeForm.price_per_night) }); setTypeForm(emptyType); setNotice("Тип номера добавлен. Теперь выберите его при создании номера."); loadTypes(); }
    catch (error) { setNotice(error.response?.data?.error?.message || "Не удалось добавить тип номера."); }
    finally { setSavingType(false); }
  }
  async function createRoom(event) {
    event.preventDefault(); setSavingRoom(true); setNotice("");
    try { await api.post("/admin/rooms", { ...roomForm, hotel_id: hotelID }); setRoomForm(emptyRoom); setNotice("Номер успешно добавлен."); }
    catch (error) { setNotice(error.response?.data?.error?.message || "Не удалось добавить номер."); }
    finally { setSavingRoom(false); }
  }

  return <div><div className="page-header"><div><span className="eyebrow">НОМЕРНОЙ ФОНД</span><h1>Номера и типы</h1><p className="muted">Сначала создайте тип номера, затем добавьте конкретные номера этого типа.</p></div></div><section className="room-hotel-picker"><label>Гостиница<select value={hotelID} onChange={(event) => { setHotelID(event.target.value); setNotice(""); }}><option value="">Выберите гостиницу</option>{hotels.map((hotel) => <option key={hotel.id} value={hotel.id}>{hotel.name}</option>)}</select></label></section>{notice && <p className="form-error" role="status">{notice}</p>}{hotelID && <div className="catalog-grid room-management"><section className="catalog-card"><h2>Добавить тип номера</h2><p className="muted">Например: «Стандарт», «Семейный» или «Люкс».</p><form onSubmit={createType}><label>Название<input name="name" value={typeForm.name} onChange={changeType} placeholder="Стандарт" required /></label><div className="form-grid"><label>Вместимость<input name="capacity" type="number" min="1" value={typeForm.capacity} onChange={changeType} required /></label><label>Цена за ночь, ₽<input name="price_per_night" type="number" min="0" step="0.01" value={typeForm.price_per_night} onChange={changeType} required /></label></div><label>Описание<input name="description" value={typeForm.description} onChange={changeType} placeholder="Кратко об особенностях номера" /></label><button disabled={savingType}>{savingType ? "Добавляем…" : "Добавить тип"}</button></form></section><section className="catalog-card"><h2>Добавить номер</h2>{types.length ? <form onSubmit={createRoom}><label>Тип номера<select name="room_type_id" value={roomForm.room_type_id} onChange={changeRoom} required><option value="">Выберите тип номера</option>{types.map((type) => <option key={type.id} value={type.id}>{type.name} · до {type.capacity} гостей · {Number(type.price_per_night).toLocaleString("ru-RU")} ₽</option>)}</select></label><div className="form-grid"><label>Номер<input name="number" value={roomForm.number} onChange={changeRoom} placeholder="Например, 203" required /></label><label>Этаж<input name="floor" value={roomForm.floor} onChange={changeRoom} placeholder="Например, 2" required /></label></div><button disabled={savingRoom}>{savingRoom ? "Добавляем…" : "Добавить номер"}</button></form> : <div className="empty-state">Сначала создайте хотя бы один тип номера.</div>}</section></div>}{hotelID && <section className="room-types-list"><h2>Типы номеров в гостинице</h2>{types.length ? <div className="detail-grid">{types.map((type) => <article key={type.id}><span>{type.name}</span><b>до {type.capacity} гостей · {Number(type.price_per_night).toLocaleString("ru-RU")} ₽ / ночь</b><small>{type.description || "Описание не указано"}</small></article>)}</div> : null}</section>}</div>;
}
