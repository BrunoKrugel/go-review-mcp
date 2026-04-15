package prompts

// Prompt represents an MCP prompt template with its metadata and arguments.
type Prompt struct {
	// Name is the unique identifier for the prompt
	Name string
	// Description explains what the prompt does
	Description string
	// Template is the prompt text with placeholders
	Template string
	// Arguments defines the parameters this prompt accepts
	Arguments []PromptArgument
}

// PromptArgument represents a single parameter for a prompt.
type PromptArgument struct {
	// Name is the argument identifier
	Name string
	// Description explains what the argument is for
	Description string
	// Required indicates if this argument must be provided
	Required bool
}

// GetAllPrompts returns all available prompts
func GetAllPrompts() []Prompt {
	return []Prompt{
		GetReviewCodePrompt(),
		GetCheckNamingPrompt(),
		GetCheckErrorHandlingPrompt(),
		GetCheckConcurrencyPrompt(),
		GetCheckTestingPrompt(),
		GetCheckInterfacesPrompt(),
	}
}

// GetReviewCodePrompt returns the comprehensive code review prompt
func GetReviewCodePrompt() Prompt {
	return Prompt{
		Name:        "review_code",
		Description: "Comprehensive Go (Golang) code review based on Google and Uber Go style guides. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go code to review",
				Required:    true,
			},
			{
				Name:        "focus",
				Description: "Specific areas to focus on (optional): naming, errors, concurrency, testing, interfaces, performance",
				Required:    false,
			},
		},
		Template: `You are an expert Go code reviewer. Review the following Go code based on the Google Go Style Guide and Uber Go Style Guide.

CODE TO REVIEW:
{{code}}

{{if focus}}FOCUS AREAS: {{focus}}{{end}}

Please review the code for:

1. **Naming Conventions**
   - MixedCaps for multi-word names (no underscores)
   - Short, consistent receiver names (1-2 chars)
   - Clear, descriptive names that reflect purpose
   - Package names: lowercase, single-word

2. **Error Handling**
   - All errors explicitly checked (no _ discards)
   - Error wrapping with %w for context
   - Errors as last return value
   - Appropriate error types (values, custom types, wrapped)
   - Use of errors.Is and errors.As

3. **Code Structure**
   - Proper use of pointers vs values
   - No pointers to interfaces
   - Zero-value mutexes (not pointers)
   - Interface compliance verification
   - Accept interfaces, return concrete types

4. **Concurrency**
   - Proper goroutine coordination
   - Channel usage and sizing (unbuffered or size 1)
   - Context usage (first parameter, not stored)
   - sync primitives used correctly
   - Proper cleanup and shutdown

5. **Best Practices**
   - defer for cleanup
   - Copy slices/maps at boundaries
   - Small, focused interfaces
   - Table-driven tests
   - Functional options pattern for configuration
   - Doc comments for exported names

6. **Common Issues**
   - Enum starting values (should start at 1)
   - Empty slices/maps (use make or var)
   - String formatting
   - Import organization

For each issue found:
- Identify the specific line or pattern
- Explain what's wrong and why
- Reference the specific style guide rule
- Provide a corrected example
- Rate severity: Critical, High, Medium, Low

End with a summary of the overall code quality and key recommendations.`,
	}
}

// GetCheckNamingPrompt returns the naming conventions review prompt
func GetCheckNamingPrompt() Prompt {
	return Prompt{
		Name:        "check_naming",
		Description: "Review Go (Golang) code naming conventions. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go code to review for naming",
				Required:    true,
			},
		},
		Template: `Review the following Go code for naming convention compliance with Google and Uber style guides.

CODE TO REVIEW:
{{code}}

Check for:

1. **Variable and Function Names**
   - MixedCaps or mixedCaps (no underscores)
   - Clear, descriptive names
   - Appropriate length for scope
   - No single-letter names except for: i, j, k (loop counters), r (reader), w (writer)

2. **Type Names**
   - Exported types: MixedCaps starting with uppercase
   - Unexported types: mixedCaps starting with lowercase
   - No underscores in type names
   - Struct names are nouns

3. **Receiver Names**
   - Consistent across all methods
   - 1-2 characters
   - Not "this", "self", or "me"
   - Reflects the type (e.g., "c" for Client, "h" for Handler)

4. **Package Names**
   - Lowercase, single word
   - No underscores or mixedCaps
   - Short, concise, evocative
   - Not plural (e.g., "user" not "users")

5. **Interface Names**
   - Single-method interfaces end in "-er" (Reader, Writer)
   - Clear, descriptive names for multi-method interfaces

6. **Constants and Variables**
   - Exported: MixedCaps
   - Unexported: mixedCaps
   - Clear purpose from name

For each naming issue:
- Identify the problematic name
- Explain why it doesn't follow conventions
- Suggest a better name
- Reference the style guide rule`,
	}
}

// GetCheckErrorHandlingPrompt returns the error handling review prompt
func GetCheckErrorHandlingPrompt() Prompt {
	return Prompt{
		Name:        "check_error_handling",
		Description: "Review Go (Golang) code error handling patterns. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go code to review for error handling",
				Required:    true,
			},
		},
		Template: `Review the following Go code for error handling compliance with Google and Uber style guides.

CODE TO REVIEW:
{{code}}

Check for:

1. **Error Checking**
   - All errors explicitly checked (no _ discards)
   - Early returns for error cases
   - No silent error suppression
   - Error checks immediately after calls

2. **Error Creation**
   - Use errors.New for simple errors
   - Use fmt.Errorf for formatted errors
   - Use custom error types when needed
   - Sentinel errors (var ErrNotFound) for expected errors

3. **Error Wrapping**
   - Use %w verb with fmt.Errorf to wrap errors
   - Add context to errors as they propagate
   - Preserve original error information
   - Don't wrap errors that should be compared directly

4. **Error Inspection**
   - Use errors.Is for sentinel error checks
   - Use errors.As for error type checks
   - Don't compare errors with ==
   - Check specific errors before generic ones

5. **Error Messages**
   - Lowercase first letter (except proper nouns)
   - No trailing punctuation
   - Add context about what failed
   - Include relevant values (IDs, names, etc.)

6. **Panic Usage**
   - Only panic for unrecoverable errors
   - Not used for normal error handling
   - Only in init() or startup for configuration errors

For each error handling issue:
- Identify the problematic pattern
- Explain the risk or problem
- Provide corrected code
- Reference the style guide rule
- Rate severity`,
	}
}

// GetCheckConcurrencyPrompt returns the concurrency review prompt
func GetCheckConcurrencyPrompt() Prompt {
	return Prompt{
		Name:        "check_concurrency",
		Description: "Review Go (Golang) code concurrency patterns. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go code to review for concurrency",
				Required:    true,
			},
		},
		Template: `Review the following Go code for concurrency compliance with Google and Uber style guides.

CODE TO REVIEW:
{{code}}

Check for:

1. **Goroutine Management**
   - Clear goroutine lifecycle (start and stop)
   - No leaked goroutines
   - Proper coordination (WaitGroup, channels, context)
   - Error handling in goroutines

2. **Channel Usage**
   - Unbuffered or size 1 (document any other size)
   - Proper closing (by sender, not receiver)
   - No sending on closed channels
   - No receiving on nil channels
   - Select statements for multiple operations

3. **Context Usage**
   - Context as first parameter
   - Not stored in structs
   - Propagated through call chain
   - Used for cancellation and deadlines
   - Checked with ctx.Done()

4. **Synchronization Primitives**
   - Mutexes: zero-value, not pointers
   - Lock/Unlock properly paired (use defer)
   - RWMutex for read-heavy workloads
   - No mixing sync primitives with channels unnecessarily
   - Atomic operations for simple counters

5. **Data Races**
   - No shared mutable state without protection
   - Proper synchronization around maps/slices
   - Copy slices/maps at boundaries
   - Use of sync.Map for concurrent map access

6. **Common Pitfalls**
   - Loop variable capture in goroutines
   - WaitGroup counter errors
   - Mutex copying
   - Goroutines blocking forever

For each concurrency issue:
- Identify the race condition or deadlock risk
- Explain the problem and potential impact
- Provide corrected code
- Reference the style guide rule
- Rate severity (data races are Critical)`,
	}
}

// GetCheckTestingPrompt returns the testing review prompt
func GetCheckTestingPrompt() Prompt {
	return Prompt{
		Name:        "check_testing",
		Description: "Review Go (Golang) test code quality. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go test code to review",
				Required:    true,
			},
		},
		Template: `Review the following Go test code for testing best practices from Google and Uber style guides.

TEST CODE TO REVIEW:
{{code}}

Check for:

1. **Test Structure**
   - Table-driven tests for multiple cases
   - Clear test names (TestXxx format)
   - Subtests with t.Run for each case
   - Test cases in same package as code
   - Helper functions use t.Helper()

2. **Test Cases**
   - Descriptive case names
   - Cover happy path and error cases
   - Test edge cases and boundaries
   - Each case is independent
   - Clear expected vs actual

3. **Assertions and Checks**
   - Clear error messages
   - Use t.Errorf, not t.Fatalf (unless setup failed)
   - Check specific error types
   - Verify all return values
   - Use appropriate comparison

4. **Test Fixtures**
   - Setup and teardown in TestMain if needed
   - Use t.TempDir() for temporary directories
   - Clean up resources with defer
   - Don't rely on external state
   - Mock external dependencies

5. **Testing Utilities**
   - Use testing.T properly
   - Leverage testify or similar if present
   - Create reusable test helpers
   - Use t.Parallel() for parallelizable tests
   - Skip tests conditionally with t.Skip()

6. **Common Issues**
   - Tests that depend on execution order
   - Not testing error cases
   - Overly complex test logic
   - Missing test coverage for critical paths
   - Flaky tests (timing, randomness)

For each testing issue:
- Identify the problem
- Explain why it's problematic
- Provide improved test code
- Reference testing best practices
- Suggest additional test cases if missing`,
	}
}

// GetCheckInterfacesPrompt returns the interface design review prompt
func GetCheckInterfacesPrompt() Prompt {
	return Prompt{
		Name:        "check_interfaces",
		Description: "Review Go (Golang) interface design and usage. Use this only for Go code.",
		Arguments: []PromptArgument{
			{
				Name:        "code",
				Description: "The Go code to review for interfaces",
				Required:    true,
			},
		},
		Template: `Review the following Go code for interface design compliance with Google and Uber style guides.

CODE TO REVIEW:
{{code}}

Check for:

1. **Interface Design**
   - Small, focused interfaces (1-3 methods)
   - Single-method interfaces named with -er suffix
   - Defined at point of use, not with implementation
   - Clear, descriptive names
   - No unnecessary interfaces

2. **Interface Usage**
   - Accept interfaces, return concrete types
   - No pointers to interfaces
   - Interface values passed by value
   - Proper nil checks for interface values
   - Compile-time interface compliance checks

3. **Interface Compliance**
   - Use var _ Interface = (*Type)(nil) pattern
   - All interface methods implemented
   - Receiver types consistent
   - Documentation for implemented interfaces

4. **Common Patterns**
   - io.Reader, io.Writer for I/O
   - fmt.Stringer for string representation
   - error interface for errors
   - context.Context for cancellation
   - Standard library interfaces used appropriately

5. **Anti-patterns**
   - Interface pollution (too many interfaces)
   - Interfaces with too many methods
   - Defining interfaces with implementations
   - Empty interfaces overused
   - Pointer to interface types

6. **Best Practices**
   - Consumer defines interface, not producer
   - Interface segregation (keep small)
   - Composition over large interfaces
   - Implicit implementation (no "implements" keyword needed)

For each interface issue:
- Identify the design problem
- Explain why it violates principles
- Suggest improved design
- Reference the style guide
- Show example of proper usage`,
	}
}
