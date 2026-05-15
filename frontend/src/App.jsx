import { BrowserRouter, Routes, Route } from "react-router-dom";

import Sidebar from "./components/Sidebar";

import IncidentsPage from "./pages/IncidentsPage";
import IncidentDetail from "./pages/IncidentDetail";
import IncidentForm from "./pages/IncidentForm";

import EmployeesPage from "./pages/EmployeesPage";
import EmployeeForm from "./pages/EmployeeForm";

function App() {
  return (
    <BrowserRouter>
      <div className="layout">
        <Sidebar />

        <main className="content">
          <Routes>
            <Route path="/" element={<IncidentsPage />} />

            <Route path="/incidents/:id" element={<IncidentDetail />} />

            <Route path="/incidents/add" element={<IncidentForm />} />

            <Route
              path="/incidents/edit/:id"
              element={<IncidentForm />}
            />

            <Route path="/employees" element={<EmployeesPage />} />

            <Route path="/employees/add" element={<EmployeeForm />} />

            <Route
              path="/employees/edit/:id"
              element={<EmployeeForm />}
            />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
