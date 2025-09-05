#!/usr/bin/env python3
"""
MCP Server for Tejas Knowledge Base API using FastMCP
"""

import os
import sys
import json
import requests
import logging

# Import the local fastmcp implementation
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from fastmcp import FastMCP, Tool, ToolParameter

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Create MCP server
app = FastMCP()

@Tool(description="Query the Tejas knowledge base with a question")
def query_tejas_kb(query: str = ToolParameter(description="The question or query to search for in the knowledge base")) -> str:
    """Query the Tejas knowledge base API."""
    api_endpoint = "https://kcr1d1gr9a.execute-api.us-west-2.amazonaws.com/prod/query"
    
    try:
        response = requests.post(
            api_endpoint,
            headers={"Content-Type": "application/json"},
            json={"query": query},
            timeout=30
        )
        response.raise_for_status()
        
        result = response.json()
        
        # Format the response nicely
        if isinstance(result, dict):
            formatted_result = json.dumps(result, indent=2)
        else:
            formatted_result = str(result)
        
        return f"Query: {query}\n\nResponse:\n{formatted_result}"
        
    except requests.exceptions.Timeout:
        return "Error: Request timed out"
    except requests.exceptions.HTTPError as e:
        return f"Error: HTTP {e.response.status_code} - {e.response.text}"
    except Exception as e:
        return f"Error: {str(e)}"

# Register the tool
app.add_tool(query_tejas_kb)

def main():
    """Main entry point for the Tejas KB server."""
    logger.info("Starting Tejas KB server in stdio mode")
    app.run_stdio()

if __name__ == "__main__":
    main()