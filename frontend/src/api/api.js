import axios from "axios";

let accessToken = null;
let currentUser = null;
export const setAccessToken = (token) => { accessToken = token; };
export const setCurrentUser = (user) => { currentUser = user; };
export const getCurrentUser = () => currentUser;
export const can = (permission) => (currentUser?.permissions || []).includes("*") || (currentUser?.permissions || []).includes(permission);
const api = axios.create({ baseURL: "/api/v1", withCredentials: true, headers: { "Content-Type": "application/json" } });
api.interceptors.request.use((config) => { if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`; return config; });
api.interceptors.response.use((response) => response, async (error) => {
  const request = error.config;
  if (error.response?.status === 401 && !request?._retried && !request?.url?.includes("/auth/refresh")) {
    request._retried = true;
    try { const response = await api.post("/auth/refresh"); setAccessToken(response.data.data.access_token); request.headers.Authorization = `Bearer ${accessToken}`; return api(request); } catch { setAccessToken(null); window.location.replace("/login"); }
  }
  return Promise.reject(error);
});
export default api;
