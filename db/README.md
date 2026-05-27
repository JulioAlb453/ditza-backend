# Ditza — Base de datos y peticiones API

Guía de referencia que relaciona cada **tabla de PostgreSQL** con los **endpoints HTTP** disponibles para operar sobre esas entidades.

> **Base URL por defecto:** `http://localhost:8080`  
> **Autenticación:** JWT en header `Authorization: Bearer <access_token>`

---

## Índice

1. [Preparar la base de datos](#preparar-la-base-de-datos)
2. [Diagrama de tablas](#diagrama-de-tablas)
3. [Autenticación (obtener token)](#autenticación-obtener-token)
4. [Peticiones por entidad](#peticiones-por-entidad)
   - [users](#1-users)
   - [habits](#2-habits)
   - [habit_completions](#3-habit_completions)
   - [cosmetics](#4-cosmetics)
   - [user_cosmetics](#5-user_cosmetics)
   - [friendships](#6-friendships)
   - [seasons](#7-seasons)
   - [season_scores](#8-season_scores)
   - [pets](#9-pets)
   - [point_transactions](#10-point_transactions)
   - [leaderboard](#11-leaderboard-vista)
5. [Estado de persistencia](#estado-de-persistencia)
6. [Consultas SQL útiles](#consultas-sql-útiles)

---

## Preparar la base de datos

```bash
psql -U postgres -c "CREATE DATABASE ditza;"
psql -U postgres -d ditza -f schema.sql
```

El script `schema.sql` crea todas las tablas, índices y datos semilla (temporada activa de 15 días + 4 cosméticos).

---

## Diagrama de tablas

```
users (UUID)
 ├── habits
 ├── habit_completions
 ├── pets (1:1)
 ├── user_cosmetics ──► cosmetics
 ├── friendships (requester / addressee)
 ├── season_scores ──► seasons
 └── point_transactions

seasons
cosmetics
```

| Tabla | PK | Descripción |
|---|---|---|
| `users` | `id` (UUID) | Cuenta: alias, email, password (bcrypt), coins |
| `habits` | `id` (BIGSERIAL) | Hábitos del usuario |
| `habit_completions` | `id` | Registro diario de completado + recompensas |
| `cosmetics` | `id` | Catálogo de la tienda |
| `user_cosmetics` | `(user_id, cosmetic_id)` | Inventario de compras |
| `pets` | `user_id` | Mascota virtual (1 por usuario) |
| `friendships` | `id` | Solicitudes y amistades |
| `seasons` | `id` | Temporadas de 15 días |
| `season_scores` | `(user_id, season_id)` | Puntos de temporada por usuario |
| `point_transactions` | `id` | Historial de movimientos de coins/puntos |

---

## Autenticación (obtener token)

Todos los endpoints marcados con 🔒 requieren JWT.

### Registro → tabla `users`

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "alias": "Juli",
    "email": "juli@example.com",
    "password": "password123"
  }'
```

### Login → tabla `users`

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juli@example.com",
    "password": "password123"
  }'
```

**Respuesta (registro y login):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_at": "2026-05-29T23:00:00Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "alias": "Juli",
  "email": "juli@example.com"
}
```

Guarda el `access_token` para las peticiones siguientes:

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIs..."
```

**PowerShell:**

```powershell
$auth = Invoke-RestMethod -Uri "http://localhost:8080/auth/login" `
  -Method POST `
  -Body '{"email":"juli@example.com","password":"password123"}' `
  -ContentType "application/json"
$token = $auth.access_token
```

---

## Peticiones por entidad

### 1. `users`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `POST` | `/auth/register` | No | `INSERT INTO users` |
| `POST` | `/auth/login` | No | `SELECT` + verificar password |
| `GET` | `/me` | 🔒 | `SELECT * FROM users WHERE id = ?` |

**Perfil del usuario autenticado:**

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta:**

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "alias": "Juli",
  "email": "juli@example.com",
  "coins": 0
}
```

---

### 2. `habits`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `GET` | `/habits` | 🔒 | `SELECT * FROM habits WHERE user_id = ? AND is_active = true` |
| `POST` | `/habits` | 🔒 | `INSERT INTO habits` |
| `DELETE` | `/habits/{id}` | 🔒 | `UPDATE habits SET is_active = false` |

**Listar hábitos activos:**

```bash
curl http://localhost:8080/habits \
  -H "Authorization: Bearer $TOKEN"
```

**Crear hábito:**

```bash
curl -X POST http://localhost:8080/habits \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Beber 2L de agua"}'
```

**Respuesta (201):**

```json
{
  "habit_id": 1,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Beber 2L de agua",
  "is_active": true,
  "current_streak": 0,
  "best_streak": 0
}
```

**Desactivar hábito:**

```bash
curl -X DELETE http://localhost:8080/habits/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### 3. `habit_completions`

No tiene endpoints CRUD directos. Se crea al completar un hábito (transacción que también actualiza `users`, `habits`, `season_scores`, `pets`, `point_transactions`).

| Método | Ruta | Auth | Tablas afectadas |
|---|---|---|---|
| `PATCH` | `/habits/{id}/complete` | 🔒 | `habit_completions`, `habits`, `users`, `season_scores`, `pets`, `point_transactions` |

**Completar hábito (mínimo):**

```bash
curl -X PATCH http://localhost:8080/habits/1/complete \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Con nota y emoji (bonus opcional):**

```bash
curl -X PATCH http://localhost:8080/habits/1/complete \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "note": "Me sentí con mucha energía hoy",
    "emoji": "💪"
  }'
```

**Respuesta:**

```json
{
  "coins_earned": 13,
  "season_points_earned": 13,
  "current_streak": 1,
  "wallet_coins": 13,
  "current_season_points": 13,
  "pet_level": 1,
  "pet_mood": "happy"
}
```

---

### 4. `cosmetics`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `GET` | `/shop/items` | No | `SELECT * FROM cosmetics WHERE is_active = true` |

**Catálogo de la tienda:**

```bash
curl http://localhost:8080/shop/items
```

**Respuesta (ejemplo):**

```json
[
  {
    "cosmetic_id": 1,
    "name": "Gorra Verde",
    "slot": "hat",
    "price_coins": 120,
    "rarity": "common",
    "asset_key": "hat_green_cap",
    "is_active": true
  }
]
```

**Valores válidos de `slot`:** `hat`, `shirt`, `background`, `accessory`  
**Valores válidos de `rarity`:** `common`, `rare`

---

### 5. `user_cosmetics`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `POST` | `/shop/buy` | 🔒 | `INSERT INTO user_cosmetics` + `UPDATE users` + `INSERT point_transactions` |
| `GET` | `/shop/inventory` | 🔒 | `SELECT * FROM user_cosmetics WHERE user_id = ?` |

**Comprar cosmético:**

```bash
curl -X POST http://localhost:8080/shop/buy \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cosmetic_id": 1}'
```

**Respuesta:**

```json
{
  "cosmetic_id": 1,
  "wallet_coins": 880
}
```

**Ver inventario:**

```bash
curl http://localhost:8080/shop/inventory \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta:**

```json
[
  {
    "cosmetic_id": 1,
    "purchased_at": "2026-05-26T23:00:00Z"
  }
]
```

---

### 6. `friendships`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `POST` | `/friends/request` | 🔒 | `INSERT INTO friendships` (status = pending) |
| `PATCH` | `/friends/{id}/accept` | 🔒 | `UPDATE friendships SET status = accepted` |
| `PATCH` | `/friends/{id}/reject` | 🔒 | `UPDATE friendships SET status = rejected` |
| `GET` | `/friends` | 🔒 | `SELECT * FROM friendships WHERE status = accepted` |
| `GET` | `/friends/pending` | 🔒 | `SELECT * FROM friendships WHERE status = pending` |

**Enviar solicitud de amistad:**

```bash
curl -X POST http://localhost:8080/friends/request \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"addressee_id": "uuid-del-otro-usuario"}'
```

**Aceptar solicitud:**

```bash
curl -X PATCH http://localhost:8080/friends/1/accept \
  -H "Authorization: Bearer $TOKEN"
```

**Rechazar solicitud:**

```bash
curl -X PATCH http://localhost:8080/friends/1/reject \
  -H "Authorization: Bearer $TOKEN"
```

**Listar amigos:**

```bash
curl http://localhost:8080/friends \
  -H "Authorization: Bearer $TOKEN"
```

**Solicitudes pendientes recibidas:**

```bash
curl http://localhost:8080/friends/pending \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta de amistad (ejemplo):**

```json
{
  "friendship_id": 1,
  "requester_id": "uuid-a",
  "addressee_id": "uuid-b",
  "status": "pending",
  "created_at": "2026-05-26T23:00:00Z"
}
```

---

### 7. `seasons`

| Método | Ruta | Auth | Operación SQL equivalente |
|---|---|---|---|
| `GET` | `/seasons/current` | No | `SELECT * FROM seasons WHERE is_active = true` |

**Temporada activa:**

```bash
curl http://localhost:8080/seasons/current
```

**Respuesta:**

```json
{
  "season_id": 1,
  "starts_at": "2026-05-26T23:00:00Z",
  "ends_at": "2026-06-10T23:00:00Z",
  "is_active": true
}
```

---

### 8. `season_scores`

Sin endpoint directo. Se consulta/actualiza indirectamente al:

- Completar un hábito → `PATCH /habits/{id}/complete`
- Consultar ranking → `GET /leaderboard/friends`

**Ver puntos en SQL:**

```sql
SELECT u.alias, ss.points, s.starts_at, s.ends_at
FROM season_scores ss
JOIN users u ON u.id = ss.user_id
JOIN seasons s ON s.id = ss.season_id
WHERE s.is_active = true
ORDER BY ss.points DESC;
```

---

### 9. `pets`

Sin endpoint HTTP directo. La mascota se crea/actualiza al completar hábitos.

**Ver mascota en SQL:**

```sql
SELECT p.*, u.alias
FROM pets p
JOIN users u ON u.id = p.user_id
WHERE u.email = 'juli@example.com';
```

**Campos relevantes:** `level`, `xp`, `mood` (`happy`, `neutral`, `sad`, `sleeping`), cosméticos equipados.

---

### 10. `point_transactions`

Sin endpoint HTTP directo. Se genera al completar hábitos o comprar cosméticos.

**Ver historial en SQL:**

```sql
SELECT type, coins_delta, season_delta, created_at
FROM point_transactions
WHERE user_id = '550e8400-e29b-41d4-a716-446655440000'
ORDER BY created_at DESC
LIMIT 20;
```

**Tipos válidos:** `habit`, `streak_bonus`, `note_bonus`, `purchase`, `season_reset`

---

### 11. Leaderboard (vista)

Agrega datos de `users`, `season_scores`, `friendships` y `seasons`. No es una tabla propia.

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/leaderboard/friends` | 🔒 | Ranking de amigos en la temporada activa |

```bash
curl http://localhost:8080/leaderboard/friends \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta (ejemplo):**

```json
[
  {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "alias": "Juli",
    "season_points": 130,
    "rank": 1,
    "is_current_user": true
  }
]
```

---

## Estado de persistencia

| Entidad | Tabla | Repositorio PostgreSQL | Endpoint HTTP |
|---|---|---|---|
| Usuario | `users` | ✅ | ✅ |
| Hábito | `habits` | ✅ | ✅ |
| Completado | `habit_completions` | ✅ | ✅ (vía complete) |
| Cosmético | `cosmetics` | ✅ | ✅ |
| Inventario | `user_cosmetics` | ✅ | ✅ |
| Amistad | `friendships` | ✅ | ✅ |
| Temporada | `seasons` | ✅ | ✅ |
| Puntos temporada | `season_scores` | ✅ | ❌ (indirecto) |
| Mascota | `pets` | ✅ | ❌ (indirecto) |
| Transacciones | `point_transactions` | ✅ | ❌ (indirecto) |
| Ranking | — | ✅ (consulta SQL) | ✅ |

Todos los endpoints HTTP persisten en PostgreSQL. Las operaciones que modifican varias tablas (completar hábito, comprar cosmético) se ejecutan dentro de una transacción.

---

## Consultas SQL útiles

**Usuarios registrados:**

```sql
SELECT id, alias, email, coins, created_at FROM users;
```

**Hábitos de un usuario:**

```sql
SELECT h.*
FROM habits h
JOIN users u ON u.id = h.user_id
WHERE u.email = 'juli@example.com';
```

**Completados de hoy:**

```sql
SELECT hc.*, h.title
FROM habit_completions hc
JOIN habits h ON h.id = hc.habit_id
WHERE hc.completion_date = CURRENT_DATE;
```

**Amistades aceptadas:**

```sql
SELECT f.*, r.alias AS requester, a.alias AS addressee
FROM friendships f
JOIN users r ON r.id = f.requester_id
JOIN users a ON a.id = f.addressee_id
WHERE f.status = 'accepted';
```

**Cosméticos semilla:**

```sql
SELECT id, name, slot, price_coins, rarity FROM cosmetics;
```

**Reiniciar datos (solo desarrollo):**

```sql
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
-- Luego volver a ejecutar schema.sql
```

---

## Flujo de prueba completo

```bash
# 1. Registrar
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"alias":"Juli","email":"juli@test.com","password":"password123"}' | jq .

# 2. Exportar token (Linux/macOS)
export TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"juli@test.com","password":"password123"}' | jq -r .access_token)

# 3. Perfil
curl -s http://localhost:8080/me -H "Authorization: Bearer $TOKEN" | jq .

# 4. Crear hábito
curl -s -X POST http://localhost:8080/habits \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Meditar 5 min"}' | jq .

# 5. Ver tienda
curl -s http://localhost:8080/shop/items | jq .

# 6. Temporada activa
curl -s http://localhost:8080/seasons/current | jq .
```

> Sustituye `jq` por inspección manual si no lo tienes instalado.
