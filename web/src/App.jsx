import { useState, useEffect } from 'react'
import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import './App.css'
import Home from './routes/Home.jsx'
import Lobby from './routes/Lobby.jsx'
import NotFound from './routes/NotFound.jsx';

function App() {


  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/lobby/:id" element={<Lobby />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Router>
  )
}

export default App
