# Kiro RAG Server

This directory contains a Retrieval-Augmented Generation (RAG) server for Kiro that implements the Model Context Protocol (MCP). The RAG server allows Kiro to search and retrieve information from your knowledge base.

## Directory Structure

```
.kiro/
├── knowledge_base/     # Your knowledge base documents go here
├── venv/              # Python virtual environment
├── .vector_store/     # Vector store for embeddings (created automatically)
├── fastmcp.py         # Simple implementation of the MCP protocol
├── rag_server.py      # RAG server implementation
└── settings/
    └── mcp.json       # MCP configuration for Kiro
```

## Setup Instructions

1. **Create a Python virtual environment**:
   ```bash
   python3 -m venv .kiro/venv
   ```

2. **Install dependencies**:
   ```bash
   .kiro/venv/bin/pip install langchain langchain-community sentence-transformers chromadb unstructured pdf2image pytesseract docx2txt
   ```

3. **Add documents to your knowledge base**:
   Place your documents in the `.kiro/knowledge_base/` directory. The server supports:
   - Text files (.txt)
   - PDF documents (.pdf)
   - CSV files (.csv)
   - Word documents (.docx)
   - Markdown files (.md)

4. **Configure MCP in Kiro**:
   The MCP configuration is already set up in `.kiro/settings/mcp.json`.

## Running the RAG Server

The RAG server will be started automatically by Kiro when needed. You can also run it manually for testing:

```bash
.kiro/venv/bin/python .kiro/rag_server.py --kb-path .kiro/knowledge_base --persist-dir .kiro/.vector_store
```

## Available Tools

The RAG server provides the following tools to Kiro:

1. **search**: Search the knowledge base for relevant information
   - Parameters:
     - `query`: The search query
     - `k`: Number of results to return (default: 5)
     - `filter_metadata`: Optional metadata filter

2. **knowledge_base_info**: Get information about the knowledge base
   - Returns:
     - `document_count`: Number of documents in the knowledge base
     - `knowledge_base_path`: Path to the knowledge base
     - `embeddings_model`: Embeddings model being used

3. **reindex_knowledge_base**: Reindex the knowledge base
   - Parameters:
     - `force`: Force reindexing (default: false)

## Using RAG in Kiro

Once everything is set up, you can use the RAG system in Kiro by:

1. Opening Kiro in your IDE
2. The MCP server should connect automatically
3. Ask Kiro questions about your knowledge base, for example:
   - "What programming languages are supported in the AWS SDK examples?"
   - "Tell me about the AWS service examples in this repository"
   - "How can I contribute to the AWS SDK examples?"

## Troubleshooting

If you encounter issues:

1. Check that the virtual environment is properly set up
2. Ensure your knowledge base directory contains documents
3. Check the Kiro logs for any error messages
4. Try running the RAG server manually to see any error output