# Fantasy League webapp

- [Fantasy League webapp](#fantasy-league-webapp)
  - [Description](#description)
  - [App Diagram](#app-diagram)
  - [Tech Stack](#tech-stack)
  - [Deployment](#deployment)
    - [Docker with Docker Compose](#docker-with-docker-compose)
  - [Revenue](#revenue)

---
## Description
Basically I want to take what I've done in Python writing my `leagueOfGoons` discord bot and create a webapp that is public facing.

---
## App Diagram
![App Diagram](fantasyAppDesign.drawio.svg)

---
## Tech Stack
Off the top of my head from *backend to frontend*
    
1. **MySQL database**
    - Storing each user.
        - Generate a unique hash to represent a user. I do not want to deal with username and passwords.
    - Storing the lobby created and all the users associated with it.
    - Deciding the amount of games played for the lobby
    - Store the points configuration that the lobby decided on. 
        - I.e. how many points for Kills, Assists, Towers, Wins, etc...
    - Track every game played of every user in a lobby

2. **REST Api**
    - An api to interact with the MySQL database and [Riot Games API](https://developer.riotgames.com/)
    - Written in `Go` as I just learned it and figure it would perform better than Python (we'll see)
    - The api will only be reachable by the domain of my web app.
        - I believe this is done by `CORS` header. Have to look into how to set that up with whatever Go API package I use.
            - `Gin` looks like an option.

3. **NGINX webserver**
    - Will used to host my website

4. **Vue.js webapp**
    - Webapp for the users to interface
    - Every user will need to enter their `region` and `league account` name at some point.
    - Make it look pretty?
    - How I think it will work
        - Every user visiting the page will generate a `user UID` cookie.
        - If creating a lobby, the API will create a `lobby UID` (MySQL DB) and the *creater/owner*'s `user UID`.
            - User is directed to a `webapp.com/lobby/<lobby UID>` where the URL can be shared with others to join.
                - The *creator/owner* of the lobby can set the points and settings for pre-defined metrics available (add more over time).
                - Track users as they join and once they have the people they want in the lobby, can start the lobby.
        - Users joining a lobby will generate a `user UID` in their cookies which will show them the status of the lobby.
        - Once lobby is started, another backend service will be periodically checking league accounts registered in the DB and updating the DB with games played.
        - Announce the **WINNER!** in some pretty way

---
## Deployment
It would be nice to use the cloud, but I don't want to pay for anything.

I'll probably use a couple of the Raspberry Pi's. Maybe this could lead into a dedicated big boy server (rack mounted) or moving onto the cloud if generating enough revenue.

### Docker with Docker Compose
Kubernetes might be overkill but I do want to containerize everything to make it portable. That way if we need to scale vertically we can.

If the app becomes bigger we can look into horizontal scaling with kuberentes. But if every service is an already built image, migration should be ***less*** painful.

---
## Revenue
We can throw advertisements on the page. Not sure how much money that brings in, would have to investigate