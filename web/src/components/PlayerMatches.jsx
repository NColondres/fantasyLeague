import Icon from "./Icon"
import Results from "./Results"
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
                    <ul id='matchStats'>
                        <div className='vertical_separator'></div>
                        
                        <Results type='KDA' tallies={{
                            kill: match.kills, death: match.deaths, assist:match.assists,
                            double: match.doubles, triple: match.triples, quadra: match.quadras,
                            penta: match.pentas}}/>
                            <div className='vertical_separator'></div>

                        <Results type='Objectives' tallies={{
                            turret: match.turrets, inhib: match.inhibs, dragon: match.dragons,
                             rift: match.rifts, baron: match.barons}}/>
                            <div className='vertical_separator'></div>

                        <Results type='Vision - CS' tallies={{
                            vision: match.vision_score, cs: match.creep_score}}/>
                    </ul>
                </div>
                
            )}
        </>
    )
}

export default PlayerMatches