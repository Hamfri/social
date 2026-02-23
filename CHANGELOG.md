# Changelog

## 1.0.0 (2026-02-23)


### Features

* impl activation email sending functionality and a background task runner. ([7f7f6ff](https://github.com/Hamfri/social/commit/7f7f6ff6099cd18ef6f047fddc6299349a820c51))
* impl authentication ([68e2f8e](https://github.com/Hamfri/social/commit/68e2f8e5454cdf4ac57afcb03962b73096f1e1a1))
* impl authorization ([ad94678](https://github.com/Hamfri/social/commit/ad94678867dd3f98cad6849e06d916bf3982ca69))
* impl logic to handle CORS, graceful srv shutdown, expose srv metrics and basic tests. ([469c545](https://github.com/Hamfri/social/commit/469c5458b563bc76f787ce48366fb433badaa7a5))
* implement account activation ([2996419](https://github.com/Hamfri/social/commit/299641925b3371e0a62fb5a81abc92fce447409c))
* implement functionality to seed development db and add db query timeouts. ([a91bb76](https://github.com/Hamfri/social/commit/a91bb764b598d01cda93d1b3e3bad6dccbcf5202))
* implement pagination and filtering. ([5ebbb3d](https://github.com/Hamfri/social/commit/5ebbb3dc6df3b71a0b1df83fb0d2e9c2f9b63cf5))
* implement users (follow, unfollow, feed) endpoints and db-seeding. ([7913139](https://github.com/Hamfri/social/commit/79131399cd858269d93e933229a7e48ae08d3520))
* implement validation, http error handling, create and get posts endpoints. ([16b12dc](https://github.com/Hamfri/social/commit/16b12dc63d8aa892cbeceb4dbf65537dbc1ea970))
* implemented posts CRUD ([1bfb644](https://github.com/Hamfri/social/commit/1bfb64455f9ce6ddf940981a857f5d46ede44d67))
* optimistic concurrency control ([39a4618](https://github.com/Hamfri/social/commit/39a4618ff15a60ee3a9ef0ca458575a5fcde8f1e))
* use redis to cache frequently used data. ([946c798](https://github.com/Hamfri/social/commit/946c7989fd22483e6d1bb1b33b14bd855fea2730))


### Bug Fixes

* ident ([0b84d9a](https://github.com/Hamfri/social/commit/0b84d9a92b3fd322a32d4133bfbbc4cf82e02b52))
* remove duplicate ([7dfaedf](https://github.com/Hamfri/social/commit/7dfaedf701b27a488229b60cf28e0cd81f50ccb9))
* removed unused variables and ci script. ([1d08619](https://github.com/Hamfri/social/commit/1d0861921927796fb0ab9ef5df75bb90fc41f629))
