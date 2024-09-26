import Icon from "./Icon"
function PlayerMatches({matches, dataDragonVersion}) {

    return (
        <>
            {matches.map(match =>
                <div id='matchInfo' key={match.match_id} className={match.win ? 'winColors' : 'loseColors'}>
                    <Icon id='championIcon' icon={match.champion} type='champion' dataDragonVersion={dataDragonVersion}/>
                    <div id="first_box">
                        <div className='noWrap'>{match.champion} {match.position}</div>
                        <div id='match_id'>{match.match_id}</div>
                    </div>
                    <div className='horizontal_separator'></div>
                    <ul id='matchStats'>
                        <li>KDA {match.kills}/{match.deaths}/{match.assists}</li>
                        {match.doubles > 0 ? <li>Doubles {match.doubles}</li> : null}
                        {match.triples > 0 ? <li>Triples {match.triples}</li> : null}
                        {match.quadras > 0 ? <li>Quadras {match.quadras}</li> : null}
                        {match.pentas > 0 ? <li>Pentas {match.pentas}</li> : null}
                        {match.turrets > 0 ? <li>Turrets {match.turrets}</li> : null}
                        {match.inhibs > 0 ? <li>Inhibs {match.inhibs}</li> : null}
                        {match.dragons > 0 ? <li>Dragons {match.dragons}</li> : null}
                        {match.rifts > 0 ? <li>Rifts {match.rifts}</li> : null}
                        {match.barons > 0 ? <li>Barons {match.barons}</li> : null}
                        {match.vision_score > 0 ? <li>Vision {match.vision_score}</li> : null}
                        {match.creep_score > 0 ? <li>Creeps {match.creep_score}</li> : null}
                    </ul>
                </div>
                
            )}
        </>
    )
}

export default PlayerMatches