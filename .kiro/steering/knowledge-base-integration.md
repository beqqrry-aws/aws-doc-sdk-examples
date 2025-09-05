# Knowledge Base Integration for Code Example Development

## 🚨 CRITICAL REQUIREMENT 🚨
**BEFORE CREATING ANY CODE IN ANY LANGUAGE, YOU MUST:**
1. Search the Quality Code Examples Knowledge Base for "[language] code example standards"
2. Example: `search("dotnet C# .NET code example standards")` for .NET projects
3. Example: `search("Java code example standards")` for Java projects
4. Example: `search("Python code example standards")` for Python projects

**FAILURE TO DO THIS WILL RESULT IN INCORRECT CODE STRUCTURE AND REJECTED WORK**

## Overview
When developing AWS SDK code examples, agents should leverage the available knowledge base resources to ensure accuracy, consistency, and adherence to best practices. This repository has access to multiple knowledge bases that contain curated information about AWS services, implementation patterns, and proven code examples.

## Available Knowledge Base Resources

### 1. Tejas Knowledge Base (Primary AWS Service Reference)
- **Tool**: `query_tejas_kb`
- **Purpose**: Comprehensive AWS service documentation, API references, and implementation guidance
- **Usage**: Query for AWS service-specific information, API details, parameter requirements, and service capabilities
- **Auto-approved**: Yes - can be used automatically without user confirmation

### 2. Quality Code Examples Knowledge Base (MANDATORY for Language Structure)
- **Tool**: `search` (from Quality Code Examples Knowledge Base MCP server)
- **Purpose**: Search existing code examples, patterns, and implementations within this repository
- **Usage**: Find similar implementations, established patterns, and proven code structures
- **Auto-approved**: Yes - for searching existing examples and patterns
- **CRITICAL**: This MUST be consulted before creating any language-specific code structure

## Mandatory Knowledge Base Consultation Workflow

### Before Creating Any AWS Service Code Example:

1. **Query Tejas KB for Service Understanding**
   ```
   Use query_tejas_kb to research:
   - Service overview and core concepts
   - Key API operations and methods
   - Required parameters and optional configurations
   - Common use cases and implementation patterns
   - Service-specific best practices and limitations
   ```

2. **MANDATORY: Search Quality Code Examples Knowledge Base for Language Structure**
   ```
   REQUIRED: Use search to find:
   - Language-specific code example standards and patterns
   - Similar service implementations in the target language
   - Established code patterns and structures
   - Error handling approaches
   - Testing methodologies
   - Documentation formats
   - Language level standards and best practices
   
   YOU MUST SEARCH FOR "[language] code example standards" BEFORE CREATING ANY CODE
   ```

3. **Cross-Reference Implementation Approaches**
   ```
   Compare findings from both knowledge bases to:
   - Identify the most appropriate implementation approach
   - Ensure consistency with existing repository patterns
   - Validate service-specific requirements
   - Confirm best practice adherence
   ```

## Required Knowledge Base Queries by Development Phase

### Phase 1: Service Research and Planning
**Always query Tejas KB first with questions like:**
- "What is [AWS Service] and what are its primary use cases?"
- "What are the key API operations for [AWS Service]?"
- "What are the required parameters for [specific operation]?"
- "What are common implementation patterns for [AWS Service]?"
- "What are the best practices for [AWS Service] error handling?"

### Phase 2: Implementation Pattern Discovery (MANDATORY QUALITY CODE EXAMPLES KB SEARCH)
**REQUIRED: Search Quality Code Examples Knowledge Base for:**
- "[language] code example standards" - THIS IS MANDATORY
- Existing implementations of the same service in other languages
- Similar service patterns within the target language
- Established error handling and testing approaches
- Documentation and metadata patterns
- Language-specific directory structure requirements

**FAILURE TO SEARCH FOR LANGUAGE STANDARDS WILL RESULT IN INCORRECT CODE STRUCTURE**

### Phase 3: Code Structure Validation
**Use both knowledge bases to verify:**
- Implementation aligns with AWS service capabilities
- Code structure follows repository conventions
- Error handling covers service-specific scenarios
- Testing approach matches established patterns

## Knowledge Base Query Examples

### Service Overview Queries
```
query_tejas_kb("What is Amazon S3 and what are its core features?")
query_tejas_kb("What are the main DynamoDB operations for CRUD functionality?")
query_tejas_kb("What authentication methods does Lambda support?")
```

### Implementation-Specific Queries
```
query_tejas_kb("How do I configure S3 bucket policies programmatically?")
query_tejas_kb("What are the required parameters for DynamoDB PutItem operation?")
query_tejas_kb("How do I handle pagination in EC2 DescribeInstances?")
```

### Best Practices Queries
```
query_tejas_kb("What are S3 security best practices for SDK implementations?")
query_tejas_kb("How should I handle DynamoDB throttling in production code?")
query_tejas_kb("What are Lambda function timeout considerations?")
```

### Local Pattern Discovery (MANDATORY LANGUAGE STANDARDS SEARCH)
```
search("[language] code example standards") - REQUIRED FIRST SEARCH
search("S3 bucket creation examples in Python")
search("DynamoDB error handling patterns")
search("Lambda function deployment code examples")
search("dotnet C# .NET code example standards") - Example for .NET
search("Java code example standards") - Example for Java
```

## Integration Requirements

### For All Code Example Development:
1. **MANDATORY Quality Code Examples KB consultation**: Every code example must begin with searching for "[language] code example standards" in the Quality Code Examples Knowledge Base
2. **Mandatory Tejas KB consultation**: Every code example must begin with Tejas knowledge base research for service understanding
3. **Document KB findings**: Include relevant information from KB queries in code comments
4. **Validate against KB**: Ensure final implementation aligns with KB recommendations
5. **Reference KB sources**: When applicable, reference specific KB insights in documentation

**CRITICAL**: You CANNOT create language-specific code without first consulting the Quality Code Examples Knowledge Base for that language's standards.

### For Service-Specific Examples:
1. **Service capability verification**: Confirm all used features are supported by the service
2. **Parameter validation**: Verify all required parameters are included and optional ones are documented
3. **Error scenario coverage**: Include error handling for service-specific failure modes
4. **Best practice adherence**: Follow service-specific best practices identified in KB research

### For Cross-Service Examples:
1. **Service interaction patterns**: Research how services integrate with each other
2. **Data flow validation**: Ensure data formats are compatible between services
3. **Authentication consistency**: Verify authentication approaches work across all services
4. **Performance considerations**: Research any service-specific performance implications

## Quality Assurance Through Knowledge Base

### Before Code Completion:
- [ ] **MANDATORY**: Searched Quality Code Examples KB for "[language] code example standards"
- [ ] Queried Tejas KB for comprehensive service understanding
- [ ] Searched Quality Code Examples KB for similar implementation patterns
- [ ] Validated implementation against KB recommendations
- [ ] Confirmed error handling covers KB-identified scenarios
- [ ] Verified best practices from KB are implemented
- [ ] **CRITICAL**: Confirmed code structure follows language-specific standards from Quality Code Examples KB

### Documentation Requirements:
- Include KB-sourced service descriptions in code comments
- Reference specific KB insights in README files
- Document any deviations from KB recommendations with justification
- Provide KB-validated parameter explanations

## Troubleshooting with Knowledge Base

### When Implementation Issues Arise:
1. **Query Tejas KB** for service-specific troubleshooting guidance
2. **Search local KB** for similar issues and their resolutions
3. **Cross-reference solutions** to ensure compatibility with repository patterns
4. **Validate fixes** against KB recommendations before implementation

### Common Troubleshooting Queries:
```
query_tejas_kb("Common errors when working with [AWS Service] and how to resolve them")
query_tejas_kb("Why might [specific operation] fail and how to handle it?")
search("error handling examples for [service] in [language]")
```

This knowledge base integration ensures that all code examples are built on a foundation of accurate, comprehensive AWS service knowledge while maintaining consistency with established repository patterns and best practices.