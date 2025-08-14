# orm2hex

`orm2hex` is a Golang code generator that automatically creates repository interfaces and CRUD implementations for your structs following the **Hexagonal Architecture** pattern. It supports both **GORM** and **Bun** ORMs, with an option to generate **raw SQL** queries instead of ORM-specific methods. The tool can optionally include `context.Context` in generated methods.

---

## Features

- Automatically scans your Go project to detect all structs.
- Generates repository interfaces and concrete implementations.
- Supports **GORM** and **Bun** ORMs.
- Optionally generates **raw SQL** for all CRUD operations.
- Optionally adds `context.Context` to method signatures.
- Uses templates for easy customization.
- Generates files in a structured way and avoids overwriting logic inadvertently.

---

## Installation

You can install `orm2hex` locally without publishing to a repository:

```bash
git clone github.com/rnikrozoft/orm2hex/cmd/orm2hex
cd orm2hex
go install ./...
```
Or

```bash
go install github.com/rnikrozoft/orm2hex/cmd/orm2hex@<version>
```

This installs the `orm2hex` CLI tool into your `$GOPATH/bin` (or Go bin path).

---

## Usage

```bash
orm2hex [flags]
```

### Flags

| Flag      | Type    | Default       | Description                                           |
|-----------|---------|---------------|-------------------------------------------------------|
| `-out`    | string  | `./repository`| Output directory for generated repository files       |
| `-ctx`    | bool    | false         | Include `context.Context` in method signatures       |
| `-raw`    | bool    | false         | Generate raw SQL queries instead of ORM functions    |
| `-orm`    | string  | "gorm"        | Select ORM: `"gorm"` or `"bun"`                      |

### Example

Generate repositories for all structs in the current project using GORM and context:

```bash
orm2hex -out ./internal/repository -ctx=true -raw=false -orm=gorm
```

Generate repositories using Bun with raw SQL (no ORM helper functions):

```bash
orm2hex -out ./internal/repository -ctx=true -raw=true -orm=bun
```

---

## Output

For a struct like:

```go
type User struct {
    ID   int    `primaryKey`
    Name string
    Age  int
}
```

`orm2hex` generates:

- `UserRepository` interface with CRUD methods.
- Concrete implementation using selected ORM (`gorm` or `bun`) or raw SQL.
- Methods optionally include `context.Context` if `-ctx` flag is set.

---

## Template Customization

All templates are located in `templates/` directory:

- `raw.tmpl` – generates repositories using raw SQL.
- `helper.tmpl` – generates repositories using ORM helper functions.

You can modify these templates to match your project’s coding conventions or add custom logic.

---

## Notes

- `orm2hex` **overwrites generated repository files** each time it runs, so avoid manual edits directly in generated files. If custom logic is needed, consider wrapping or extending the generated repositories.
- The tool requires your project to be compilable and Go packages to be loadable.

---