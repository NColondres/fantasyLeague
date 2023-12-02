import { useState, useEffect, useCallback, forceUpdate } from 'react'
import { useParams } from 'react-router-dom'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
import SummonerEnroll from '../components/SummonerEnroll'


function Lobby(){

    const [cookies, setCookie, removeCookie] = useCookies()
    const navigate = useNavigate()
    const { id } = useParams()
    const [lobbyData, setLobbyData] = useState({})
    const [playersData, setPlayersData] = useState([])
    const [isPuuidSet, setIsPuuidSet] = useState(false)
    const [isCreatorPuuid, setIsCreatorPuuid] = useState(false)

    if (!cookies.hasOwnProperty('lobby_id') || cookies.lobby_id !== id) {
        setCookie('lobby_id', id, {
            path: '/',
            maxAge: 864000,
            domain: 'localhost',
            sameSite: 'strict'
        }) 
        removeCookie('puuid', {path: '/'})
    }

    useEffect(() => {
        getLobbyInfo()
    },[])

    useEffect(() => {
        cookies.hasOwnProperty('puuid') ? setIsPuuidSet(true) : setIsPuuidSet(false)
        isPuuidSet && lobbyData.creator_puuid === cookies.puuid ? setIsCreatorPuuid(true) : setIsCreatorPuuid(false)
    },[lobbyData, cookies])

  async function getLobbyInfo(){

    const response = await fetch('http://localhost:8080/lobby', {
        mode: 'cors',
        credentials: 'include'
      })
    
    const data = await response.json()

    setLobbyData(data.lobby)
    setPlayersData(data.players)

    console.log(data)

  }

  async function startLobby(){
    const response = await fetch('http://localhost:8080/start', {
        method: "POST",
        mode: 'cors',
        credentials: 'include',
        headers: {
            "Content-Type": "application/json"
        },
    })
    const data = await response.json()

    if (!response.ok) {
        alert(`${response.statusText}: ${data.error}`)
    } else {
        getLobbyInfo()
        console.log("Starting lobby:", data)
    }
  }

  async function deleteLobby(){
    console.log("Deleting lobby")
    removeCookie('lobby_id', {path: '/'})
    removeCookie('puuid', {path: '/'})
    navigate('/')
    
  }

    return (
        <>
        <div>
            <h3>Lobby Id: {lobbyData.id} </h3>
            
            {!isPuuidSet && <SummonerEnroll text='Join' setPlayersData={setPlayersData}/>}
            
            {isCreatorPuuid && !lobbyData.started && <button type="submit" onClick={startLobby}>Start</button>}
            
            {playersData.map(player => 
                <div key={player.puuid}>
                    <h4>Name: {player.name}</h4>
                </div>
            )}

            <button type="submit" onClick={deleteLobby}>Delete lobby</button>
        </div>
        </>
    )
}

export default Lobby