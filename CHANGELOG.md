# Changelog

## [1.5.1](https://github.com/Hamfri/social/compare/v1.5.0...v1.5.1) (2026-02-27)


### Bug Fixes

* gcloud auth issue ([2d641d6](https://github.com/Hamfri/social/commit/2d641d6a984e9c2eda0a06e5065918ab85a6c8ee))

## [1.5.0](https://github.com/Hamfri/social/compare/v1.4.6...v1.5.0) (2026-02-27)


### Features

* Trigger deployments based on tags ([676ff50](https://github.com/Hamfri/social/commit/676ff504d3e5f4001430ec199078cdce5ceb8ea1))


### Bug Fixes

* Typo and format yml files ([78b5c43](https://github.com/Hamfri/social/commit/78b5c43342f702a0268780acd4ac6cf0397be67c))

## [1.4.6](https://github.com/Hamfri/social/compare/v1.4.5...v1.4.6) (2026-02-26)


### Bug Fixes

* update models and swagger docs. ([fba770d](https://github.com/Hamfri/social/commit/fba770d8416865cacb9706d9c2beed905cab9a94))

## [1.4.5](https://github.com/Hamfri/social/compare/v1.4.4...v1.4.5) (2026-02-26)


### Bug Fixes

* Docker image issues, allow ipv6 connections from container and fix bug in migration. ([c7e758f](https://github.com/Hamfri/social/commit/c7e758f1626199ed2d6b84df03f9e94a6103a14e))

## [1.4.4](https://github.com/Hamfri/social/compare/v1.4.3...v1.4.4) (2026-02-24)


### Bug Fixes

* Add SSL certs. ([a1eb98e](https://github.com/Hamfri/social/commit/a1eb98e78584be9d29639d99da39174bc244a72a))

## [1.4.3](https://github.com/Hamfri/social/compare/v1.4.2...v1.4.3) (2026-02-24)


### Bug Fixes

* change port ([4737a7a](https://github.com/Hamfri/social/commit/4737a7a691abb289c95b07058bbe5ca35c62a4d4))

## [1.4.2](https://github.com/Hamfri/social/compare/v1.4.1...v1.4.2) (2026-02-24)


### Bug Fixes

* remove user ([0951a15](https://github.com/Hamfri/social/commit/0951a15be3f6f77998bb2b5b32e2c0f323b11e7d))

## [1.4.1](https://github.com/Hamfri/social/compare/v1.4.0...v1.4.1) (2026-02-24)


### Bug Fixes

* build docs before building the app ([e07813b](https://github.com/Hamfri/social/commit/e07813b7cc4c21a00623130267d78a8f1b0e3252))

## [1.4.0](https://github.com/Hamfri/social/compare/v1.3.0...v1.4.0) (2026-02-24)


### Features

* run migrations before starting the application ([db335e6](https://github.com/Hamfri/social/commit/db335e6b9fb0058163a42af82b431275becd4649))

## [1.3.0](https://github.com/Hamfri/social/compare/v1.2.0...v1.3.0) (2026-02-24)


### Features

* Expose port ([b593010](https://github.com/Hamfri/social/commit/b593010e5c59f84c435f9781de47a80b1d943abe))

## [1.2.0](https://github.com/Hamfri/social/compare/v1.1.2...v1.2.0) (2026-02-24)


### Features

* Add Dockerfile ([6f5a4ea](https://github.com/Hamfri/social/commit/6f5a4ea26c419ee5493a5f72e0b712b88a4474c3))

## [1.1.2](https://github.com/Hamfri/social/compare/v1.1.1...v1.1.2) (2026-02-24)


### Bug Fixes

* read version from release PR title. ([e472143](https://github.com/Hamfri/social/commit/e4721433d42567b4cc3e017970f05bf6c844ba7a))

## [1.1.1](https://github.com/Hamfri/social/compare/v1.1.0...v1.1.1) (2026-02-24)


### Bug Fixes

* failing version bump script. ([06cf424](https://github.com/Hamfri/social/commit/06cf4243d44d3a83ebf49e700689c01484d6c54c))

## [1.1.0](https://github.com/Hamfri/social/compare/v1.0.1...v1.1.0) (2026-02-24)


### Features

* trigger release ([e01ac50](https://github.com/Hamfri/social/commit/e01ac50003cac6ca96ebd73167e9164148d602be))

## [1.0.1](https://github.com/Hamfri/social/compare/v1.0.0...v1.0.1) (2026-02-23)


### Bug Fixes

* bug in version bump script ([c63c54f](https://github.com/Hamfri/social/commit/c63c54f845b66b21e156c8cc53be5c080838e96b))

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
