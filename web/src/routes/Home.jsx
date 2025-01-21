import { useEffect } from 'react'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
import SummonerEnroll from '../components/SummonerEnroll'

function Home(){
    
    const [cookies] = useCookies()
    const navigate = useNavigate()

    useEffect(() => {
        if (Object.prototype.hasOwnProperty.call(cookies, "lobby_id")) {
            navigate(`/lobby/${cookies.lobby_id}`)
        }
    },[cookies.lobby_id])

    return (
        <>
        <div>
            <h2 id="title">Welcome to League of Goons</h2>
            <SummonerEnroll text='Create lobby'/>
        </div>
        </>
    )
   
}

export default Home