import { useLocation } from "react-router-dom"

function NotFound(){

    const location = useLocation()
    console.log("URL path:", location.pathname)

    return (
        <>
        <div>
            <h1>404 Page not Found</h1>
        </div>
        </>
    )
}

export default NotFound