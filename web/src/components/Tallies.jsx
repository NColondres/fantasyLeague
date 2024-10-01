function Tallies({type, tallies}){

    const list_items = []

    if (type === 'KDA'){
        list_items.push(<li key={type}>{tallies.kill}/{tallies.death}/{tallies.assist}</li>)
    }
    for (const [key, value] of Object.entries(tallies)){
        
        if (key === 'kill' || key === 'death' || key === 'assist'){
            continue
        } else if (key === 'vision' || key === 'cs') {
            list_items.push(<li key={key}>{value}</li>)
        } else if (value === 0) {
            continue
        } else {
            list_items.push(<li key={key}>{key} x{value}</li>)
        }
    }

    return (
        <div>
            <div id='header'>{type}</div>
            <ul id='results'>
                {list_items}
            </ul>
        </div>
    )
}
export default Tallies