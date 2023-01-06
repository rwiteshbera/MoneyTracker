import React from "react";
import "./App.css";
import 'bootstrap/dist/css/bootstrap.min.css';
import {Route, Routes} from "react-router-dom"

import Signup from "./components/home/Signup";
import Login from "./components/home/Login";
import Dashboard from "./components/dashboard/Dashboard";

const App = () => {
  return (
    <>
    <h3 className="title">Money Tracker</h3>
      <Routes>
      <Route exact path="/" element={<Login />}></Route>
        <Route exact path="/register" element={<Signup />}></Route>
        <Route exact path="/dashboard" element={<Dashboard />}></Route>
      </Routes>
    </>
  );
};

export default App;
