function Results({type, tallies}){

    return (
        <div>
            <div id='header'>{type}</div>
            <ul id='results'>
                {tallies.map((tally) => (
                tally === 0 ? null : <li>{tally}</li>
                ))}
            </ul>
        </div>
    )
}
export default Results