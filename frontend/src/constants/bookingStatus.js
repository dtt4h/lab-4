export const BOOKING_STATUSES = {
  pending: {
    label: "Ожидает оплаты",
    className: "status-pending",
    actions: [{ value: "paid", label: "Отметить оплату" }, { value: "cancelled", label: "Отменить" }],
  },
  paid: {
    label: "Оплачено",
    className: "status-paid",
    actions: [{ value: "confirmed", label: "Подтвердить" }, { value: "cancelled", label: "Отменить" }],
  },
  confirmed: {
    label: "Подтверждено",
    className: "status-confirmed",
    actions: [{ value: "cancelled", label: "Отменить" }],
  },
  cancelled: { label: "Отменено", className: "status-cancelled", actions: [] },
};

export const bookingStatusLabel = (status) => BOOKING_STATUSES[status]?.label || status;
export const bookingStatusActions = (status) => BOOKING_STATUSES[status]?.actions || [];
