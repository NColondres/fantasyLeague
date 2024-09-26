import { useState, useEffect } from 'react'
import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import './App.css'
import Home from './routes/Home.jsx'
import Lobby from './routes/Lobby/Lobby.jsx'
import NotFound from './routes/NotFound.jsx';

function App() {

  const [dataDragonLatestAPIVersion, setdataDragonLatestAPIVersion] = useState('')

  useEffect(() =>{

    async function getDataDragonLatestAPIVersion(){
      const response = await fetch('https://ddragon.leagueoflegends.com/api/versions.json')
      const  data = await response.json()
      console.log(" Response Status:", response.status, "Data Dragon Version:", data[0])
      if (response.status === 200){
        setdataDragonLatestAPIVersion(data[0])
      }
    }
    
    getDataDragonLatestAPIVersion()
  }, [dataDragonLatestAPIVersion])

  

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/lobby/:id" element={<Lobby dataDragonVersion={dataDragonLatestAPIVersion} />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Router>
  )
}

export default App
