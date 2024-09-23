import { useState, useEffect, useCallback, forceUpdate } from 'react'
import { useParams } from 'react-router-dom'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
import SummonerEnroll from '../../components/SummonerEnroll'
import PlayerMatches from '../../components/PlayerMatches'
import Icon from '../../components/Icon'
import './Lobby.css'


function Lobby(){

    const [cookies, setCookie, removeCookie] = useCookies()
    const navigate = useNavigate()
    const { id } = useParams()
    const [lobbyData, setLobbyData] = useState({})
    const [playersData, setPlayersData] = useState([])
    const [isPuuidSet, setIsPuuidSet] = useState(false)
    const [isCreatorPuuid, setIsCreatorPuuid] = useState(false)

    

    useEffect(() => {
        if (!Object.prototype.hasOwnProperty.call(cookies, 'lobby_id') || cookies.lobby_id !== id) {
            setCookie('lobby_id', id, {
                path: '/',
                maxAge: 864000,
                domain: 'localhost',
                sameSite: 'strict'
            }) 
            removeCookie('puuid', {path: '/'})
        }
        getLobbyInfo()
    },[])

    useEffect(() => {
        Object.prototype.hasOwnProperty.call(cookies, 'puuid') ? setIsPuuidSet(true) : setIsPuuidSet(false)
        isPuuidSet && lobbyData.creator_puuid === cookies.puuid ? setIsCreatorPuuid(true) : setIsCreatorPuuid(false)
    },[lobbyData, cookies])

    useEffect(()=> {
        const interval =  setInterval(() => {
            getLobbyInfo()
        }, 60000)

        return () => clearInterval(interval)
    },[])

  async function getLobbyInfo(){

    const response = await fetch('http://localhost:8080/lobby', {
        mode: 'cors',
        credentials: 'include'
      })
    
    const data = await response.json()

    setLobbyData(data.lobby)
    setPlayersData(data.players)

    if ('error' in data && data.error == 'lobby not found'){
        deleteLobby()
    }
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
        body: {}
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

            {!isPuuidSet && <SummonerEnroll text='Join' setPlayersData={setPlayersData}/>}
            
            {isCreatorPuuid && !lobbyData.started && <button type="submit" onClick={startLobby}>Start</button>}
            
            {playersData.map(player => 
            
                <div key={player.puuid}>
                    <div id="playerInfo">
                    <div>
                    <Icon id='profileIcon' icon={player.profile_icon_id} type='profileicon'/>
                    </div>
                    <div>
                    {player.total_score > 0 && <strong>{player.total_score}</strong>}
                    <h3 id="playerName">{player.name}</h3>
                    </div>
                    </div>
                    
                    {player.matches?.length && <PlayerMatches matches={player.matches}/>}
                </div>
            )}

            <button type="submit" onClick={deleteLobby}>Delete lobby</button>
        </div>
        </>
    )
}

export default Lobby