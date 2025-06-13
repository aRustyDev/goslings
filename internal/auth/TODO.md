# TODOs

## `auth.go`

- [ ] `NewAuthManager` Implement Vault Store Logic
- [ ] `NewAuthManager` Implement K8s Store Logic
- [ ] `loadFromStore` Implement Error Type for Failing this
    - StoreNotFound
    - CredsNotFound
    - StoreAccessDenied
- [ ] `Authenticate` review this logic to try and make it extensible

## `auth_test.go`

- [ ] Test `error` types
- [ ] Test `store/File` store
- [ ] Test `store/Vault` store
- [ ] Test `store/K8s` store
- [ ] Test `lease/azure` provider
- [ ] Test `lease/m365` provider
- [ ] Test `lease/mde` provider
- [ ] Test `lease/d4iot` provider
- [ ] Test `lease/app` provider
- [ ] test auth behaviour with a valid token.
- [ ] test auth behaviour with a expired token.
- [ ] test auth behaviour with a credential that fails to retrieve a token.
- [ ] test auth behaviour with a credential that experiences transient errors.

## Validation Stage

### Validate Azure Authentication

- [ ] Token Acquisition in `Confidential` Environment
    - [ ] Via `Credential`
        - [ ] Via `Assertion`
        - [ ] Via `Token Provider`
        - [ ] Via `Certificate`
        - [ ] Via `Secret`
    - [ ] Via `Auth Code`
    - [ ] Via `Username Password`
    - [ ] Via `Silent`
    - [ ] Via `On Behalf of`
- [ ] Token Acquisition in `Public` Environment
    - [ ] Via `Device Code`
    - [ ] Via `Auth Code`
    - [ ] Via `Username Password`
    - [ ] Via `Silent`
    - [ ] Via `Interactive`
- [ ] Token Acquisition in `Managed` Environment
