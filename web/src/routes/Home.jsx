import { useState, useEffect } from 'react'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
function Home(){
    
    const [cookies] = useCookies()
    const navigate = useNavigate()
    const [input, setInput] = useState("")

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
            <input name='summonerName' placeholder="Summoner Name"></input>
            <select id="selectedRegion" name="selectedRegion">
                <option value="na1">NA1</option>
                <option value="euw1">EUW1</option>
                <option value="eun1">ENU1</option>
                <option value="br1">BR1</option>
                <option value="jp1">JP1</option>
                <option value="kr">KR</option>
                <option value="la1">LA1</option>
                <option value="oc1">OC1</option>
                <option value="tr1">TR1</option>
                <option value="ru">RU</option>
                <option value="ph2">PH2</option>
                <option value="sg2">SG2</option>
                <option value="th2">TH2</option>
                <option value="tw2">TW2</option>
                <option value="vn2">VN2</option>
            </select>
        </div>
        </>
    )
   
}

export default Home