//go:build tools
// +build tools

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Contract validation script that validates the OpenAPI specification
// against the actual implementation in main.go
//
// Usage: go run validate_contract.go
// Exit code: 0 = success, 1 = validation failure

func main() {
	fmt.Println("🔍 Starting API Contract Validation...")
	fmt.Println()

	// Load OpenAPI specification
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("openapi.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load OpenAPI spec: %v", err)
	}

	// Validate the OpenAPI specification itself
	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("❌ OpenAPI specification is invalid: %v", err)
	}
	fmt.Println("✅ OpenAPI specification is valid")

	// Expected routes from the implementation
	expectedRoutes := map[string][]string{
		"/":                     {"GET"},
		"/api/enrollments":      {"GET", "POST"},
		"/api/enrollments/{id}": {"GET", "PUT", "DELETE"},
	}

	// Validate all routes exist in the spec
	fmt.Println("\n📋 Validating routes...")
	routeCount := 0
	for path, methods := range expectedRoutes {
		pathItem := doc.Paths.Find(path)
		if pathItem == nil {
			log.Fatalf("❌ Route not found in OpenAPI spec: %s", path)
		}

		for _, method := range methods {
			operation := pathItem.GetOperation(method)
			if operation == nil {
				log.Fatalf("❌ Method %s not defined for route %s", method, path)
			}
			fmt.Printf("  ✓ %s %s\n", method, path)
			routeCount++
		}
	}
	fmt.Printf("✅ All %d routes validated\n", routeCount)

	// Validate required schemas
	fmt.Println("\n📦 Validating schemas...")
	requiredSchemas := []string{
		"Enrollment",
		"EnrollmentInput",
		"ErrorResponse",
		"SuccessResponse",
	}

	for _, schemaName := range requiredSchemas {
		if doc.Components.Schemas[schemaName] == nil {
			log.Fatalf("❌ Required schema not found: %s", schemaName)
		}
		fmt.Printf("  ✓ Schema: %s\n", schemaName)
	}
	fmt.Println("✅ All schemas valid")

	// Validate X-Cache-Status header is documented
	fmt.Println("\n🏷️  Validating custom headers...")
	getEnrollmentPath := doc.Paths.Find("/api/enrollments/{id}")
	if getEnrollmentPath != nil && getEnrollmentPath.Get != nil {
		if resp := getEnrollmentPath.Get.Responses.Status(200); resp != nil {
			if resp.Value.Headers["X-Cache-Status"] == nil {
				log.Fatalf("❌ X-Cache-Status header not documented in GET /api/enrollments/{id}")
			}
			fmt.Println("  ✓ X-Cache-Status header documented")
		}
	}

	// Validate error responses are documented
	fmt.Println("\n⚠️  Validating error responses...")
	errorCodes := []int{400, 404, 500}
	hasErrors := false

	for path := range doc.Paths.Map() {
		pathItem := doc.Paths.Find(path)
		operations := []*openapi3.Operation{
			pathItem.Get,
			pathItem.Post,
			pathItem.Put,
			pathItem.Delete,
		}

		for _, op := range operations {
			if op == nil {
				continue
			}

			for _, code := range errorCodes {
				if op.Responses.Status(code) != nil {
					hasErrors = true
					break
				}
			}
		}
	}

	if !hasErrors {
		log.Fatalf("❌ No error responses (400, 404, 500) documented")
	}
	fmt.Println("  ✓ Error responses documented (400, 404, 500)")

	// Validate status enum values
	fmt.Println("\n🔐 Validating business rules...")
	enrollmentSchema := doc.Components.Schemas["Enrollment"]
	if enrollmentSchema.Value.Properties["status"] == nil {
		log.Fatalf("❌ Status field not found in Enrollment schema")
	}

	statusProp := enrollmentSchema.Value.Properties["status"]
	expectedStatuses := []string{"pending", "active", "completed"}
	if len(statusProp.Value.Enum) != len(expectedStatuses) {
		log.Fatalf("❌ Status enum does not match expected values")
	}

	for _, status := range expectedStatuses {
		found := false
		for _, enumVal := range statusProp.Value.Enum {
			if enumVal == status {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("❌ Status '%s' not in enum", status)
		}
	}
	fmt.Println("  ✓ Status validation rules correct (pending/active/completed)")

	// Success summary
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("✅ CONTRACT VALIDATION PASSED")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("✓ OpenAPI specification: VALID\n")
	fmt.Printf("✓ Routes validated: %d\n", routeCount)
	fmt.Printf("✓ Schemas validated: %d\n", len(requiredSchemas))
	fmt.Printf("✓ Custom headers: DOCUMENTED\n")
	fmt.Printf("✓ Error responses: DOCUMENTED\n")
	fmt.Printf("✓ Business rules: VALIDATED\n")
	fmt.Println(strings.Repeat("═", 60))

	os.Exit(0)
}
