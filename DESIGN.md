# Design Overview

## Primary Tools

**Frontend**: `Javascript`/`React`/`Webpack`/Airbnb-styleguide -- a standard frontend stack in keeping with Teleport's typical stack. For the sake of expediency I'll build for the latest stable Chrome release (and can therefore eschew `Babel` and a variety of other headaches) and will ignore mobile optimization.

**Backend**: `Go 1.15` standard library plus external dependencies (discussed below)

## Router/Dispatcher

The router/dispatcher will use the popular [gorilla/mux](https://github.com/gorilla/mux) package rather than the standard `http.ServeMux`. The primary reason for this choice is that `gorilla` provides a cousin package [gorilla/websocket](https://github.com/gorilla/websocket) for handling websockets. Technically `gorilla/websocket` should be compatible with `go`'s standard `http` package, however using the two `gorilla` packages together ensures superior online support and avoids any possible idiosyncratic incompatibilities. Additionally `gorilla` seemingly has a reputation for better developer ergonomics than the standard package.

Note: My plan is to initially implement the system with client-side polling (the level 2 project), and then time/energy permitting attempt to refactor the system to hit the level 3/4 project which would make use of websockets.

## Browser-client Security

Browser-client's will be managed via cookie-based authentication and sessions will be stored on the server. The server will use the [scs](https://github.com/alexedwards/scs) package for managing sessions, which implements a session management pattern following the [OWASP security guidelines](https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md). `scs` appears relatively popular and actively supported (813 stars, most recent commit 1 month ago at the time of writing), and offers an easily comprehensible interface for handling user sessions. Initially the session store will reside in memory, but the package supports a variety of other server-side session stores (i.e. Postgres, Redis) which should make future technology choices relatively trivial to implement as the app scales.

**Note on `gorilla/sessions`**: Because I'm using multiple other `gorilla` packages, it's natural to wonder why I'm not also using [gorilla/sessions](https://github.com/gorilla/sessions) for session management. While `gorilla/sessions` ostensibly provides extra security with its inbuilt support for encrypted/signed [secure cookies](https://github.com/gorilla/securecookie), its default `CookieStore` implementation apparently uses a [non-OWASP-compliant](https://curtisvermeeren.github.io/2018/05/13/Golang-Gorilla-Sessions.html) storage mechanism. I don't have particular reason to doubt the security of this implementation, but generally it seems wiser to go with standard, battle-tested methods. On top of that the `gorilla/sessions` pakcage doesn't provide a straight forward API for handling session timeouts like `scs` does, which is nice to have.

#### Session settings

The key security parameters for session management are the idle and absolute timeout values. Because our application has an "Upgrade" feature which would likely be extended to trigger a payment, I'm going to consider it a ["high-value application"](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#session-expiration) and set the idle timeout to 5 minutes and the absolute timeout to 1 hour. These are admittedly extremely conservative values and I imagine might be iterated on in the future based on whether or not they're annoying our users.

#### Session cookie attributes

- `Secure`: Protects against MITM attacks
- `HttpOnly`: Mitigates XSS attacks (though a motivated attacker with access to the hardware could theoretically still steal your cookie)
- `SameSite=Strict`: Mitigates CSRF attacks
- `Expires` not set: Non-persistent cookie. Session lifetimes are managed by `scs` and the session cookie contains no sensitive information, so I'm going with [standard practice](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#expire-and-max-age-attributes).
- `Domain` not set: Restrict subdomains

#### CSRF protection

Although there is some discussion on StackOverflow as to whether `SameSite=Strict` eliminates the need for CSRF tokens, the [OWASP CSRF cheat sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#samesite-cookie-attribute) doesn't mince words:

> It is important to note that this attribute should be implemented as an additional layer defense in depth concept. This attribute protects the user through the browsers supporting it, and it contains as well 2 ways to bypass it as mentioned in the following section. This attribute should not replace having a CSRF Token. Instead, it should co-exist with that token in order to protect the user in a more robust way.

Therefore I will use [`gorilla/csrf`](https://github.com/gorilla/csrf) middleware for handling CSRF tokens as well as build a custom middleware for [verifying the origin with standard headers](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#verifying-origin-with-standard-headers) as the [minimal suggested additional "in depth" mitigation](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#introduction).

## Fakeiot-client security

Given the design of the [fakeiot](https://github.com/gravitational/fakeiot) simulation where the bearer token is passed from the command line, the security model between the fakeiot-client and the server is technically a [pre-shared key](https://en.wikipedia.org/wiki/Pre-shared_key) model.

#### API key design

Interestingly I wasn't able to find any definitive standard (on OWASP or elsewhere) for "API key entropy", but using CSRF tokens as a [reasonable proxy](https://security.stackexchange.com/a/54126) I landed on a 32 byte (base64 encoded) string generated by a CSPRNG (`crypto/rand`).

#### Considerations

This is a relatively insecure security model, because if a single client or the server is compromised (or MITM'd) and the key is stolen, the whole system is compromised (not just a single client). There are a variety of methods for [hardening](https://hackernoon.com/improve-the-security-of-api-keys-v5kp3wdu) our keys with varying degrees of engineering effort required to implement.

The primary advantage of this model is that its extremely cheap to implement. MITM risk can be significantly reduced by requiring SSL connections so that the token is never mistakenly sent in an unencrypted header. And given that the fakeiot-clients are only talking to the `/metrics` endpoint, compromise is relatively low risk -- a successful attacker could try to spam new accounts or force some unfortunate account to upgrade before they'd actually hit their user limit, but these would likely be caught manually and are reversible; or they could try some species of SQL/Log injection attack, but these can be mitigated by standard methods; or they could attempt a DoS/DDoS attack, but this can be defended against via rate limiting protection.

## Database

For easy installation and usage I will use SQLite3 as the RDBMS. If this project was expected to scale up massively, I would elect to migrate over to Postgres. For additional security I might add password protection and encryption to the database file (or SSL in the case of Postgres), but for the sake of avoiding extra complexity and scope creep in this project I will simply label these as theoretical TODO's.

Database queries will be sent with the `database/sql` package which uses [prepared statements](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html#defense-option-1-prepared-statements-with-parameterized-queries) [behind the scenes](http://go-database-sql.org/prepared.html)

#### Data model

| account    |      |          |          |
| ---------- | ---- | -------- | -------- |
| account_id | plan | username | password |

| user    |            |           |
| ------- | ---------- | --------- |
| user_id | account_id | is_active |

| logins     |         |           |
| ---------- | ------- | --------- |
| account_id | user_id | timestamp |

## Endpoints

#### `/login`

**POST**: Public, CSRF protected. Valid username/pwd combo creates a new session and gives the user a corresponding session cookie.

#### `/logout`

**DELETE**: Session cookie protected. Deletes corresponding session.

#### `/metrics`

**POST**: API key protected. Updates the `logins` table with a new row. If a new `account_id` is recieved, updates the `account` table with a corresponding account on the startup plan. For each new `account_id`/`user_id` combination that's recieved, a new entry in the `user` table is created; the `is_active` column is determined by whether the corresponding account has exceeded it's plan's usage limits.

**GET**: Session cookie protected. Returns the plan's current number of active users and plan type/user-limit.

#### `/updgrade`

**PATCH**: Session cookie protected. Upgrades `account`'s `plan` column, and updates all previously inactive users on that account to active.

As a final security consideration, these endpoints might be protected by a rate limiting middleware. Unless directed otherwise I will consider this out of scope for this project.
