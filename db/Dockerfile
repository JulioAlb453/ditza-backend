FROM postgres:16-alpine

# Scripts en orden alfabético: extensiones primero, luego esquema y semillas
COPY init/00-extensions.sql /docker-entrypoint-initdb.d/00-extensions.sql
COPY schema.sql /docker-entrypoint-initdb.d/01-schema.sql

EXPOSE 5432
