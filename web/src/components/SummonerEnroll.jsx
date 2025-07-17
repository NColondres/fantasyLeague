import { useState } from 'react'
import { useCookies } from 'react-cookie'

function SummonerEnroll({text, setPlayersData}){

    const [gameName, setgameName] = useState('')
    const [tagLine, settagLine] = useState('#NA1')
    const [cookies] = useCookies()

  async function enrollUser(){
        if (gameName){
            const response = await fetch(`${import.meta.env.VITE_API_URL}/enroll`, {
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

            const data = await response.json() 
            switch(true){

              case response.status == 404:
                alert(data.error)
                break

              case response.status == 409:
                if (window.confirm(`${data.name} is already in lobby: ${data.lobby_id}\nClick OK to navigate to lobby`)) {
                    window.location.href=`${import.meta.env.VITE_WEB_URL}/lobby/${data.lobby_id}`
                }
                break

              case response.status >= 400:
                alert(JSON.stringify(data))
                break
            }
            if (Object.prototype.hasOwnProperty.call(cookies, "lobby_id")){
                getLobbyInfo()
            }
        
        } else {
            alert('Summoner name empty')
        }
    }
        
    async function getLobbyInfo(){

        const response = await fetch(`${import.meta.env.VITE_API_URL}/lobby`, {
            mode: 'cors',
            credentials: 'include'
          })
        
        const data = await response.json()
    
        console.log("Lobby Reponse Data: ", data)

        setPlayersData(data.players)
      }


    return (
        <div onKeyDown={e => {
            if (e.key === 'Enter') {enrollUser()}
        }}>
        <input name='gameName' placeholder="Summoner Name"
            value={gameName} onChange={e => setgameName(e.target.value)}
        />
        <input id="selectedRegion" placeholder='#NA1'
            value={tagLine} onChange={e => settagLine(e.target.value)}
            
        >
        </input>
        <button type="submit" onClick={enrollUser}>{text}</button>
        </div>
    )

}

export default SummonerEnroll
