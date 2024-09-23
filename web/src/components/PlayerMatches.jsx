import Icon from "./Icon"
function PlayerMatches({matches}) {
    return (
        <>
            {matches.map(match =>
                <div id='matchInfo' key={match.match_id} class={match.win ? 'winColors' : 'loseColors'}>
                    
                    <Icon id='championIcon' icon={match.champion} type='champion'/>
                    <div>
                        <div>{match.match_id}</div>
                        <div class='noWrap'>{match.champion} {match.position}</div>
                    </div>
                    <ul id='matchStats'>
                        <li>KDA {match.kills}/{match.deaths}/{match.assists}</li>
                        <li>Doubles {match.doubles}</li>
                        <li>Triples {match.triples}</li>
                        <li>Quadras {match.quadras}</li>
                        <li>Pentas {match.pentas}</li>
                        <li>Turrets {match.turrets}</li>
                        <li>Inhibs {match.inhibs}</li>
                        <li>Dragons {match.dragons}</li>
                        <li>Rifts {match.rifts}</li>
                        <li>Barons {match.barons}</li>
                        <li>Vision {match.vision_score}</li>
                        <li>Creeps {match.creep_score}</li>
                    </ul>
                </div>
                
            )}
        </>
    )
}

export default PlayerMatches