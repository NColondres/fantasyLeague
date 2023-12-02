import { useCookies } from 'react-cookie'
import { useState } from 'react'

function SummonerEnroll({text, setPlayersData}){

    const [summonerName, setSummonerName] = useState('')
    const [region, setRegion] = useState('na1')

  async function enrollUser(){
        if (summonerName){

            fetch('http://localhost:8080/enroll', {
                method: "POST",
                mode: 'cors',
                credentials: 'include',
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    summoner: summonerName,
                    region: region
                })
              })
              .then((response) => response.json())
              .then((data) => {
                console.log(data)
                getLobbyInfo()
              })


        
        } else {
            alert('Summoner name empty')
        }
    }
        
    async function getLobbyInfo(){

        const response = await fetch('http://localhost:8080/lobby', {
            mode: 'cors',
            credentials: 'include'
          })
        
        const data = await response.json()
    
        setPlayersData(data.players)
    
        console.log(data)
    
      }


    return (
        <>
        <input name='summonerName' placeholder="Summoner Name"
            value={summonerName} onChange={e => setSummonerName(e.target.value)}
            onKeyDown={e => {
                if (e.key === 'Enter') {enrollUser()}
            }}
        />
        <select id="selectedRegion" name="selectedRegion"
            value={region} onChange={e => setRegion(e.target.value)}
        >
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
        <button type="submit" onClick={enrollUser}>{text}</button>
        </>
    )

}

export default SummonerEnroll