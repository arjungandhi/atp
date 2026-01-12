# ATP Code Review Fixes

## 🚨 Critical Issues (Must Fix)

- [x] **1. Remove hardcoded organization data** (`config/config.go:72-90`) ✅
  - ✅ Removed Pattern-Labs hardcoded defaults
  - ✅ Empty defaults now, users configure via config.toml
  - ✅ Added documentation to README.md

- [x] **2. Fix URL reconstruction logic** (`github/sync.go:292-318`) ✅
  - ✅ Now uses stored `url` label
  - ✅ Removed redundant URL building logic

- [x] **3. Fix dangerous file overwrites** (`todo/todo.go:184-212`) ✅
  - ✅ Implemented atomic write with temp file + rename pattern
  - ✅ Proper fsync to ensure data reaches disk
  - ✅ Single .bak file with proper backup management
  - ✅ Complete error handling and rollback

## ⚠️ High Priority Issues

- [ ] **4. Refactor SyncIssues() god function** (`github/sync.go:50-202`)
  - Split into focused functions
  - Extract conflict resolution
  - Separate GitHub fetch logic
  - Separate local sync logic

- [ ] **5. Fix Project/Todo relationship** (`project/project.go:12-52`)
  - Clarify ownership model
  - Make ToTodo() create new instance instead of mutating
  - Document the relationship clearly

- [ ] **6. Add concurrency protection** (`github/sync.go`)
  - Add file locking for sync operations
  - Prevent concurrent syncs
  - Protect against user edits during sync

- [x] **7. Remove debug code from production** (`github/sync.go:467-486`) ✅
  - ✅ Removed debug prints
  - ✅ Removed test-specific "test issue" logic
  - ✅ Priority sync now works for ALL GitHub todos

- [ ] **8. Fix inconsistent error messages**
  - Standardize to lowercase
  - Use consistent wrapping pattern
  - Follow Go conventions

## 📋 Medium Priority Issues

- [x] **9. Fix path expansion bug** (`cmd/atp/cli/load.go:15-33`) ✅
  - ✅ Properly expands ~ to home directory
  - ✅ Handles both ~ and ~/ patterns
  - ✅ Added proper error handling

- [ ] **10. Optimize todo iterations** (`github/sync.go:114-222`)
  - Combine multiple iterations into single pass
  - Cache reconstructGitHubURL results

- [ ] **11. Document magic numbers** (`repo/repo.go:74-111`)
  - Add constants for depth levels
  - Document expected directory structure
  - Add validation for git repositories

- [ ] **12. Fix inconsistent naming conventions**
  - Change snake_case to camelCase
  - Document naming patterns
  - Make exported/unexported decisions clear

- [ ] **13. Compile regex patterns once** (`todo/todo.go:39-109`)
  - Move regex compilation to package level
  - Consider using parser instead of regex

- [ ] **14. Make struct fields private with validation** (`todo/todo.go:13-22`)
  - Add private fields with getters/setters
  - Add validation for priority, dates, etc.

- [ ] **15. Rename unclear functions** (`project/project.go:106-128`)
  - Make function names match their actual behavior
  - Return structured data for active/done projects

## 💡 Suggestions (Nice to Have)

- [ ] **16. Add validation layer**
  - Validate phase values
  - Validate priority format
  - Validate date formats
  - Add bounds checking

- [ ] **17. Define magic string constants**
  - Create constants for "github", "In Progress", "A", etc.
  - Use throughout codebase

- [ ] **18. Add structured logging**
  - Replace fmt.Printf with proper logging
  - Add log levels
  - Make debug output controllable

- [ ] **19. Add test utilities**
  - Helper for temp ATP directories
  - Mock GitHub API responses
  - Test data generators

- [ ] **20. Document pointer usage patterns**
  - Add to CONTRIBUTING.md
  - Make consistent throughout

## Architecture Improvements

- [ ] **21. Split github/sync.go**
  - github/client.go - API client
  - github/sync_strategy.go - Conflict resolution
  - github/todo_sync.go - Todo synchronization
  - github/pr_sync.go - PR synchronization

- [ ] **22. Add abstraction layers**
  - Storage interface for file operations
  - GitHub interface for API operations

- [ ] **23. Add package documentation**
  - Document each package's purpose
  - Add examples for complex operations
