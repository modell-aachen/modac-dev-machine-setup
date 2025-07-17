# Team Development Guidelines

This document contains our team's core principles for writing maintainable, high-quality code. These guidelines apply across all projects and programming languages.

## Core Principles

### 1. Code Quality Standards
- **Functions must be focused and small** - Each function should have a single, clear responsibility
- **Follow existing patterns exactly** - Maintain consistency with the current codebase
- **Eliminate code duplication** - Extract common patterns into helper methods and shared constants
- **Use descriptive names** - Names should clearly communicate intent and purpose
- **Only write comments that add value** - Never comment obvious code sections
- **Always run linting** - Use project linting tools to check for code quality issues before committing
- **Fix all linting issues** - Address linting problems immediately, never ignore them
- **Use proper error handling** - Never ignore error return values or exceptions

### 2. Testing Requirements (Test-Driven Development)
- **New features require tests** - No code goes to production without proper test coverage
- **Test edge cases and errors** - Don't just test the happy path
- **Use table-driven tests** - For repetitive test patterns, use structured test data
- **Create helper methods** - Reduce test code duplication through reusable utilities
- **Use descriptive test documentation** - Test methods should clearly describe what they're testing
- **Focus on behavior, not implementation** - Tests should validate outcomes, not internal mechanics
- **Do not test external libraries** - Never write tests for third-party code from GitHub or other external sources
- **Use proper test assertions** - Use testing framework assertions instead of manual output inspection
- **Test both success and failure cases** - Ensure comprehensive coverage of all code paths

### 3. Architecture Preferences
- **Simplicity over complexity** - Always prefer simpler solutions with less code
- **Direct integration** - Avoid unnecessary abstraction layers and wrapper managers
- **Minimal abstractions** - Embed functionality directly where it's used rather than creating separate managers
- **Popular libraries** - Use well-established libraries instead of custom implementations
- **Question complexity** - If a solution involves many files or complex wrappers, ask if there's a simpler approach
- **Reduce code volume** - 80% less code through direct integration is better than complex layered architectures
- **Replace functionality completely** - When implementing new systems, replace old ones entirely rather than adding flags
- **Simplify function interfaces** - Extract dependencies from context rather than passing them as parameters when possible

### 4. Development Workflow
- **Use existing tooling** - Check for existing build/test/deploy scripts before creating new ones
- **Follow project conventions** - Each project may have specific tooling patterns (e.g., Taskfile.yml, Makefile)
- **Understand before changing** - Read existing code and understand patterns before making modifications
- **Incremental improvements** - Make small, focused changes rather than large refactors
- **Ask permission before adding external libraries** - Never add third-party dependencies without explicit approval
- **Create feature branches** - Always work in feature branches, never directly on main
- **Integration over replacement** - Build on existing systems rather than completely replacing them
- **Plan backward compatibility** - Ensure changes don't break existing interfaces

### 5. Code Organization
- **Organize imports consistently** - Group standard library, third-party, and local imports
- **Use meaningful constants** - Replace magic numbers and strings with named constants
- **Logical file structure** - Group related functionality together
- **Clear separation of concerns** - Keep different responsibilities in separate modules/files

### 6. DRY Principle (Don't Repeat Yourself)
- **Extract common patterns** - Create reusable helper methods for repeated code
- **Use shared constants** - Centralize configuration values and test data
- **Create reusable utilities** - Build libraries of common functionality
- **Avoid copy-paste programming** - If you're copying code, consider abstraction

### 7. Error Handling and Robustness
- **Handle edge cases** - Consider what happens when inputs are invalid or systems fail
- **Fail fast and clearly** - Use descriptive error messages that help with debugging
- **Test error conditions** - Ensure error paths are tested and work correctly
- **Graceful degradation** - Systems should handle failures without complete breakdown

### 8. Code Review and Collaboration
- **Small, focused changes** - Make changes that are easy to review and understand
- **Short, descriptive commit messages** - Keep main message concise (< 50 chars), put detailed explanations in PR descriptions
- **Document design decisions** - Explain non-obvious choices in code or documentation
- **Be consistent with team style** - Match the existing codebase style and patterns
- **Be critical and provide counter-arguments** - Don't just say yes; challenge requests with good evidence and reasoning when appropriate
- **Use PR descriptions for detail** - Put comprehensive explanations in pull request bodies, not commit messages

## Important Reminders

- **Never create files unless absolutely necessary** - Always prefer editing existing files
- **Don't proactively create documentation** - Only create documentation when explicitly requested
- **Focus on the task at hand** - Do what has been asked; nothing more, nothing less
- **Maintain professional standards** - Never mention tools used for code generation or similar aspects

## Implementation Guidelines

When implementing these principles:

1. **Start with tests** - Write tests first to clarify requirements
2. **Keep it simple** - Choose the most straightforward approach
3. **Follow existing patterns** - Look at how similar problems are solved in the codebase
4. **Refactor incrementally** - Make small improvements over time
5. **Ask questions** - If complexity seems necessary, discuss alternatives first

These guidelines help ensure our code is maintainable, testable, and consistent across all team members and projects.

# Personal Preferences
@personal-CLAUDE.md
