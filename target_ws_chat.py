import websockets
import asyncio
import json
import logging

logger = logging.getLogger("arishem.target_ws_chat")

# Cache mapping each parallel asyncio.Task (session) to its own websocket connection
_ws_connections = {}

async def _get_ws():
    global _ws_connections
    
    # 1. Clean up finished tasks from the cache to prevent memory leaks
    for task in list(_ws_connections.keys()):
        if task.done():
            ws_to_close = _ws_connections.pop(task, None)
            if ws_to_close:
                try:
                    await ws_to_close.close()
                except Exception:
                    pass
                
    # 2. Get/create connection for the current running Task (session context)
    current_task = asyncio.current_task()
    ws = _ws_connections.get(current_task)
    
    if ws is None:
        logger.info(f"Establishing isolated persistent connection for Task {id(current_task)}...")
        ws = await websockets.connect("ws://localhost:5000/ws")
        _ws_connections[current_task] = ws
        
    return ws

async def arishem_chat(user_message: str) -> str:
    """
    an expert sports AI assistant with deep knowledge across all major sports
    including football, basketball, soccer, tennis, cricket, baseball, and more.
    """
    current_task = asyncio.current_task()
    
    # Try up to 2 attempts: first attempt using cached connection,
    # second attempt with a fresh connection if the cached one was closed/dead.
    for attempt in range(2):
        ws = None
        try:
            ws = await _get_ws()
            # Send the payload to the server
            payload = {"message": user_message}
            await ws.send(json.dumps(payload))
            
            # Receive and parse response
            print("Waiting response")
            res = await ws.recv()
            print("Received response", res)
            data = json.loads(res)
            return data.get("output", "")
            
        except (websockets.exceptions.ConnectionClosed, ConnectionError, AttributeError) as e:
            logger.warning(
                f"Connection error on Task {id(current_task)} (attempt {attempt + 1}/2): {e}. "
                "Discarding cached connection and retrying..."
            )
            # Discard the dead connection from cache
            _ws_connections.pop(current_task, None)
            if ws:
                try:
                    await ws.close()
                except Exception:
                    pass
            # If we already retried and failed, return the error
            if attempt == 1:
                return f"WebSocket Error: {e}"
                
        except Exception as e:
            # Handle non-connection errors (e.g. JSON decode error) immediately without retry
            logger.error(f"WebSocket communication error on Task {id(current_task)}: {e}")
            _ws_connections.pop(current_task, None)
            if ws:
                try:
                    await ws.close()
                except Exception:
                    pass
            return f"WebSocket Error: {e}"