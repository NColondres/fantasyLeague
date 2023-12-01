import { useState, useEffect } from 'react'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
import SummonerEnroll from '../components/SummonerEnroll'

function Home(){
    
    const [cookies] = useCookies()
    const navigate = useNavigate()

    useEffect(() => {
        if (cookies.hasOwnProperty("lobby_id")) {
            console.log(cookies)
            navigate(`/lobby/${cookies.lobby_id}`)
        }
    },[cookies.lobby_id])

    return (
        <>
        <div>
            <h3>Welcome to League of Goons</h3>
            <SummonerEnroll text='Create lobby'/>
        </div>
        </>
    )
   
}

export default Home