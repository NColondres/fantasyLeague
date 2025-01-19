## How to Run
Make sure you're on the root directory of this project. You'll be using `docker compose` with the `docker-compose.yaml` located in this project's root directory.

### .env file required!!!
You must create a `.env` in the root directory with the following:
```
#MySQL config
MYSQL_URL=db
MYSQL_PORT=3306
MYSQL_DATABASE=lol
MYSQL_USER=fanstasyleague
MYSQL_PASSWORD=fantasyleague123
MYSQL_RANDOM_ROOT_PASSWORD=true

#Riot_API
BASE_URL=api.riotgames.com
API_KEY=<INSERT_API_KEY_FROM_RIOT_HERE>

# Services URLS
DOMAIN=localhost
VITE_DOMAIN=localhost

API_URL=http://localhost:8080
VITE_API_URL=http://localhost:8080

WEB_URL=http://localhost:4173
#Globals
MINIMUM_PLAYERS=2
DEFAULT_MATCHES=10
```
You can leave these defaults for testing locally. You must provide an `API_KEY` or else the `sentinel` and `api` services will not be able to make queries to **Riot's API**. You will get an error at runtime.

#### Run cmd
`docker compose up -d --build`

#### View logs
`docker compose logs -f`

