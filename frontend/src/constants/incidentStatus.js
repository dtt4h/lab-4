export const INCIDENT_STATUS_LABELS = {
  open: "Открыт",
  in_progress: "В работе",
  resolved: "Решён",
  closed: "Закрыт",
};

export const incidentStatusLabel = (status) => INCIDENT_STATUS_LABELS[status] || status;
