# Ditza — Consultas SQL en PostgreSQL

Guía para inspeccionar los datos de la base **`ditza`** desde **pgAdmin**, **DBeaver**, **DataGrip** o la consola **`psql`**.

> Para peticiones HTTP (API), ver [`README.md`](./README.md).

---

## Conexión

| Parámetro | Valor típico (`.env`) |
|---|---|
| Host | `localhost` |
| Puerto | `5432` |
| Base de datos | `ditza` |
| Usuario | valor de `DB_USER` |
| Contraseña | valor de `DB_PASSWORD` |
| SSL | `disable` |

**psql:**

```bash
psql -h localhost -p 5432 -U postgres -d ditza
```

**Cadena de conexión:**

```
postgresql://postgres:TU_PASSWORD@localhost:5432/ditza?sslmode=disable
```

---

## Índice

1. [Panel general (resumen)](#1-panel-general-resumen)
2. [users](#2-users)
3. [habits](#3-habits)
4. [habit_completions](#4-habit_completions)
5. [cosmetics](#5-cosmetics)
6. [user_cosmetics](#6-user_cosmetics)
7. [pets](#7-pets)
8. [friendships](#8-friendships)
9. [seasons](#9-seasons)
10. [season_scores](#10-season_scores)
11. [point_transactions](#11-point_transactions)
12. [Consultas cruzadas útiles](#12-consultas-cruzadas-útiles)
13. [Mantenimiento (solo desarrollo)](#13-mantenimiento-solo-desarrollo)

---

## 1. Panel general (resumen)

Ejecuta esto al abrir la base para ver cuántos registros hay en cada tabla:

```sql
SELECT 'users'              AS tabla, COUNT(*) AS registros FROM users
UNION ALL SELECT 'habits',              COUNT(*) FROM habits
UNION ALL SELECT 'habit_completions',   COUNT(*) FROM habit_completions
UNION ALL SELECT 'cosmetics',           COUNT(*) FROM cosmetics
UNION ALL SELECT 'user_cosmetics',      COUNT(*) FROM user_cosmetics
UNION ALL SELECT 'pets',                COUNT(*) FROM pets
UNION ALL SELECT 'friendships',         COUNT(*) FROM friendships
UNION ALL SELECT 'seasons',             COUNT(*) FROM seasons
UNION ALL SELECT 'season_scores',       COUNT(*) FROM season_scores
UNION ALL SELECT 'point_transactions',  COUNT(*) FROM point_transactions
ORDER BY tabla;
```

**Última actividad por usuario** (completados + compras + transacciones):

```sql
SELECT
    u.alias,
    u.email,
    u.coins,
    MAX(hc.completed_at)  AS ultimo_habito_completado,
    MAX(uc.purchased_at)  AS ultima_compra,
    MAX(pt.created_at)    AS ultima_transaccion
FROM users u
LEFT JOIN habit_completions hc ON hc.user_id = u.id
LEFT JOIN user_cosmetics uc    ON uc.user_id = u.id
LEFT JOIN point_transactions pt ON pt.user_id = u.id
GROUP BY u.id, u.alias, u.email, u.coins
ORDER BY GREATEST(
    COALESCE(MAX(hc.completed_at), '1970-01-01'),
    COALESCE(MAX(uc.purchased_at), '1970-01-01'),
    COALESCE(MAX(pt.created_at), '1970-01-01')
) DESC;
```

---

## 2. users

**Qué guarda:** cuentas de usuario (alias, email, monedas). El campo `password` es hash bcrypt — no se puede leer la contraseña en claro.

| Columna | Tipo | Notas |
|---|---|---|
| `id` | UUID | Identificador del usuario (JWT) |
| `alias` | VARCHAR(60) | Nombre visible |
| `email` | VARCHAR(255) | Único |
| `password` | TEXT | Hash bcrypt |
| `coins` | INTEGER | Monedas disponibles (≥ 0) |
| `created_at` / `updated_at` | TIMESTAMPTZ | Auditoría |

**Ver todos los usuarios (sin password):**

```sql
SELECT id, alias, email, coins, created_at, updated_at
FROM users
ORDER BY created_at DESC;
```

**Buscar por email:**

```sql
SELECT id, alias, email, coins, created_at
FROM users
WHERE email = 'juli@test.com';
```

**Buscar por alias (parcial):**

```sql
SELECT id, alias, email, coins
FROM users
WHERE alias ILIKE '%juli%';
```

**Usuarios con más monedas:**

```sql
SELECT alias, email, coins
FROM users
ORDER BY coins DESC
LIMIT 10;
```

**Usuarios registrados hoy:**

```sql
SELECT alias, email, created_at
FROM users
WHERE created_at::date = CURRENT_DATE;
```

---

## 3. habits

**Qué guarda:** hábitos creados por cada usuario. `is_active = false` significa eliminado (soft delete).

| Columna | Tipo | Notas |
|---|---|---|
| `id` | BIGSERIAL | ID del hábito |
| `user_id` | UUID | Dueño |
| `title` | VARCHAR(80) | Nombre del hábito |
| `is_active` | BOOLEAN | `true` = visible en la app |
| `current_streak` / `best_streak` | INTEGER | Racha actual y récord |
| `last_completed_date` | DATE | Último día completado (UTC) |

**Todos los hábitos activos con alias del dueño:**

```sql
SELECT
    h.id,
    u.alias,
    u.email,
    h.title,
    h.current_streak,
    h.best_streak,
    h.last_completed_date,
    h.created_at
FROM habits h
JOIN users u ON u.id = h.user_id
WHERE h.is_active = TRUE
ORDER BY u.alias, h.created_at;
```

**Hábitos de un usuario por email:**

```sql
SELECT h.*
FROM habits h
JOIN users u ON u.id = h.user_id
WHERE u.email = 'juli@test.com'
ORDER BY h.is_active DESC, h.id;
```

**Hábitos desactivados (eliminados):**

```sql
SELECT h.id, u.alias, h.title, h.updated_at
FROM habits h
JOIN users u ON u.id = h.user_id
WHERE h.is_active = FALSE;
```

**Conteo de hábitos activos por usuario:**

```sql
SELECT u.alias, COUNT(*) AS habitos_activos
FROM habits h
JOIN users u ON u.id = h.user_id
WHERE h.is_active = TRUE
GROUP BY u.id, u.alias
ORDER BY habitos_activos DESC;
```

---

## 4. habit_completions

**Qué guarda:** cada vez que un usuario completa un hábito. Máximo **1 completado por hábito por día** (`completion_date` es columna generada en UTC).

| Columna | Tipo | Notas |
|---|---|---|
| `habit_id` / `user_id` | FK | Referencias |
| `completed_at` | TIMESTAMPTZ | Momento exacto |
| `completion_date` | DATE | Día UTC (generado) |
| `note` / `emoji` | Opcional | Detalle del completado |
| `coins_awarded` / `season_points_awarded` | INTEGER | Recompensas otorgadas |

**Completados de hoy:**

```sql
SELECT
    u.alias,
    h.title,
    hc.completed_at,
    hc.coins_awarded,
    hc.season_points_awarded,
    hc.note,
    hc.emoji
FROM habit_completions hc
JOIN habits h ON h.id = hc.habit_id
JOIN users u ON u.id = hc.user_id
WHERE hc.completion_date = CURRENT_DATE
ORDER BY hc.completed_at DESC;
```

**Historial de un usuario:**

```sql
SELECT
    hc.completion_date,
    h.title,
    hc.coins_awarded,
    hc.season_points_awarded,
    hc.note
FROM habit_completions hc
JOIN habits h ON h.id = hc.habit_id
JOIN users u ON u.id = hc.user_id
WHERE u.email = 'juli@test.com'
ORDER BY hc.completion_date DESC;
```

**Completados de un hábito específico:**

```sql
SELECT completion_date, completed_at, coins_awarded, season_points_awarded
FROM habit_completions
WHERE habit_id = 1
ORDER BY completion_date DESC;
```

**Total de recompensas otorgadas por día:**

```sql
SELECT
    completion_date,
    COUNT(*)              AS completados,
    SUM(coins_awarded)    AS coins_totales,
    SUM(season_points_awarded) AS puntos_temporada
FROM habit_completions
GROUP BY completion_date
ORDER BY completion_date DESC;
```

---

## 5. cosmetics

**Qué guarda:** catálogo de la tienda (semilla: 4 ítems). Solo los activos aparecen en `GET /shop/items`.

| Columna | Valores posibles |
|---|---|
| `slot` | `hat`, `shirt`, `background`, `accessory` |
| `rarity` | `common`, `rare` |
| `is_active` | `true` = visible en tienda |

**Ver catálogo completo:**

```sql
SELECT id, name, slot, price_coins, rarity, asset_key, is_active, created_at
FROM cosmetics
ORDER BY slot, price_coins;
```

**Solo ítems disponibles en tienda:**

```sql
SELECT id, name, slot, price_coins, rarity
FROM cosmetics
WHERE is_active = TRUE
ORDER BY price_coins;
```

**Cosméticos por slot:**

```sql
SELECT slot, COUNT(*) AS cantidad, MIN(price_coins) AS precio_min, MAX(price_coins) AS precio_max
FROM cosmetics
WHERE is_active = TRUE
GROUP BY slot;
```

---

## 6. user_cosmetics

**Qué guarda:** inventario de compras (relación usuario ↔ cosmético).

**Inventario con nombres de ítems:**

```sql
SELECT
    u.alias,
    u.email,
    c.name       AS cosmetic,
    c.slot,
    c.rarity,
    uc.purchased_at
FROM user_cosmetics uc
JOIN users u     ON u.id = uc.user_id
JOIN cosmetics c ON c.id = uc.cosmetic_id
ORDER BY uc.purchased_at DESC;
```

**Inventario de un usuario:**

```sql
SELECT c.id, c.name, c.slot, c.rarity, uc.purchased_at
FROM user_cosmetics uc
JOIN cosmetics c ON c.id = uc.cosmetic_id
JOIN users u     ON u.id = uc.user_id
WHERE u.email = 'juli@test.com'
ORDER BY uc.purchased_at DESC;
```

**Cosméticos más comprados:**

```sql
SELECT c.name, c.slot, COUNT(*) AS veces_comprado
FROM user_cosmetics uc
JOIN cosmetics c ON c.id = uc.cosmetic_id
GROUP BY c.id, c.name, c.slot
ORDER BY veces_comprado DESC;
```

**Usuarios que tienen un ítem concreto:**

```sql
SELECT u.alias, u.email, uc.purchased_at
FROM user_cosmetics uc
JOIN users u ON u.id = uc.user_id
WHERE uc.cosmetic_id = 1;
```

---

## 7. pets

**Qué guarda:** mascota virtual (1 por usuario). Se crea/actualiza al completar hábitos.

| Columna | Notas |
|---|---|
| `level` / `xp` | Progresión |
| `mood` | `happy`, `neutral`, `sad`, `sleeping` |
| `equipped_*_id` | Cosméticos equipados (nullable) |

**Ver todas las mascotas:**

```sql
SELECT
    u.alias,
    p.name,
    p.level,
    p.xp,
    p.mood,
    p.last_interaction_at,
    p.updated_at
FROM pets p
JOIN users u ON u.id = p.user_id
ORDER BY p.level DESC, p.xp DESC;
```

**Mascota de un usuario con ítems equipados:**

```sql
SELECT
    u.alias,
    p.name,
    p.level,
    p.xp,
    p.mood,
    hat.name       AS gorro,
    shirt.name     AS camiseta,
    bg.name        AS fondo,
    acc.name       AS accesorio
FROM pets p
JOIN users u ON u.id = p.user_id
LEFT JOIN cosmetics hat   ON hat.id   = p.equipped_hat_id
LEFT JOIN cosmetics shirt ON shirt.id = p.equipped_shirt_id
LEFT JOIN cosmetics bg    ON bg.id    = p.equipped_background_id
LEFT JOIN cosmetics acc   ON acc.id   = p.equipped_accessory_id
WHERE u.email = 'juli@test.com';
```

**Usuarios sin mascota** (aún no completaron ningún hábito):

```sql
SELECT u.alias, u.email
FROM users u
LEFT JOIN pets p ON p.user_id = u.id
WHERE p.user_id IS NULL;
```

---

## 8. friendships

**Qué guarda:** solicitudes y amistades entre usuarios.

| `status` | Significado |
|---|---|
| `pending` | Solicitud enviada, sin responder |
| `accepted` | Amigos |
| `rejected` | Solicitud rechazada |

**Todas las relaciones con alias legibles:**

```sql
SELECT
    f.id,
    r.alias AS solicitante,
    a.alias AS destinatario,
    f.status,
    f.created_at,
    f.responded_at
FROM friendships f
JOIN users r ON r.id = f.requester_id
JOIN users a ON a.id = f.addressee_id
ORDER BY f.created_at DESC;
```

**Solicitudes pendientes:**

```sql
SELECT
    f.id,
    r.alias AS de,
    a.alias AS para,
    f.created_at
FROM friendships f
JOIN users r ON r.id = f.requester_id
JOIN users a ON a.id = f.addressee_id
WHERE f.status = 'pending'
ORDER BY f.created_at;
```

**Amigos de un usuario:**

```sql
SELECT
    CASE
        WHEN f.requester_id = u.id THEN amigo.alias
        ELSE solicitante.alias
    END AS amigo,
    f.status,
    f.responded_at
FROM friendships f
JOIN users u ON u.id IN (f.requester_id, f.addressee_id)
JOIN users solicitante ON solicitante.id = f.requester_id
JOIN users amigo        ON amigo.id = f.addressee_id
WHERE u.email = 'juli@test.com'
  AND f.status = 'accepted';
```

**Conteo por estado:**

```sql
SELECT status, COUNT(*) AS cantidad
FROM friendships
GROUP BY status;
```

---

## 9. seasons

**Qué guarda:** temporadas de 15 días. Debe haber **una sola** con `is_active = true`.

**Ver temporadas:**

```sql
SELECT
    id,
    starts_at,
    ends_at,
    is_active,
    ends_at - starts_at AS duracion,
    CASE
        WHEN NOW() BETWEEN starts_at AND ends_at THEN 'en curso'
        WHEN NOW() < starts_at THEN 'futura'
        ELSE 'finalizada'
    END AS estado_calculado,
    created_at
FROM seasons
ORDER BY starts_at DESC;
```

**Temporada activa:**

```sql
SELECT *
FROM seasons
WHERE is_active = TRUE;
```

**Verificar que solo hay una activa** (debería devolver 1 fila):

```sql
SELECT COUNT(*) AS temporadas_activas
FROM seasons
WHERE is_active = TRUE;
```

**Días restantes de la temporada actual:**

```sql
SELECT
    id,
    ends_at,
    GREATEST(0, EXTRACT(DAY FROM ends_at - NOW()))::int AS dias_restantes_aprox
FROM seasons
WHERE is_active = TRUE;
```

---

## 10. season_scores

**Qué guarda:** puntos de temporada por usuario. Se actualiza al completar hábitos.

**Ranking de la temporada activa:**

```sql
SELECT
    ss.points,
    u.alias,
    u.email,
    s.starts_at,
    s.ends_at
FROM season_scores ss
JOIN users u   ON u.id = ss.user_id
JOIN seasons s ON s.id = ss.season_id
WHERE s.is_active = TRUE
ORDER BY ss.points DESC, u.alias;
```

**Puntos de un usuario en la temporada actual:**

```sql
SELECT ss.points, ss.updated_at
FROM season_scores ss
JOIN users u   ON u.id = ss.user_id
JOIN seasons s ON s.id = ss.season_id
WHERE u.email = 'juli@test.com'
  AND s.is_active = TRUE;
```

**Usuarios con puntos > 0 sin amigos** (útil para depurar ranking):

```sql
SELECT u.alias, ss.points
FROM season_scores ss
JOIN users u ON u.id = ss.user_id
JOIN seasons s ON s.id = ss.season_id
WHERE s.is_active = TRUE
  AND ss.points > 0
ORDER BY ss.points DESC;
```

---

## 11. point_transactions

**Qué guarda:** historial de movimientos de monedas y puntos de temporada.

| `type` | Origen típico |
|---|---|
| `habit` | Completar hábito |
| `streak_bonus` | Bonus por racha |
| `note_bonus` | Bonus por nota |
| `purchase` | Compra en tienda |
| `season_reset` | Reinicio de temporada |

**Últimas transacciones (todas):**

```sql
SELECT
    pt.created_at,
    u.alias,
    pt.type,
    pt.coins_delta,
    pt.season_delta,
    pt.reference_id
FROM point_transactions pt
JOIN users u ON u.id = pt.user_id
ORDER BY pt.created_at DESC
LIMIT 50;
```

**Historial de un usuario:**

```sql
SELECT type, coins_delta, season_delta, reference_id, created_at
FROM point_transactions pt
JOIN users u ON u.id = pt.user_id
WHERE u.email = 'juli@test.com'
ORDER BY created_at DESC;
```

**Balance neto de coins por usuario** (debe coincidir con `users.coins`):

```sql
SELECT
    u.alias,
    u.coins                    AS coins_en_users,
    COALESCE(SUM(pt.coins_delta), 0) AS coins_desde_transacciones,
    u.coins - COALESCE(SUM(pt.coins_delta), 0) AS diferencia
FROM users u
LEFT JOIN point_transactions pt ON pt.user_id = u.id
GROUP BY u.id, u.alias, u.coins
HAVING u.coins <> COALESCE(SUM(pt.coins_delta), 0)
ORDER BY ABS(u.coins - COALESCE(SUM(pt.coins_delta), 0)) DESC;
```

**Resumen por tipo de transacción:**

```sql
SELECT
    type,
    COUNT(*)           AS operaciones,
    SUM(coins_delta)   AS coins_total,
    SUM(season_delta)  AS puntos_total
FROM point_transactions
GROUP BY type
ORDER BY operaciones DESC;
```

---

## 12. Consultas cruzadas útiles

**Perfil completo de un usuario** (una sola vista):

```sql
SELECT
    u.id,
    u.alias,
    u.email,
    u.coins,
    (SELECT COUNT(*) FROM habits h WHERE h.user_id = u.id AND h.is_active)     AS habitos_activos,
    (SELECT COUNT(*) FROM habit_completions hc WHERE hc.user_id = u.id)         AS total_completados,
    (SELECT COUNT(*) FROM user_cosmetics uc WHERE uc.user_id = u.id)            AS cosmeticos,
    (SELECT COUNT(*) FROM friendships f
     WHERE f.status = 'accepted'
       AND u.id IN (f.requester_id, f.addressee_id))                           AS amigos,
    p.level AS pet_nivel,
    p.mood  AS pet_estado,
    ss.points AS puntos_temporada
FROM users u
LEFT JOIN pets p ON p.user_id = u.id
LEFT JOIN season_scores ss ON ss.user_id = u.id
LEFT JOIN seasons s ON s.id = ss.season_id AND s.is_active = TRUE
WHERE u.email = 'juli@test.com';
```

**Ranking entre amigos** (equivalente a `GET /leaderboard/friends`):

```sql
WITH amigos AS (
    SELECT
        CASE
            WHEN f.requester_id = u.id THEN f.addressee_id
            ELSE f.requester_id
        END AS friend_id
    FROM friendships f
    JOIN users u ON u.email = 'juli@test.com'
    WHERE f.status = 'accepted'
),
participantes AS (
    SELECT u.id, u.alias
    FROM users u
    WHERE u.email = 'juli@test.com'
    UNION
    SELECT u.id, u.alias
    FROM users u
    JOIN amigos a ON a.friend_id = u.id
)
SELECT
    RANK() OVER (ORDER BY COALESCE(ss.points, 0) DESC, p.alias) AS rank,
    p.alias,
    COALESCE(ss.points, 0) AS season_points,
    (p.id = (SELECT id FROM users WHERE email = 'juli@test.com')) AS es_usuario_actual
FROM participantes p
LEFT JOIN season_scores ss ON ss.user_id = p.id
LEFT JOIN seasons s ON s.id = ss.season_id AND s.is_active = TRUE
ORDER BY rank;
```

**Actividad reciente del sistema** (últimas 24 h):

```sql
SELECT 'completado' AS evento, hc.completed_at AS fecha, u.alias, h.title AS detalle
FROM habit_completions hc
JOIN users u ON u.id = hc.user_id
JOIN habits h ON h.id = hc.habit_id
WHERE hc.completed_at > NOW() - INTERVAL '24 hours'

UNION ALL

SELECT 'compra', uc.purchased_at, u.alias, c.name
FROM user_cosmetics uc
JOIN users u ON u.id = uc.user_id
JOIN cosmetics c ON c.id = uc.cosmetic_id
WHERE uc.purchased_at > NOW() - INTERVAL '24 hours'

UNION ALL

SELECT 'registro', u.created_at, u.alias, u.email
FROM users u
WHERE u.created_at > NOW() - INTERVAL '24 hours'

ORDER BY fecha DESC;
```

---

## 13. Mantenimiento (solo desarrollo)

> ⚠️ Estas operaciones **borran datos**. Usar solo en entorno local.

**Vaciar todas las tablas conservando el esquema:**

```sql
TRUNCATE
    point_transactions,
    season_scores,
    habit_completions,
    user_cosmetics,
    friendships,
    habits,
    pets,
    users,
    seasons,
    cosmetics
RESTART IDENTITY CASCADE;
```

Después de truncar, vuelve a ejecutar la semilla de `schema.sql` (temporada + cosméticos) o inserta manualmente.

**Reinicio total (borra esquema y datos):**

```sql
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
-- Luego: psql -U postgres -d ditza -f schema.sql
```

**Eliminar un usuario y todo lo relacionado** (CASCADE):

```sql
DELETE FROM users WHERE email = 'juli@test.com';
```

---

## Tips en la interfaz gráfica

| Herramienta | Acción rápida |
|---|---|
| **pgAdmin** | Clic derecho en tabla → *View/Edit Data* → *All Rows* |
| **DBeaver** | Doble clic en tabla → pestaña *Data* |
| **psql** | `\dt` lista tablas · `\d users` describe columnas · `\x` modo fila expandida |

**Refrescar datos:** en pgAdmin/DBeaver pulsa el botón de refrescar o vuelve a ejecutar la consulta después de probar la API.

**Filtrar por usuario en cualquier tabla:** casi siempre puedes unir con `users` usando `user_id` o buscar el UUID con:

```sql
SELECT id FROM users WHERE email = 'juli@test.com';
```

Sustituye `'juli@test.com'` por el email que quieras inspeccionar en cualquier consulta de este documento.
