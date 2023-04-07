## Hyperpipe-backend

Until NewPipeExtractor Supports Youtube Music Browse endpoints.

## Self-Hosting

*Please see [Hyperpipe/Docker](https://codeberg.org/Hyperpipe/Docker) if hosting both the frontend and the backend*

### Docker

Run the [container](https://codeberg.org/Hyperpipe/-/packages/container/hyperpipe-backend/latest) with `HYP_PROXY` set as an env var

## Official Frontend

https://codeberg.org/Hyperpipe/Hyperpipe

## Docs

### `GET` `/explore`

### `GET` `/genres`

### `GET` `/genres/:id`

- `:id` -> `ggMPO*`

### `GET` `/charts?params=${id}&code=${code}`

### `GET` `/next/:id?queue=${queue}`

- `:id` -> `song id (same as /watch?v=:id)`
- `$queue` -> `avoid` (optional)

### `GET` `/lyrics/:id`

- `:id` -> `MPLY*`

### `GET` `/channel/:id`

- `:id` -> `UC*`

### `GET` `/next/channel/:id/:params?ct=${click}&v=${visit}`

*Use with caution*

- `:id` -> `UC*`

## License

Hyperpipe-Backend

Copyright (C) 2022-23  Shiny Nematoda

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
