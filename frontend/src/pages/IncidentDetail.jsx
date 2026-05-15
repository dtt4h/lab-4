import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/api";

function IncidentDetail() {
  const { id } = useParams();

  const [incident, setIncident] = useState(null);

  useEffect(() => {
    fetchIncident();
  }, []);

  async function fetchIncident() {
    const res = await api.get(`/incidents/${id}`);
    setIncident(res.data);
  }

  if (!incident) return <h2>Загрузка...</h2>;

  return (
    <div className="detail-card">
      <h1>{incident.title}</h1>

      <p>
        <b>Отель:</b> {incident.hotel}
      </p>

      <p>
        <b>Описание:</b> {incident.description}
      </p>

      <p>
        <b>Уровень:</b> {incident.level}
      </p>

      <p>
        <b>Статус:</b> {incident.status}
      </p>
    </div>
  );
}

export default IncidentDetail;
