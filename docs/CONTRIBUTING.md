## Development

Integration testing is implemented with Go's `testscript` harness.

To run the full suite:

```sh
go test ./...
```

To run unit tests only:

```sh
go test -short ./... 
```

To run only integration scripts:

```sh
go test ./internal/cmd -run TestScripts
```

If you prefer a containerized environment, use the provided VS Code
`.devcontainer`, which is configured for this workflow.
