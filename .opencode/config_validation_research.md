# Configuration Validation Research

## Current Problem
The current config validation in `server/internal/config/config.go` has a monolithic `validate()` function that manually checks relationships between fields. The user is concerned that as more configuration options are added, this function will become unwieldy and hard to maintain.

Current validation logic:
- Database: exactly one of Local or Turso must be configured
- Local database: requires Path and Migrations fields
- Turso database: requires URL and Token fields
- Storage: requires ReposDir
- Server: requires Host and valid Port (1-65535)

## Research Findings

### 1. go-playground/validator (Most Popular)
**Pros:**
- Most mature and widely adopted (18.7k stars)
- Extensive built-in validators
- Cross-field validation support via tags like `eqfield`, `gtfield`, etc.
- Custom validation functions
- Struct-level validation functions
- Tag-based approach keeps validation logic close to field definitions
- Good performance (benchmarked)
- i18n support for error messages

**Cons:**
- Complex tag syntax for advanced validations
- Cross-field validation limited to same struct level
- Custom validators require registration
- Learning curve for complex scenarios

**Tag Examples:**
```go
type Config struct {
    Database Database `validate:"required"`
    Storage  Storage  `validate:"required"`
    Server   Server   `validate:"required"`
}

type Database struct {
    Local *LocalDatabase `validate:"required_without=Turso"`
    Turso *TursoDatabase `validate:"required_without=Local"`
}
```

### 2. asaskevich/govalidator (Alternative)
**Pros:**
- Simpler API
- Good for basic validations
- Struct tag support
- Custom validators

**Cons:**
- Less active development
- Limited cross-field validation
- Fewer built-in validators
- Less sophisticated than go-playground/validator

### 3. Custom Validation Approaches

#### Interface-based Validation
```go
type Validator interface {
    Validate() error
}

func (c *Config) Validate() error {
    // Custom validation logic
}
```

#### Functional Validation
```go
type ValidationRule func(*Config) error

var rules = []ValidationRule{
    validateDatabase,
    validateStorage,
    validateServer,
}
```

#### Builder Pattern
```go
type ConfigValidator struct {
    rules []ValidationRule
}

func NewConfigValidator() *ConfigValidator {
    return &ConfigValidator{
        rules: []ValidationRule{
            validateExactlyOneDatabase,
            validateRequiredFields,
        },
    }
}
```

## Recommendations

### Option 1: go-playground/validator with Custom Validators (RECOMMENDED)
This approach combines the power of tag-based validation with custom logic for complex relationships.

```go
type Config struct {
    Database Database `validate:"required,database_exclusive"`
    Storage  Storage  `validate:"required"`
    Server   Server   `validate:"required"`
}

type Database struct {
    Local *LocalDatabase `validate:"omitempty,dive"`
    Turso *TursoDatabase `validate:"omitempty,dive"`
}

// Custom validator for database exclusivity
func validateDatabaseExclusive(fl validator.FieldLevel) bool {
    db := fl.Field().Interface().(Database)
    localSet := db.Local != nil
    tursoSet := db.Turso != nil
    return localSet != tursoSet // exactly one must be set
}
```

**Benefits:**
- Declarative validation close to field definitions
- Extensible with custom validators
- Handles both simple field validation and complex relationships
- Industry standard approach
- Good error messages and i18n support

### Option 2: Interface-based with Validation Registry
```go
type Validator interface {
    Validate() error
}

type ValidationRegistry struct {
    validators []Validator
}

func (c *Config) Validate() error {
    registry := &ValidationRegistry{
        validators: []Validator{
            &DatabaseValidator{c.Database},
            &StorageValidator{c.Storage},
            &ServerValidator{c.Server},
        },
    }
    return registry.ValidateAll()
}
```

**Benefits:**
- Clean separation of concerns
- Easy to test individual validators
- No external dependencies
- Full control over validation logic

### Option 3: Hybrid Approach
Use go-playground/validator for simple field validation and custom validators for complex relationships:

```go
type Config struct {
    Database Database `validate:"required"`
    Storage  Storage  `validate:"required"`
    Server   Server   `validate:"required"`
}

func (c *Config) Validate() error {
    // First run standard validation
    if err := validator.New().Struct(c); err != nil {
        return err
    }
    
    // Then run custom relationship validation
    return c.validateRelationships()
}

func (c *Config) validateRelationships() error {
    // Custom logic for database exclusivity, etc.
}
```

## Implementation Strategy

1. **Start with go-playground/validator** for basic field validation
2. **Add custom validators** for relationship constraints
3. **Keep validation logic modular** - one validator per concern
4. **Use struct tags** for simple validations
5. **Use custom functions** for complex business rules

This approach will scale well as new configuration options are added, keeps validation logic organized, and leverages battle-tested libraries while maintaining flexibility for custom requirements.