# Claude Development Guidelines for Lookie

## Project Overview
Lookie is a quantum computing intelligence service that monitors industry news, classifies content using AI, and provides alerts for important developments. Built in Go with SQLite, Gemini AI integration, and REST API.

## Development Approach

### Research-First Methodology
- **Industry Research**: Study similar products (news aggregators, intelligence services, monitoring tools) before implementing
- **Best Practices**: Research Go patterns, API design standards, and monitoring service architectures
- **Benchmarking**: Learn from established products in adjacent spaces (news aggregation, competitive intelligence, etc.)

### Epic-Driven Development
- Break complex work into logical epics using TodoWrite tool
- Each epic should be a cohesive unit of functionality
- Regular checkpoints between epics to validate direction
- Track progress and dependencies clearly

### Quality Assurance Process
1. **Pre-Implementation**: Research patterns and validate approach
2. **During Development**: Follow established Go conventions from existing codebase
3. **Post-Implementation**: Review against requirements and industry standards
4. **Testing**: Validate functionality works as expected

### Communication Pattern
- **Research Findings**: Present key insights before implementation
- **Epic Planning**: Break down work and get approval before starting
- **Checkpoint Reviews**: Pause at milestones for feedback
- **Technical Decisions**: Explain rationale based on research

## Technical Standards

### Code Quality
- Follow existing Go conventions in codebase
- Maintain consistency with current architecture patterns
- Use structured logging (already established)
- Handle errors appropriately with context

### Dependencies
- Prefer existing libraries already in go.mod
- Research thoroughly before adding new dependencies
- Consider maintenance and security implications

### Testing Strategy
- Research existing test patterns in codebase
- Ensure new features have appropriate test coverage
- Validate against real-world usage patterns

## Collaboration Expectations

### From Human Partner
- Define vision and requirements clearly
- Set priorities and constraints
- Provide domain expertise in quantum computing industry
- Review epic breakdowns and research findings
- Give feedback at checkpoints

### From Claude
- Thorough research before implementation
- Clear communication of findings and decisions
- Structured epic planning
- Regular progress updates
- Pattern validation against industry standards

## Current Project Status
- **Language**: Go 1.24.5
- **Database**: SQLite with WAL mode  
- **AI Integration**: Google Gemini API
- **Architecture**: Microservice-style with scheduler, scraper, classifier, and API components
- **Deployment**: Docker-ready, supports multiple run modes

## Commands for Development
- **Build**: `go build -o bin/lookie cmd/lookie/main.go`
- **Test**: `go test ./...`
- **Lint**: Check existing patterns (no standardized linter configured yet)
- **Migration**: `go run cmd/migrate/main.go`