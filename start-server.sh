#!/bin/bash

# Kill any existing server processes
pkill -9 -f bin/server 2>/dev/null
sleep 1

# Start server
cd /home/zzf/projects/goinvent/warehouse
nohup ./bin/server > logs/server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Check if server is running
if ps -p $SERVER_PID > /dev/null 2>&1; then
    echo "✅ Server started successfully"
    echo "   PID: $SERVER_PID"
    echo "   Port: 8080"
    echo "   Access: http://192.168.1.13:8080"
    echo "   Logs: tail -f /home/zzf/projects/goinvent/warehouse/logs/server.log"
    
    # Test if API is responding
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/warehouses 2>/dev/null)
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "401" ]; then
        echo "   API Status: ✅ Responding (HTTP $HTTP_CODE)"
    else
        echo "   API Status: ⚠️  May need attention (HTTP $HTTP_CODE)"
    fi
else
    echo "❌ Failed to start server"
    echo "   Check logs: tail -50 /home/zzf/projects/goinvent/warehouse/logs/server.log"
fi
