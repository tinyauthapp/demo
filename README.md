# Tinyauth Demo

A small web app for demonstrating [Tinyauth](https://github.com/tinyauthapp/tinyauth).

It sits behind Tinyauth via Traefik forward auth and does zero authentication
of its own. It serves two pages:

- **`/` - public.** Tinyauth is told to skip authentication for this exact path
  (through ACLs), so anyone can open it. It explains
  what will happen when you click through to the protected page.
- **`/protected` - behind Tinyauth.** Opening it without a session bounces you
  to the Tinyauth login page. Once signed in it shows a "You're authenticated!"
  hero, the `Remote-User` / `Remote-Name` / `Remote-Email` / `Remote-Groups`
  headers the proxy injected, and a "Log out" button pointing at Tinyauth's
  logout endpoint.

Both pages include a collapsible table of every request header, so you can
compare exactly what the proxy forwarded in each case.

## Run the full demo stack

1. Point `app.example.com` and `auth.example.com` at the Docker host. For a
   local test, add to `/etc/hosts`:

   ```
   127.0.0.1 app.example.com auth.example.com
   ```

2. Start everything:

   ```sh
   docker compose up -d --build
   ```

3. Open <http://app.example.com>. You get the public landing page, no login,
   no identity headers.

4. Click **Open protected page**. Traefik asks Tinyauth whether you're logged
   in, Tinyauth redirects you to its login page, and after signing in
   (`demo` / `password`) you land on `/protected` showing your identity.

## Configuration

| Variable     | Default   | Description                                  |
| ------------ | --------- | -------------------------------------------- |
| `LOGOUT_URL` | `/logout` | Target of the "Log out" button.              |

The server always listens on `:3000` and serves `/` (public) and `/protected`.
