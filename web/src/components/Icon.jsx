function Icon({icon, type, id, dataDragonVersion}){
    
    const source = `https://ddragon.leagueoflegends.com/cdn/${dataDragonVersion}/img/${type}/${icon}.png`

    return (
        <img id={id} src={source} alt={type}></img>
    )
}
export default Icon