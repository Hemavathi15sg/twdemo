---
applyTo: '**'
---

# Go Clean Architecture Guidelines

## Project Structure
- Follow Clean Architecture principles with clear separation of concerns
- Organize code into layers: `models`, `repositories`, `usecases`, `handlers`, `delivery`
- Keep dependencies pointing inward (handlers -> usecases -> repositories -> models)

## Models Layer
- Define domain entities in `models/` directory
- Use appropriate Go data types:
  - `string` for text fields
  - `int64`, `int32`, `int` for numeric IDs and counters
  - `float64`, `float32` for decimal values
  - `bool` for boolean flags
  - `time.Time` for timestamps
  - `uuid.UUID` for unique identifiers
  - Pointers (`*Type`) for optional/nullable fields
- Add JSON tags for serialization: `json:"field_name"`
- Add validation tags where needed: `validate:"required,email"`
- Include struct tags for database mapping: `db:"column_name"`

## Repository Layer
- Create interfaces in `repositories/` for data access operations
- Follow naming convention: `ModelNameRepository` interface
- Implement CRUD operations: `Create`, `GetByID`, `Update`, `Delete`, `List`
- Return `(result, error)` tuple from all methods
- Use context as first parameter: `ctx context.Context`
- Keep repositories focused on single entity/aggregate

## Use Cases/Services Layer
- Business logic resides in `usecases/` or `services/`
- Depend on repository interfaces, not implementations
- Handle transactions and complex operations
- Return domain models, not DTOs at this layer

## Handlers/Delivery Layer
- HTTP handlers in `handlers/` or `delivery/http/`
- Parse request, call use case, format response
- Handle HTTP-specific concerns (status codes, headers)
- Use DTOs/request-response structs when needed

## Testing Requirements

-agents should always show testing commands to test after code generation.


## General Go Best Practices
- Use meaningful variable and function names
- Handle errors explicitly, never ignore them
- Use `defer` for cleanup operations
- Prefer composition over inheritance
- Keep functions small and focused
- Use Go modules for dependency management
- Follow standard Go formatting (`gofmt`, `goimports`)