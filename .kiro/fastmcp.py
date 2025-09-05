"""
Simple implementation of the MCP protocol for the RAG server.
"""

import json
import logging
import inspect
import asyncio
import argparse
from typing import Any, Callable, Dict, List, Optional, Type, get_type_hints, Union
from dataclasses import dataclass, field

logger = logging.getLogger(__name__)

@dataclass
class ToolParameter:
    """Parameter for a tool."""
    description: str
    default: Any = None

@dataclass
class Tool:
    """Decorator for tools."""
    description: str = ""
    
    def __call__(self, func):
        func._is_tool = True
        func._tool_description = self.description
        return func

class FastMCP:
    """Simple implementation of the MCP protocol."""
    
    def __init__(self):
        self.tools = {}
        
    def add_tool(self, func: Callable):
        """Add a tool to the MCP server."""
        if not hasattr(func, '_is_tool'):
            # If not decorated, add the attributes
            func._is_tool = True
            func._tool_description = func.__doc__ or ""
        
        name = func.__name__
        self.tools[name] = func
        logger.info(f"Added tool: {name}")
        
    async def handle_request(self, request_data: Dict[str, Any]) -> Dict[str, Any]:
        """Handle an MCP request."""
        try:
            method = request_data.get("method")
            request_id = request_data.get("id")
            params = request_data.get("params", {})
            
            if method == "initialize":
                result = self._handle_initialize(params)
                return {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "result": result
                }
            elif method == "initialized":
                # This is a notification, no response needed
                return None
            elif method == "tools/list":
                result = self._handle_tools_list()
                return {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "result": result
                }
            elif method == "tools/call":
                result = await self._handle_tool_call(params)
                return {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "result": result
                }
            else:
                return {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "error": {"code": -32601, "message": f"Method not found: {method}"}
                }
        except Exception as e:
            logger.error(f"Error handling request: {e}")
            return {
                "jsonrpc": "2.0",
                "id": request_data.get("id"),
                "error": {"code": -32603, "message": str(e)}
            }
    
    def _handle_initialize(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle an initialize request."""
        return {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "tools": {}
            },
            "serverInfo": {
                "name": "knowledge-base-server",
                "version": "1.0.0"
            }
        }
    
    def _handle_tools_list(self) -> Dict[str, Any]:
        """Handle a tools/list request."""
        tools = []
        
        for name, func in self.tools.items():
            tool_info = {
                "name": name,
                "description": getattr(func, "_tool_description", func.__doc__ or ""),
                "inputSchema": {
                    "type": "object",
                    "properties": self._get_parameters(func),
                    "required": self._get_required_parameters(func)
                }
            }
            tools.append(tool_info)
        
        return {
            "tools": tools
        }
    
    def _get_parameters(self, func: Callable) -> Dict[str, Dict[str, Any]]:
        """Get the parameters for a function."""
        params = {}
        
        sig = inspect.signature(func)
        for param_name, param in sig.parameters.items():
            param_info = {
                "type": "string"  # Default to string, could be enhanced
            }
            
            # Handle ToolParameter objects
            if isinstance(param.default, ToolParameter):
                param_info["description"] = param.default.description
                if param.default.default is not None:
                    param_info["default"] = param.default.default
            elif param.default != inspect.Parameter.empty:
                param_info["default"] = param.default
            
            # Set type based on annotation
            if param.annotation != inspect.Parameter.empty:
                if param.annotation == str:
                    param_info["type"] = "string"
                elif param.annotation == int:
                    param_info["type"] = "integer"
                elif param.annotation == bool:
                    param_info["type"] = "boolean"
                elif hasattr(param.annotation, '__origin__'):
                    # Handle Optional types
                    if param.annotation.__origin__ is Union:
                        args = param.annotation.__args__
                        if len(args) == 2 and type(None) in args:
                            # This is Optional[T]
                            non_none_type = args[0] if args[1] is type(None) else args[1]
                            if non_none_type == str:
                                param_info["type"] = "string"
                            elif non_none_type == int:
                                param_info["type"] = "integer"
                            elif non_none_type == bool:
                                param_info["type"] = "boolean"
                            elif non_none_type == dict:
                                param_info["type"] = "object"
            
            params[param_name] = param_info
        
        return params
    
    def _get_required_parameters(self, func: Callable) -> List[str]:
        """Get the required parameters for a function."""
        required = []
        
        sig = inspect.signature(func)
        for param_name, param in sig.parameters.items():
            # Parameter is required if it has no default value, or if it's a ToolParameter without a default
            if param.default == inspect.Parameter.empty:
                required.append(param_name)
            elif isinstance(param.default, ToolParameter) and param.default.default is None:
                required.append(param_name)
        
        return required
    
    async def _handle_tool_call(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle a tools/call request."""
        tool_name = params.get("name")
        arguments = params.get("arguments", {})
        
        if tool_name not in self.tools:
            return {"error": f"Tool not found: {tool_name}"}
        
        try:
            func = self.tools[tool_name]
            
            # Call the function with the provided arguments
            if inspect.iscoroutinefunction(func):
                result = await func(**arguments)
            else:
                result = func(**arguments)
            
            return {"content": [{"type": "text", "text": str(result)}]}
        except Exception as e:
            logger.error(f"Error executing tool {tool_name}: {e}")
            return {"error": str(e)}
        sig = inspect.signature(func)
        type_hints = get_type_hints(func)
        
        for name, param in sig.parameters.items():
            if name == "self":
                continue
                
            param_info = {
                "type": self._get_type_name(type_hints.get(name, Any)),
                "description": "",
                "required": param.default is inspect.Parameter.empty
            }
            
            # Check if the parameter has a default value
            if param.default is not inspect.Parameter.empty:
                if not isinstance(param.default, ToolParameter):
                    param_info["default"] = param.default
            
            # Check if the parameter is annotated with ToolParameter
            if param.default is not inspect.Parameter.empty and isinstance(param.default, ToolParameter):
                param_info["description"] = param.default.description
                if param.default.default is not None:
                    param_info["default"] = param.default.default
            
            params[name] = param_info
        
        return params
    
    def _get_type_name(self, type_hint: Type) -> str:
        """Get the name of a type."""
        if hasattr(type_hint, "__origin__"):
            if type_hint.__origin__ is Union:
                # Handle Optional[T] which is Union[T, None]
                if len(type_hint.__args__) == 2 and type_hint.__args__[1] is type(None):
                    return self._get_type_name(type_hint.__args__[0])
                # Handle other Union types
                return "union"
            elif type_hint.__origin__ is list:
                return "array"
            elif type_hint.__origin__ is dict:
                return "object"
        
        if type_hint is str:
            return "string"
        elif type_hint is int:
            return "integer"
        elif type_hint is float:
            return "number"
        elif type_hint is bool:
            return "boolean"
        elif type_hint is Any:
            return "any"
        
        return "object"
    

    
    async def _server_loop(self, host: str, port: int):
        """Run the MCP server."""
        server = await asyncio.start_server(
            self._handle_client, host, port
        )
        
        addr = server.sockets[0].getsockname()
        logger.info(f'Serving on {addr}')
        
        async with server:
            await server.serve_forever()
    
    async def _handle_client(self, reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
        """Handle a client connection."""
        addr = writer.get_extra_info('peername')
        logger.info(f"Client connected: {addr}")
        
        while True:
            # Read the length of the message
            length_bytes = await reader.read(4)
            if not length_bytes:
                break
                
            length = int.from_bytes(length_bytes, byteorder='little')
            
            # Read the message
            data = await reader.read(length)
            if not data:
                break
                
            # Parse the message
            try:
                request = json.loads(data.decode())
                logger.debug(f"Received request: {request}")
                
                # Handle the request
                response = await self.handle_request(request)
                
                # Send the response
                response_bytes = json.dumps(response).encode()
                length_bytes = len(response_bytes).to_bytes(4, byteorder='little')
                
                writer.write(length_bytes + response_bytes)
                await writer.drain()
            except Exception as e:
                logger.error(f"Error handling client: {e}")
                break
        
        logger.info(f"Client disconnected: {addr}")
        writer.close()
    
    def run(self, host: str = "127.0.0.1", port: int = 8000):
        """Run the MCP server."""
        asyncio.run(self._server_loop(host, port))
    
    def run_stdio(self):
        """Run the MCP server using stdio."""
        asyncio.run(self._stdio_loop())
    
    async def _stdio_loop(self):
        """Run the MCP server using stdio."""
        import sys
        
        logger.info("Starting MCP server in stdio mode")
        
        while True:
            try:
                # Read from stdin
                line = await asyncio.get_event_loop().run_in_executor(None, sys.stdin.readline)
                if not line:
                    break
                
                line = line.strip()
                if not line:
                    continue
                
                # Parse the JSON-RPC request
                request = json.loads(line)
                logger.debug(f"Received request: {request}")
                
                # Handle the request
                response = await self.handle_request(request)
                
                # Send the response to stdout (only if there's a response)
                if response is not None:
                    response_json = json.dumps(response)
                    print(response_json, flush=True)
                
            except json.JSONDecodeError as e:
                logger.error(f"JSON decode error: {e}")
                error_response = {
                    "jsonrpc": "2.0",
                    "id": None,
                    "error": {"code": -32700, "message": "Parse error"}
                }
                print(json.dumps(error_response), flush=True)
            except Exception as e:
                logger.error(f"Error in stdio loop: {e}")
                break