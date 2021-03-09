import React from 'react';
import {BrowserRouter} from "react-router-dom";
import axios from "axios";

import PrivateRoute from './components/PrivateRoute';
import PublicRoute from './components/PublicRoute';

import Dashboard from './components/Dashboard'
import Home from './components/Home'

import './App.scss'

axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

function App() {
  return (
    <div className="h-100">
      <BrowserRouter>
        <PublicRoute exact path="/" component={Home}/>
        <PrivateRoute exact path="/dashboard" component={Dashboard}/>
      </BrowserRouter>
    </div>
  );
}

export default App;
