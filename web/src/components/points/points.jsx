// Display the default points.
// ToDo:
// Let the user change the values (within a set range) 
// so that they can customize the scoring

import { useEffect, useState } from "react"
import "./points.css"

function Points() {

    const [points, setPoints] = useState({})

    useEffect(() => {

        getDefaultPoints()

    },[])

    async function getDefaultPoints() {

        const response = await fetch(`${import.meta.env.VITE_API_URL}/points-default`, {
            method: 'GET',
            mode: 'cors',
            credentials: 'include',
        })

        console.log(`${import.meta.env.VITE_API_URL}/points-default`)
        console.log(response)

        const body = await response.json()

        console.log(body)

        setPoints(body)

    }


    return (
        <div id="points">
            <h2 id="points">Points</h2>
            <ul id='points-list'>
                {Object.entries(points).map(([key, value]) =>
                    <li id='points-list-item'><span className="points-values">{key}</span><span className="points-values">{value}</span></li>
                )}
            </ul>
        </div>
    )
    
}

export default Points
