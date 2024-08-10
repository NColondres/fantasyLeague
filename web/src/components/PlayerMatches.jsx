import Icon from "./Icon"
function PlayerMatches({matches}) {

    return (
        <>
            {matches.map(match =>
                <div key={match.match_id}>
                    
                    <div>
                        <Icon icon={match.champion} type='champion'/>
                        {match.champion} {match.position}
                    </div>
                    <div>{match.match_id}</div>
                    <ul>
                        {match.win ? <li>Win</li> : <li>Loss</li> }
                        <li>Kills {match.kills}</li>
                        <li>Deaths {match.deaths}</li>
                        <li>Assists {match.assists}</li>
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