import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/api";
import { incidentStatusLabel } from "../constants/incidentStatus";

const priorityLabels = { low: "Низкий", medium: "Средний", high: "Высокий", critical: "Критичный" };
const locationLabel = (location) => location ? [location.area_name, location.floor && `этаж ${location.floor}`, location.room && `помещение ${location.room}`].filter(Boolean).join(" · ") : "Не указано";

export default function IncidentDetail() {
  const { id } = useParams();
  const [data, setData] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.get(`/incidents/${id}`).then(async (response) => {
      const incident = response.data.data;
      const [hotels, categories, locations] = await Promise.all([
        api.get("/hotels"), api.get("/categories"), incident.hotel_id ? api.get("/locations", { params: { hotel_id: incident.hotel_id } }) : Promise.resolve({ data: { data: [] } }),
      ]);
      setData({ incident, hotel: hotels.data.data?.find((item) => item.id === incident.hotel_id), category: categories.data.data?.find((item) => item.id === incident.category_id), location: locations.data.data?.find((item) => item.id === incident.location_id) });
    }).catch(() => setError("Не удалось загрузить инцидент."));
  }, [id]);

  if (error) return <p className="form-error" role="alert">{error}</p>;
  if (!data) return <div className="empty-state">Загружаем инцидент…</div>;
  const { incident, hotel, category, location } = data;

  return <section className="detail-card incident-detail"><span className="eyebrow">ИНЦИДЕНТ</span><div className="detail-heading"><div><h1>{incident.title}</h1><p className="muted">Зафиксирован {new Date(incident.occurred_at).toLocaleString("ru-RU")}</p></div><span className={`status-badge incident-status-${incident.status}`}>{incidentStatusLabel(incident.status)}</span></div><div className="detail-grid"><article><span>Гостиница</span><b>{hotel?.name || "Не указана"}</b></article><article><span>Место</span><b>{locationLabel(location)}</b></article><article><span>Категория</span><b>{category?.name || "Не указана"}</b></article><article><span>Приоритет</span><b>{priorityLabels[incident.priority] || incident.priority}</b></article></div><div className="incident-description"><h2>Подробности</h2><p>{incident.description || "Подробности не указаны."}</p></div></section>;
}
