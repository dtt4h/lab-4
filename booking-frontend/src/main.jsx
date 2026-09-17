import React, { useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import axios from "axios";
import "./style.css";

const api = axios.create({ baseURL: "/api/v1" });
const today = () => new Date().toISOString().slice(0, 10);
const nextDay = (date) => new Date(new Date(`${date}T00:00:00`).getTime() + 86400000).toISOString().slice(0, 10);
const initialForm = { hotel_id: "", check_in: "", check_out: "", guests_count: 1, guest_name: "", guest_email: "", guest_phone: "" };

function App() {
  const [view, setView] = useState("home");
  return view === "home" ? <Home onChoose={() => setView("rooms")} /> : <Rooms onBack={() => setView("home")} />;
}

function Home({ onChoose }) {
  return <main className="hero"><div className="hero-content"><span className="eyebrow">ПОИСК ВЫГОДНЫХ ОТЕЛЕЙ</span><h1>Найдём отель по хорошей цене</h1><p>Сравните доступные номера, выберите даты и оформите бронь без лишних шагов.</p><div className="hero-actions"><button className="choose-button" onClick={onChoose}>Найти номер <span aria-hidden="true">→</span></button><a className="admin-button" href="http://localhost:5173">Открыть панель управления</a></div></div></main>;
}

function Rooms({ onBack }) {
  const [hotels, setHotels] = useState([]);
  const [rooms, setRooms] = useState([]);
  const [selected, setSelected] = useState(null);
  const [form, setForm] = useState(initialForm);
  const [notice, setNotice] = useState(null);
  const [loadingRooms, setLoadingRooms] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const canSearch = form.hotel_id && form.check_in && form.check_out;
  const nights = useMemo(() => canSearch ? Math.round((new Date(form.check_out) - new Date(form.check_in)) / 86400000) : 0, [canSearch, form.check_in, form.check_out]);

  useEffect(() => {
    api.get("/public/hotels").then((response) => setHotels(response.data.data || [])).catch(() => setNotice({ type: "error", text: "Не удалось загрузить список гостиниц. Попробуйте обновить страницу." }));
  }, []);

  useEffect(() => {
    if (!canSearch) { setRooms([]); return; }
    if (nights < 1) { setRooms([]); setNotice({ type: "error", text: "Дата выезда должна быть позже даты заезда." }); return; }
    setLoadingRooms(true);
    setNotice(null);
    api.get("/public/rooms", { params: { hotel_id: form.hotel_id, check_in: form.check_in, check_out: form.check_out } })
      .then((response) => setRooms(response.data.data || []))
      .catch(() => setNotice({ type: "error", text: "Не удалось проверить доступность номеров. Проверьте даты и попробуйте снова." }))
      .finally(() => setLoadingRooms(false));
  }, [canSearch, form.hotel_id, form.check_in, form.check_out, nights]);

  function change(event) {
    const { name, value } = event.target;
    setNotice(null);
    setForm((current) => {
      const next = { ...current, [name]: value };
      if (name === "check_in" && current.check_out && current.check_out <= value) next.check_out = "";
      return next;
    });
    if (name === "hotel_id" || name === "check_in" || name === "check_out") setSelected(null);
  }

  async function book(event) {
    event.preventDefault();
    if (!selected || nights < 1) return;
    setSubmitting(true);
    setNotice(null);
    try {
      const response = await api.post("/public/bookings", { ...form, room_id: selected.id, guests_count: Number(form.guests_count) });
      const booking = response.data.data;
      setNotice({ type: "success", text: `Заявка на бронирование принята. Итоговая стоимость: ${Number(booking.total_price).toLocaleString("ru-RU")} ₽. Ожидайте подтверждения от отеля.` });
      setRooms((current) => current.filter((room) => room.id !== selected.id));
      setSelected(null);
      setForm((current) => ({ ...initialForm, hotel_id: current.hotel_id, check_in: current.check_in, check_out: current.check_out }));
    } catch (error) {
      setNotice({ type: "error", text: error.response?.status === 409 ? "Этот номер только что забронировали на выбранные даты. Выберите другой вариант." : error.response?.data?.error?.message || "Не удалось оформить бронь. Попробуйте ещё раз." });
    } finally { setSubmitting(false); }
  }

  return <main className="rooms-page"><button className="back-button" onClick={onBack}>← На главную</button><section className="rooms-content">
    <span className="eyebrow dark">ПОИСК ВЫГОДНЫХ ОТЕЛЕЙ</span><h2>Выберите номер</h2><p className="subtitle">Сначала укажите даты — покажем только свободные номера и точную стоимость.</p>
    <div className="search-grid" aria-label="Параметры поиска номера">
      <label className="sr-only" htmlFor="hotel">Гостиница</label><select id="hotel" name="hotel_id" value={form.hotel_id} onChange={change}><option value="">Гостиница</option>{hotels.map((hotel) => <option key={hotel.id} value={hotel.id}>{hotel.name}</option>)}</select>
      <label className="sr-only" htmlFor="check-in">Дата заезда</label><input id="check-in" name="check_in" type="date" min={today()} value={form.check_in} onChange={change} />
      <label className="sr-only" htmlFor="check-out">Дата выезда</label><input id="check-out" name="check_out" type="date" min={form.check_in ? nextDay(form.check_in) : today()} value={form.check_out} onChange={change} disabled={!form.check_in} />
    </div>
    {notice && <p className={`message ${notice.type === "success" ? "message-success" : ""}`} role="status" aria-live="polite">{notice.text}</p>}
    {canSearch && <div className="room-list">{loadingRooms ? <p className="empty">Проверяем доступность номеров…</p> : rooms.map((room) => <article className="room-card" key={room.id}><div><span className="room-label">Номер</span><strong>{room.number}</strong></div><span>{room.room_type_name}</span><span>до {room.capacity} гостей</span><b>{Number(room.price_per_night).toLocaleString("ru-RU")} ₽ <small>/ ночь</small></b><button className="book-button" onClick={() => setSelected(room)}>Выбрать</button></article>)}{!loadingRooms && !rooms.length && <p className="empty">На эти даты свободных номеров нет. Измените даты или выберите другую гостиницу.</p>}</div>}
    {selected && <section className="booking-card" aria-labelledby="booking-title"><div className="booking-title"><div><span className="eyebrow dark">ОФОРМЛЕНИЕ</span><h2 id="booking-title">Номер {selected.number}</h2></div><button className="close-button" onClick={() => setSelected(null)} aria-label="Закрыть форму">×</button></div><p className="subtitle">{selected.room_type_name} · {form.check_in} — {form.check_out}</p><p className="total-price">Итого за {nights} {nights === 1 ? "ночь" : "ночей"}: <b>{(Number(selected.price_per_night) * nights).toLocaleString("ru-RU")} ₽</b></p><form onSubmit={book}><label>Ваше имя<input name="guest_name" value={form.guest_name} onChange={change} autoComplete="name" required /></label><label>Email для подтверждения<input name="guest_email" type="email" value={form.guest_email} onChange={change} autoComplete="email" required /></label><label>Телефон<input name="guest_phone" value={form.guest_phone} onChange={change} autoComplete="tel" required /></label><label>Гостей<input name="guests_count" type="number" min="1" max={selected.capacity} value={form.guests_count} onChange={change} required /></label><button className="book-button" disabled={submitting}>{submitting ? "Оформляем…" : "Подтвердить бронирование"}</button></form></section>}
  </section></main>;
}

createRoot(document.getElementById("root")).render(<App />);
