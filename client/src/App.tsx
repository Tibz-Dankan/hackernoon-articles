import React, { Fragment } from "react";
import { Home } from "./components/Home";
import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import "./App.css";

export const App: React.FC = () => {
  return (
    <div className="bg-(--clr-background)  min-h-screen overflow-x-hidden">
      <BrowserRouter>
        <Fragment>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/home" element={<Navigate to="/" replace />} />
            <Route path="/*" element={<Navigate to="/" replace />} />
          </Routes>
        </Fragment>
      </BrowserRouter>
    </div>
  );
};
