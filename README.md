# Dating Service

A rudimentary web service which could be used to power a dating app.

## Excuses made :)

1. First time using `slog` and std lib `http` for http routes/middleware. Logging statements came out a bit weird.
2. The requirements say `password` was to be returned for `/create` but I didn't do that, for security reasons.
3. I lazily made error http responses `text/plain` these should be structured in JSON, but I was tight on time.
4. Controversially, no tests. Honestly ran out of time.
5. An improvement could be made to inject a request id into the context in the logging middleware to enable request tracing.
6. I had a gotcha with age. I wanted to provide data of birth as `dd-mm-yyyy` but I couldn't get GORM to parse a short date form and work with MySQL.
    As a result, endpoints using date, must use the full form i.e. `1987-09-14T00:00:00Z`.
    Apologies if that breaks your tests.
7. The ranking and scoring mechanism is a bit rough around the edges but shows the premise I was going for. I also ran out of time to
    to rank profiles based on attractiveness but with what I had implemented, I don't think it's complex to add.
8. In reviewing the code, I didn't actually make use of interfaces at all. I would usually add interfaces, for items such as DB, services etc to aid testing.
9. It also looks like I totally missed any meaningful commentary in a bid to get this done. I hope it's not too confusing.

## Setup

This project was set up using `docker compose`. Ordinarily you can use a single command to spin up the various containers,
but I was having issues with getting my app and `migrate` to wait for MySQL to become ready.

As a result, you must run compose three times as follows:
1. `docker compose up mysql` Waiting for it to be ready for connections
2. `docker compose up migrate` This will set up the database and seed with user data
3. `docker compose up app --build` Builds the web service and runs it. 

### DB Schema

The schema is included in `migrations`. Nothing to do here as it is handled with docker compose.
You might have noticed the large transaction for the initial setup. I initially couldn't get `migrate` to work with MySQL
So I was manually running that script to MySQL to get me going.

## Usage

## The API

The API has been implemented as required, there was no mention of use of the authentication token.
One `/login` has been called, all other calls are authenticated.
The token that is returned by login should be added to subsequent calls as a header, as follows:
```
Authorization: Bearer 228b05b982fca3080f949cc76e2b4148f0145f010952f37e0f3041847b01dd11
```

(password for all accounts is `password`)

All other endpoints check for a valid token, and the endpoints that perform actions on the users account
ensure that the logged-in users id matches the account being updated.

### Endpoints
As per the original requirements, the following endpoints have been implemented:
* `/user/create`
* `/login`
* `/discover`
* `/swipe`

* An extra endpoint `/user/preferences` was added to enable a user to specify some preferences for matching purposes.

Request:
```json
{
    "userId": 2,
    "wantsChildren": false,
    "enjoysTravel": true,
    "educationLevel": "BSCH",
    "minAge": 30,
    "maxAge": 40,
    "genders": "Female"
}
```

* `userID` ID of user the preferences are being set for. Must be logged-in user.
* `educationLevel` this is very basic. Simple string matching on backend. I seeded my DB with things like:
* * BSCH, MSCH, HS, PHD, ASC
* `genders` is also basic. Intended to be a CSV separated string supporting multiple genders.

Response:

`201 Created` for a successful submission, and an appropriate alternative otherwise.

