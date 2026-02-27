# Building hiimtaylorjones.com in Go

A guided reconstruction of a personal Rails 8 website in Go — built step by step as a hands-on learning exercise. The original site is a blog/portfolio with posts, admin authentication, markdown rendering, and PostgreSQL.

## Steps

| Step | Topic | Status |
|------|-------|--------|
| [Step 1](./step-01-hello-world.md) | Initialize Go Module & Hello World Server | Complete |
| [Step 2](./step-02-routing.md) | Add Routing with Chi | Complete |
| [Step 3](./step-03-templates.md) | HTML Templates & Layouts | Complete |
| [Step 4](./step-04-database.md) | Connect to PostgreSQL | Complete |
| [Step 5](./step-05-models-migrations.md) | Data Models & Migrations | Complete |
| [Step 6](./step-06-posts-read.md) | Posts CRUD — Read Side | Complete |
| [Step 7](./step-07-posts-write.md) | Posts CRUD — Write Side | Complete |
| [Step 8](./step-08-auth.md) | Admin Authentication | In Progress |

## Tech Stack

| Concern | Rails (original) | Go (rebuild) |
|---------|-----------------|--------------|
| Router | Rails router (`config/routes.rb`) | [go-chi/chi](https://github.com/go-chi/chi) |
| Templates | ERB | `html/template` (stdlib) |
| Database driver | `pg` gem | [jackc/pgx](https://github.com/jackc/pgx) |
| Migrations | `rails db:migrate` | [pressly/goose](https://github.com/pressly/goose) |
| Markdown | Commonmarker | [yuin/goldmark](https://github.com/yuin/goldmark) |
| Auth | Devise | `golang.org/x/crypto/bcrypt` + [alexedwards/scs](https://github.com/alexedwards/scs) |
| Sessions | Rails session / cookies | scs session manager |
