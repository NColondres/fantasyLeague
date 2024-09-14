import { useState } from 'react'
import { useCookies } from 'react-cookie'

function SummonerEnroll({text, setPlayersData}){

    const [gameName, setgameName] = useState('')
    const [tagLine, settagLine] = useState('#NA1')
    const [cookies] = useCookies()

  async function enrollUser(){
        if (gameName){
            fetch('http://localhost:8080/enroll', {
                method: "POST",
                mode: 'cors',
                credentials: 'include',
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    gameName: gameName,
                    tagLine: tagLine.slice(1)
                })
              })
              .then((response) => response.json())
              .then((data) => {
                if (Object.prototype.hasOwnProperty.call(cookies, "lobby_id")){
                    getLobbyInfo()
                }
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
    
        console.log("Lobby Reponse Data: ", data)

        setPlayersData(data.players)
      }


    return (
        <>
        <input name='gameName' placeholder="Summoner Name"
            value={gameName} onChange={e => setgameName(e.target.value)}
            onKeyDown={e => {
                if (e.key === 'Enter') {enrollUser()}
            }}
        />
        <input id="selectedRegion" placeholder='#NA1'
            value={tagLine} onChange={e => settagLine(e.target.value)}
        >
        </input>
        <button type="submit" onClick={enrollUser}>{text}</button>
        </>
    )

}

export default SummonerEnroll