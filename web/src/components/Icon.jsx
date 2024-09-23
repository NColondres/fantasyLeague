function Icon({icon, type, id}){
    
    const path = `/img/${type}/${icon}.png`


    return (
        <img id={id} src={path} alt={type}></img>
    )
}
export default Icon