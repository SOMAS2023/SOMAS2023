# SOMAS 2023
[comment]: <> (TEST)

## Official Documents
- [Rules and Implentation](./docs/Rules%20and%20Implentation.md)
- [SOMAS Base Platform User Manual](https://imperiallondon.sharepoint.com/sites/elec70071-202310/Shared%20Documents/General/basePlatformSOMAS_User_Manual.pdf?CT=1698662908166&OR=ItemsView)
- [Introduction to SOMAS Base Platform](https://imperiallondon.sharepoint.com/sites/elec70071-202310/Shared%20Documents/General/basePlatformSOMAS.pdf?CT=1699098039015&OR=ItemsView)
- [SOMAS ACW Rules](https://imperiallondon.sharepoint.com/sites/elec70071-202310/Class%20Materials/Coursework/SOMAS%20ACW%202023.pdf?CT=1699098083591&OR=ItemsView)

## Useful Links
- [Previous Year's Github](https://github.com/SOMAS2020/SOMAS2020/tree/main)

## Running code
See [Setup & Rules](./docs/SETUP.md) for requirements - EVERYONE SHOULD READ THIS DOC.

```bash
# Approach 1
go run . # Linux and macOS: Use `sudo go run .` if you encounter any "Permission denied" errors.

# Approach 2
go build # build step
./SOMAS2023 # SOMAS2023.exe if you're on Windows. Use `sudo` on Linux and macOS as Approach 1 if required.
```

### Parameters & Help
```bash
go run . --help
```

## Structure

### [`docs`](docs)
Important documents pertaining to codebase organisation, code conventions and project management. Read before writing code.
The rules can be found here [Rules and Implementation](./docs/Rules%20and%20Implementation.md)

### [`internal`](internal)
Internal SOMAS2020 packages. Most development occurs here, including client and server code.

- [`clients`](internal/clients)
Individual team code goes into the respective folders in this directory.

- [`common`](internal/common)
Common utilities, or system-wide code such as game specification etc.

- [`server`](internal/server)
Self-explanatory.

