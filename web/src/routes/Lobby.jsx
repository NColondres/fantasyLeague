import { useState, useEffect } from 'react'
import { useParams, redirect } from 'react-router-dom'
import { useCookies } from 'react-cookie'


function Lobby(){

    const [cookies, setCookie] = useCookies()
    const { id } = useParams()
    const [lobbyData, setLobbyData] = useState({});
    const [playersData, setPlayersData] = useState([]);

    if (!cookies.hasOwnProperty('lobby_id')) {
        setCookie('lobby_id', id, {
            path: '/',
            maxAge: 864000,
            domain: 'localhost'
        })   
    } else if (cookies.lobby_id !== id) {

        setCookie('lobby_id', id, {
            path: '/',
            maxAge: 864000,
            domain: 'localhost'
        }) 

    }


  useEffect(() =>{
    fetch('http://localhost:8080/lobby', {
      mode: 'cors',
      credentials: 'include'
    })
    .then((response) => response.json())
    .then((data) => {
        setLobbyData(data.lobby)
        setPlayersData(data.players)
    })
  },[])

    return (
        <>
        <div>
            <h3>Lobby Id: {lobbyData.id} </h3>
            {playersData.map(player => 
                <div key={player.puuid}>
                    <h4>Name: {player.name}</h4>
                </div>
            )}
        </div>
        </>
    )
}

export default Lobby