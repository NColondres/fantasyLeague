function Icon({icon, type}){
    
    const path = `/img/${type}/${icon}.png`


    return (
        <img src={path} alt={type}></img>
    )
}
export default Icon