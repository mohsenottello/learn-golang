import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import { AuthProvider } from "./auth";
import { ToastProvider } from "./toast";
import { SiteProvider } from "./site";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <SiteProvider>
          <ToastProvider>
            <App />
          </ToastProvider>
        </SiteProvider>
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>,
);
