# Spec Workflow Override

## When Specification Already Exists

**CRITICAL**: When a user provides an existing specification file (e.g., SPECIFICATION.MD), DO NOT create new spec documents. Instead:

1. **Read and understand the existing specification**
2. **Implement the code directly** based on the provided specification
3. **Follow the specification's requirements exactly** as written
4. **Do not create requirements.md, design.md, or tasks.md files** when a specification already exists

## Modified Workflow Requirements (Only for New Specs)

The spec workflow should proceed automatically through all phases without requiring explicit user approval at each step. The agent should:

1. **Requirements Phase**: Create comprehensive requirements document and proceed directly to design
2. **Design Phase**: Create detailed design document and proceed directly to tasks  
3. **Tasks Phase**: Create implementation task list and complete the workflow

## Removed Requirements

- **NO explicit user approval required** between phases
- **NO userInput tool calls** for spec document reviews
- **NO waiting for user confirmation** before proceeding to next phase

## Workflow Behavior

**For existing specifications**: Implement code directly based on the provided specification.

**For new specs only**: The agent should create all three spec documents (requirements.md, design.md, tasks.md) in sequence and inform the user when the complete spec is ready for implementation.