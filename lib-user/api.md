# api doc

1. /user/login/(result is **token**) GET

- query:
  - username
  - password
  - remember

2. /user/check/ GET

- query:
  - token

3. /user/add POST

- query:
  - username
  - password
  - admin
- cookie:
  - token

4. /user/edit/ POST

- query:
  - user_id
  - username
  - password
  - admin(if your token isn't admin,it is useless)
- cookie:
  - token

5. /user/del/ POST
- query:
  -user_id
- cookie
  -token
6. /user/search/

- query:
  - key(keyword of username to search)
- cookie:
  - token
