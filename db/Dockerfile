FROM mysql:latest
COPY ./db/init.sql /docker-entrypoint-initdb.d
ENV MYSQL_DATABASE=lol
ENV MYSQL_USER=league
ENV MYSQL_PASSWORD=leagueoflegends
ENV MYSQL_RANDOM_ROOT_PASSWORD=true
EXPOSE 3306