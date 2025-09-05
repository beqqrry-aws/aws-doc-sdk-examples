#!/usr/bin/env python3
"""
RAG Server for Kiro using MCP protocol
"""

import os
import argparse
import logging
import glob
from typing import List, Optional

# Import our local fastmcp implementation
import sys
import os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from fastmcp import FastMCP, Tool, ToolParameter
from langchain_community.vectorstores import Chroma
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_community.document_loaders import (
    DirectoryLoader,
    TextLoader,
    PyPDFLoader,
    CSVLoader,
    Docx2txtLoader,
    UnstructuredMarkdownLoader
)
from langchain.text_splitter import RecursiveCharacterTextSplitter

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Create MCP server
app = FastMCP()

# Global variables
vector_store = None
knowledge_base_path = None
persist_directory = None
embeddings_model = "all-MiniLM-L6-v2"  # Default model

def setup_vector_store(kb_path: str, persist_dir: str = "./.vector_store"):
    """Set up and return the vector store with the knowledge base."""
    global vector_store, knowledge_base_path, persist_directory
    
    # Resolve relative paths
    kb_path = os.path.abspath(os.path.expanduser(kb_path))
    persist_dir = os.path.abspath(os.path.expanduser(persist_dir))
    
    knowledge_base_path = kb_path
    persist_directory = persist_dir
    logger.info(f"Setting up vector store with knowledge base at: {kb_path}")
    logger.info(f"Vector store persist directory: {persist_dir}")
    logger.info(f"Current working directory: {os.getcwd()}")
    
    # Check if the knowledge base path exists
    if not os.path.exists(kb_path):
        logger.error(f"Knowledge base path does not exist: {kb_path}")
        raise FileNotFoundError(f"Knowledge base path does not exist: {kb_path}")
    
    # Initialize embeddings
    embeddings = HuggingFaceEmbeddings(model_name=embeddings_model)
    
    # Check if we have an existing vector store
    if os.path.exists(persist_dir) and os.path.isdir(persist_dir) and os.listdir(persist_dir):
        logger.info(f"Loading existing vector store from: {persist_dir}")
        try:
            vector_store = Chroma(persist_directory=persist_dir, embedding_function=embeddings)
            # Check if the vector store has any documents
            collection = vector_store._collection
            count = collection.count()
            if count > 0:
                logger.info(f"Successfully loaded existing vector store with {count} documents")
                return vector_store
            else:
                logger.info("Existing vector store is empty, will recreate it")
        except Exception as e:
            logger.warning(f"Failed to load existing vector store: {e}")
            logger.info("Will create a new vector store instead")
    
    # Create a new vector store
    logger.info("Creating new vector store...")
    
    # Load documents from the knowledge base
    documents = load_documents(kb_path)
    logger.info(f"Loaded {len(documents)} documents from knowledge base")
    
    # Split documents into chunks
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=100)
    splits = text_splitter.split_documents(documents)
    logger.info(f"Split into {len(splits)} chunks")
    
    # Create and persist the vector store
    try:
        # Ensure the persist directory exists with proper permissions
        os.makedirs(persist_dir, mode=0o755, exist_ok=True)
        logger.info(f"Created/verified persist directory: {persist_dir}")
        
        vector_store = Chroma.from_documents(
            documents=splits,
            embedding=embeddings,
            persist_directory=persist_dir
        )
        vector_store.persist()
        logger.info(f"Vector store created and persisted to: {persist_dir}")
    except Exception as e:
        logger.error(f"Failed to create vector store: {e}")
        raise
    
    return vector_store

def load_documents(directory_path):
    """Load documents from the specified directory."""
    documents = []
    
    # Load text files
    try:
        txt_loader = DirectoryLoader(directory_path, glob="**/*.txt", loader_cls=TextLoader)
        txt_docs = txt_loader.load()
        logger.info(f"Loaded {len(txt_docs)} .txt documents")
        documents.extend(txt_docs)
    except Exception as e:
        logger.warning(f"Error loading .txt documents: {e}")
    
    # Load markdown files with TextLoader (more reliable than UnstructuredMarkdownLoader)
    try:
        md_loader = DirectoryLoader(directory_path, glob="**/*.md", loader_cls=TextLoader)
        md_docs = md_loader.load()
        logger.info(f"Loaded {len(md_docs)} .md documents")
        documents.extend(md_docs)
    except Exception as e:
        logger.error(f"Error loading .md documents: {e}")
        # Try loading markdown files individually if batch loading fails
        import glob
        md_files = glob.glob(os.path.join(directory_path, "**/*.md"), recursive=True)
        logger.info(f"Found {len(md_files)} markdown files, trying individual loading")
        for md_file in md_files:
            try:
                loader = TextLoader(md_file)
                doc = loader.load()
                documents.extend(doc)
                logger.debug(f"Loaded {md_file}")
            except Exception as file_error:
                logger.warning(f"Failed to load {md_file}: {file_error}")
    
    # Load code files (these are crucial for practical examples)
    code_extensions = [".py", ".rs", ".js", ".java", ".go", ".cpp", ".c", ".h", ".hpp", ".kt", ".php", ".rb", ".swift"]
    for ext in code_extensions:
        try:
            code_loader = DirectoryLoader(directory_path, glob=f"**/*{ext}", loader_cls=TextLoader)
            code_docs = code_loader.load()
            logger.info(f"Loaded {len(code_docs)} {ext} code files")
            documents.extend(code_docs)
        except Exception as e:
            logger.warning(f"Error loading {ext} code files: {e}")
    
    # Load configuration files
    config_extensions = [".toml", ".yaml", ".yml", ".json", ".xml"]
    for ext in config_extensions:
        try:
            config_loader = DirectoryLoader(directory_path, glob=f"**/*{ext}", loader_cls=TextLoader)
            config_docs = config_loader.load()
            logger.info(f"Loaded {len(config_docs)} {ext} config files")
            documents.extend(config_docs)
        except Exception as e:
            logger.warning(f"Error loading {ext} config files: {e}")
    
    # Load other file types
    other_loaders = {
        ".pdf": (DirectoryLoader, {"glob": "**/*.pdf", "loader_cls": PyPDFLoader}),
        ".csv": (DirectoryLoader, {"glob": "**/*.csv", "loader_cls": CSVLoader}),
        ".docx": (DirectoryLoader, {"glob": "**/*.docx", "loader_cls": Docx2txtLoader}),
    }
    
    for ext, (loader_cls, loader_args) in other_loaders.items():
        try:
            loader = loader_cls(directory_path, **loader_args)
            docs = loader.load()
            logger.info(f"Loaded {len(docs)} {ext} documents")
            documents.extend(docs)
        except Exception as e:
            logger.warning(f"Error loading {ext} documents: {e}")
    
    logger.info(f"Total documents loaded: {len(documents)}")
    return documents

@Tool(description="Search the knowledge base for relevant information")
def search(
    query: str = ToolParameter(description="The search query"),
    k: int = ToolParameter(description="Number of results to return", default=5)
) -> List[dict]:
    """Search the knowledge base for relevant information."""
    global vector_store
    
    if vector_store is None:
        return {"error": "Vector store not initialized. Please set up the knowledge base first."}
    
    # Test if vector store is accessible
    try:
        collection = vector_store._collection
        count = collection.count()
        logger.info(f"Vector store has {count} documents")
    except Exception as e:
        logger.error(f"Vector store test failed: {e}")
        return {"error": f"Vector store is not accessible: {e}"}
    
    try:
        # Perform similarity search without filter for now
        docs = vector_store.similarity_search(query, k=k)
        
        # Format results
        results = []
        for i, doc in enumerate(docs):
            results.append({
                "content": doc.page_content,
                "metadata": doc.metadata,
                "relevance_score": i+1  # Simple ranking, lower is better
            })
        
        return results
    except Exception as e:
        logger.error(f"Error during search: {e}")
        return {"error": str(e)}

@Tool(description="Get information about the knowledge base")
def knowledge_base_info() -> dict:
    """Get information about the knowledge base."""
    global vector_store, knowledge_base_path
    
    if vector_store is None:
        return {"error": "Vector store not initialized"}
    
    try:
        collection = vector_store._collection
        count = collection.count()
        return {
            "document_count": count,
            "knowledge_base_path": knowledge_base_path,
            "embeddings_model": embeddings_model
        }
    except Exception as e:
        logger.error(f"Error getting knowledge base info: {e}")
        return {"error": str(e)}

@Tool(description="Reindex the knowledge base")
def reindex_knowledge_base(force: bool = ToolParameter(description="Force reindexing", default=False)) -> dict:
    """Reindex the knowledge base."""
    global vector_store, knowledge_base_path
    
    if knowledge_base_path is None:
        return {"error": "Knowledge base path not set"}
    
    try:
        # If force is True, delete the existing vector store
        if force and vector_store is not None:
            vector_store.delete_collection()
            vector_store = None
        
        # Set up the vector store again
        setup_vector_store(knowledge_base_path, persist_directory)
        
        return {"status": "success", "message": "Knowledge base reindexed successfully"}
    except Exception as e:
        logger.error(f"Error reindexing knowledge base: {e}")
        return {"error": str(e)}

# Register tools with the MCP server
app.add_tool(search)
app.add_tool(knowledge_base_info)
app.add_tool(reindex_knowledge_base)

def main():
    """Main entry point for the RAG server."""
    parser = argparse.ArgumentParser(description="RAG Server for Kiro using MCP protocol")
    parser.add_argument("--kb-path", required=True, help="Path to the knowledge base directory")
    parser.add_argument("--persist-dir", default="./.vector_store", help="Directory to persist the vector store")
    parser.add_argument("--embeddings-model", default="all-MiniLM-L6-v2", help="Embeddings model to use")
    parser.add_argument("--host", help="Host to bind the server to (if not specified, uses stdio mode)")
    parser.add_argument("--port", type=int, default=8000, help="Port to bind the server to")
    parser.add_argument("--stdio", action="store_true", help="Use stdio mode for MCP communication")
    
    args = parser.parse_args()
    
    global embeddings_model
    embeddings_model = args.embeddings_model
    
    # Set up the vector store
    try:
        setup_vector_store(args.kb_path, args.persist_dir)
        logger.info("Vector store setup completed successfully")
    except Exception as e:
        logger.error(f"Failed to set up vector store: {e}")
        # Continue anyway to allow the server to start
    
    # Run the MCP server
    if args.stdio or args.host is None:
        logger.info("Starting RAG server in stdio mode")
        app.run_stdio()
    else:
        logger.info(f"Starting RAG server on {args.host}:{args.port}")
        app.run(host=args.host, port=args.port)

if __name__ == "__main__":
    main()