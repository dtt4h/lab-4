import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { useEffect, useState } from "react";
import api, { can, setAccessToken, setCurrentUser } from "./api/api";
import Sidebar from "./components/Sidebar";
import IncidentsPage from "./pages/IncidentsPage";
import IncidentDetail from "./pages/IncidentDetail";
import IncidentForm from "./pages/IncidentForm";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import CatalogsPage from "./pages/CatalogsPage";
import DashboardPage from "./pages/DashboardPage";
import BookingsPage from "./pages/BookingsPage";
import AccessPage from "./pages/AccessPage";
import RoomsPage from "./pages/RoomsPage";

function Denied() {
  return <div className="detail-card"><span className="eyebrow">ДОСТУП ОГРАНИЧЕН</span><h1>Недостаточно прав</h1><p className="muted">У вас нет доступа к этому разделу.</p></div>;
}

function PermissionGuard({ permission, anyOf, children }) {
  return (permission ? can(permission) : anyOf?.some((item) => can(item))) ? children : <Denied />;
}

function Protected({ children }) {
  const [ready, setReady] = useState(false);
  const [user, setUser] = useState(null);

  useEffect(() => {
    api.post("/auth/refresh")
      .then(({ data }) => { setAccessToken(data.data.access_token); return api.get("/auth/me"); })
      .then((response) => { setCurrentUser(response.data.data); setUser(response.data.data); })
      .catch(() => setUser(null))
      .finally(() => setReady(true));
  }, []);

  if (!ready) return <div className="loading-screen">Проверяем сессию…</div>;
  if (!user || !user.permissions?.length) return <Navigate to="/login" replace />;
  return children;
}

export default function App() {
  return <BrowserRouter><Routes>
    <Route path="/login" element={<LoginPage />} />
    <Route path="/register" element={<RegisterPage />} />
    <Route path="/*" element={<Protected><div className="layout"><Sidebar /><main className="content"><Routes>
      <Route path="/" element={<PermissionGuard anyOf={["hotel.read", "booking.read", "incident.read"]}><DashboardPage /></PermissionGuard>} />
      <Route path="/dashboard" element={<PermissionGuard anyOf={["hotel.read", "booking.read", "incident.read"]}><DashboardPage /></PermissionGuard>} />
      <Route path="/bookings" element={<PermissionGuard permission="booking.read"><BookingsPage /></PermissionGuard>} />
      <Route path="/rooms" element={<PermissionGuard permission="room.read"><RoomsPage /></PermissionGuard>} />
      <Route path="/access" element={<PermissionGuard permission="roles.manage"><AccessPage /></PermissionGuard>} />
      <Route path="/incidents" element={<PermissionGuard permission="incident.read"><IncidentsPage /></PermissionGuard>} />
      <Route path="/incidents/:id" element={<PermissionGuard permission="incident.read"><IncidentDetail /></PermissionGuard>} />
      <Route path="/incidents/add" element={<PermissionGuard permission="incident.manage"><IncidentForm /></PermissionGuard>} />
      <Route path="/incidents/edit/:id" element={<PermissionGuard permission="incident.manage"><IncidentForm /></PermissionGuard>} />
      <Route path="/catalogs" element={<PermissionGuard permission="catalog.manage"><CatalogsPage /></PermissionGuard>} />
    </Routes></main></div></Protected>} />
  </Routes></BrowserRouter>;
}
