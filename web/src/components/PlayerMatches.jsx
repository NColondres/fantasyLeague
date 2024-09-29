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
                        
                        <Results type='KDA' tallies={[
                            `${match.kills}/${match.deaths}/${match.assists}`, 
                            match.doubles, match.triples, match.triples, match.quadras, match.pentas]}/>
                            <div className='vertical_separator'></div>

                        <Results type='Objectives' tallies={[
                            match.turrets, match.inhibs, match.dragons, match.rifts, match.barons]}/>
                            <div className='vertical_separator'></div>

                        <Results type='Vision - CS' tallies={[
                            match.vision_score, match.creep_score]}/>
                    </ul>
                </div>
                
            )}
        </>
    )
}

export default PlayerMatches