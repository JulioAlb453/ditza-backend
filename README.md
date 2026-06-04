# Ditza Backend

API REST para **Ditza**, una app móvil de microhábitos con mascota virtual, economía de puntos (coins + season points), tienda de cosméticos y ranking entre amigos que se reinicia cada 15 días.

## Stack

| Tecnología | Uso |
|---|---|
| Go 1.26.3 | Lenguaje principal |
| PostgreSQL | Base de datos |
| `lib/pq` | Driver PostgreSQL |
| `godotenv` | Variables de entorno |
| `bcrypt` (`golang.org/x/crypto`) | Hash de contraseñas |
| `google/uuid` | IDs de usuario |
| `golang-jwt/jwt` | Tokens JWT |

## Arquitectura

El proyecto sigue **arquitectura hexagonal** organizada por entidad:

```
internal/{entidad}/
  domain/          → entidades, value objects, puertos (Repository)
  data/            → modelos de persistencia + mappers
  application/     → casos de uso (Service)
  infrastructure/
    postgres/      → adaptador de base de datos
    http/          → controladores + DTOs

cmd/api/           → bootstrap (config, container, server, main)
db/schema.sql      → esquema PostgreSQL
db/README.md       → peticiones API por entidad
db/CONSULTAS-SQL.md → consultas SQL para pgAdmin / DBeaver / psql
```

### Entidades

`user`, `habit`, `habit-completion`, `pet`, `cosmetic`, `user-cosmetic`, `friendship`, `season`, `season-score`, `point-transaction`, `leaderboard`

## Requisitos previos

- Go 1.26.3 o superior
- PostgreSQL 13+
- Git

## Configuración

### 1. Clonar e instalar dependencias

```bash
git clone <url-del-repo>
cd ditza-backend
go mod download
```

### 2. Variables de entorno

Crea un archivo `.env` en la raíz del proyecto:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=ditza
DB_SSLMODE=disable
PORT=8080
LOG_DIR=logs
JWT_SECRET=cambia-este-secreto-por-uno-seguro-de-32-chars-min
JWT_EXPIRATION_HOURS=72
```

| Variable | Descripción | Default |
|---|---|---|
| `DB_HOST` | Host de PostgreSQL | — (obligatorio) |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de PostgreSQL | — (obligatorio) |
| `DB_PASSWORD` | Contraseña | — (obligatorio) |
| `DB_NAME` | Nombre de la base de datos | — (obligatorio) |
| `DB_SSLMODE` | Modo SSL | `disable` |
| `PORT` | Puerto del API | `8080` |
| `LOG_DIR` | Directorio de logs | `logs` |
| `JWT_SECRET` | Secreto para firmar tokens (mín. 32 caracteres) | — (obligatorio) |
| `JWT_EXPIRATION_HOURS` | Duración del token en horas | `72` |

### 3. Base de datos

Crea la base de datos y aplica el esquema:

```bash
psql -U postgres -c "CREATE DATABASE ditza;"
psql -U postgres -d ditza -f db/schema.sql
```

> Si ya tenías un esquema anterior (con `BIGSERIAL` en usuarios), recrea la BD o migra manualmente antes de ejecutar el script.
> Para una BD existente previa a los campos interactivos de hábitos, ejecuta `psql -U postgres -d ditza -f db/migrations/001_add_habit_interactive_fields.sql`.

## Ejecución

```bash
go run ./cmd/api
```

El servidor quedará disponible en `http://localhost:8080` (o el puerto definido en `PORT`).

Verifica que esté activo:

```bash
curl http://localhost:8080/health
```

## Endpoints

### Salud

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/health` | Estado del servidor |

### Usuario / Auth

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `POST` | `/auth/register` | No | Registrar usuario |
| `POST` | `/auth/login` | No | Iniciar sesión |
| `GET` | `/me` | Sí | Perfil del usuario autenticado |
| `GET` | `/users/search?alias=` | Sí | Buscar usuarios por alias (mín. 2 caracteres) |

**Registro / Login — body de ejemplo:**

```json
{
  "alias": "Juli",
  "email": "user@example.com",
  "password": "password123"
}
```

> En login solo se envían `email` y `password`.

**Respuesta (201 registro / 200 login):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_at": "2026-05-29T23:00:00Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "alias": "Juli",
  "email": "user@example.com"
}
```

## Autenticación JWT

Los endpoints protegidos requieren el header:

```
Authorization: Bearer <access_token>
```

El token se obtiene al registrarse o iniciar sesión. Expira según `JWT_EXPIRATION_HOURS` (default 72 h).

### Mascota

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/pet` | Sí | Obtener mascota (404 si aún no completó un hábito) |
| `PATCH` | `/pet/equip` | Sí | Equipar o desequipar cosmético del inventario |

**Equipar:**

```json
{ "cosmetic_id": 1 }
```

**Desequipar:**

```json
{ "unequip": true, "slot": "hat" }
```

### Hábitos

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/habits` | Sí | Listar hábitos activos |
| `POST` | `/habits` | Sí | Crear hábito |
| `DELETE` | `/habits/{id}` | Sí | Desactivar hábito |
| `PATCH` | `/habits/{id}/complete` | Sí | Completar hábito |

**Crear hábito — body de ejemplo:**

```json
{
  "title": "Leer",
  "description": "Leer 10 páginas antes de dormir",
  "emoji": "📚",
  "category": "estudio",
  "color": "blue",
  "frequency": "daily",
  "target_count": 10,
  "target_unit": "páginas",
  "difficulty": "medium",
  "reminder_time": "21:00"
}
```

Campos opcionales: `description`, `emoji`, `category`, `color`, `frequency`, `target_count`, `target_unit`, `difficulty`, `reminder_time`. Si solo se envía `title`, el backend aplica valores por defecto. `frequency` acepta `daily`, `weekly` o `specific_days`; `difficulty` acepta `easy`, `medium` o `hard`; `reminder_time` usa formato `HH:MM`.

### Tienda / Cosméticos

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/shop/items` | No | Catálogo de cosméticos |
| `POST` | `/shop/buy` | Sí | Comprar cosmético |
| `GET` | `/shop/inventory` | Sí | Inventario del usuario |

### Amistades

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `POST` | `/friends/request` | Sí | Enviar solicitud de amistad |
| `PATCH` | `/friends/{id}/accept` | Sí | Aceptar solicitud |
| `PATCH` | `/friends/{id}/reject` | Sí | Rechazar solicitud |
| `GET` | `/friends` | Sí | Listar amigos |
| `GET` | `/friends/pending` | Sí | Solicitudes pendientes |

### Temporada / Ranking

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/seasons/current` | No | Temporada activa |
| `GET` | `/leaderboard/friends` | Sí | Ranking entre amigos |

## Logs

Los logs se escriben en el directorio `logs/` (configurable con `LOG_DIR`) y también se muestran en la terminal:

```
logs/
├── app.log          # eventos globales
├── http.log         # peticiones HTTP
└── models/
    ├── user.log
    ├── habit.log
    └── ...          # un archivo por entidad
```

Los archivos de log están en `.gitignore`.

## Estado de implementación

| Módulo | Repositorio PostgreSQL | HTTP |
|---|---|---|
| Usuario | ✅ Implementado | ✅ |
| Hábitos | ✅ Implementado | ✅ |
| Completar hábito | ✅ Implementado | ✅ |
| Cosméticos / Tienda | ✅ Implementado | ✅ |
| Amistades | ✅ Implementado | ✅ |
| Temporada / Ranking | ✅ Implementado | ✅ |
| Mascota (pet) | ✅ Implementado | ✅ |
| Transacciones de puntos | ✅ Implementado | ❌ (indirecto vía tienda y hábitos) |
| Puntos de temporada | ✅ Implementado | ❌ (indirecto vía ranking) |

Todos los endpoints HTTP usan repositorios PostgreSQL reales. Las operaciones transaccionales (completar hábito, comprar cosmético) usan `UnitOfWork` con transacciones de base de datos.

## Estructura del proyecto

```
ditza-backend/
├── cmd/api/                    # Punto de entrada
│   ├── main.go
│   ├── config.go
│   ├── container.go
│   ├── database.go
│   └── server.go
├── db/
│   ├── schema.sql
│   ├── README.md           # API ↔ tablas
│   └── CONSULTAS-SQL.md    # Inspección de datos en PostgreSQL
├── internal/
│   ├── user/
│   ├── habit/
│   ├── habit-completion/
│   ├── pet/
│   ├── cosmetic/
│   ├── user-cosmetic/
│   ├── friendship/
│   ├── season/
│   ├── season-score/
│   ├── point-transaction/
│   ├── leaderboard/
│   └── shared/
│       ├── domain/
│       └── infrastructure/
│           ├── database/
│           ├── logger/
│           ├── monitoring/
│           ├── password/
│           ├── httpapi/
│           ├── httpserver/
│           └── postgres/       # UnitOfWork y helpers de transacción
├── .env                        # Variables locales (no commitear)
├── .gitignore
├── go.mod
└── go.sum
```

## Solución de problemas

### Puerto ocupado (Windows)

```powershell
netstat -ano | findstr ":8080"
Stop-Process -Id <PID> -Force
```

Si el proceso es un `api.exe` huérfano de `go run`, ciérralo antes de volver a levantar el servidor.

### Error 501 al registrar

Verifica que el servidor esté corriendo con el código más reciente (`go run ./cmd/api`) y que la tabla `users` exista con el esquema UUID actual.

## Licencia

Proyecto académico / MVP — consultar al equipo para uso y distribución.
