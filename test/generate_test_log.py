#!/usr/bin/env python3
"""Generate a test log file with ANSI color codes"""

print("\033[1;32m[INFO]\033[0m Starting application...")
print("\033[1;34m[DEBUG]\033[0m Loading configuration from config.yaml")
print("\033[1;33m[WARN]\033[0m Configuration file not found, using defaults")
print("\033[1;32m[INFO]\033[0m Connecting to database at \033[4mlocalhost:5432\033[0m")
print("\033[1;31m[ERROR]\033[0m Failed to connect: Connection refused")
print("\033[1;33m[WARN]\033[0m Retrying in 5 seconds...")
print("\033[1;32m[INFO]\033[0m Connection successful")
print("\033[1;36m[TRACE]\033[0m Query executed: SELECT * FROM users WHERE id=1")
print("\033[1;32m[INFO]\033[0m User authenticated: \033[1;95mjohn.doe@example.com\033[0m")
print("\033[1;34m[DEBUG]\033[0m Session ID: \033[2;37mabc123def456\033[0m")
print("\033[1;32m[INFO]\033[0m Starting HTTP server on \033[1;4;33mhttp://localhost:8080\033[0m")
print("")
print("\033[1;46;30m=== Request Log ===\033[0m")
print("\033[90m2024-01-27 10:15:32\033[0m \033[1;32mGET\033[0m /api/users \033[1;32m200\033[0m 45ms")
print("\033[90m2024-01-27 10:15:33\033[0m \033[1;33mPOST\033[0m /api/login \033[1;33m401\033[0m 12ms")
print("\033[90m2024-01-27 10:15:34\033[0m \033[1;32mGET\033[0m /api/products \033[1;32m200\033[0m 78ms")
print("\033[90m2024-01-27 10:15:35\033[0m \033[1;31mDELETE\033[0m /api/users/5 \033[1;31m403\033[0m 5ms")
print("")
print("\033[1;41;97m!!! CRITICAL ERROR !!!\033[0m")
print("\033[1;31mStack trace:\033[0m")
print("\033[2;37m  at handleRequest (server.js:42:15)\033[0m")
print("\033[2;37m  at processRequest (router.js:128:9)\033[0m")
print("\033[2;37m  at Server.handle (http.js:456:12)\033[0m")
print("")
print("\033[7;32m Application recovered and running normally \033[0m")
print("\033[1;32m[INFO]\033[0m Health check: \033[1;42;97m OK \033[0m")
print("\033[1;34m[DEBUG]\033[0m Memory usage: \033[1;96m128MB\033[0m / \033[1;93m512MB\033[0m")
print("\033[1;32m[INFO]\033[0m Active connections: \033[1;35m42\033[0m")
print("")
print("\033[3;90mApplication log ends here...\033[0m")

# Generate many more lines for scrolling test
for i in range(100):
    color = 32 + (i % 6)  # Cycle through colors
    print(f"\033[1;{color}m[LOG {i+1:03d}]\033[0m This is log line number {i+1} with some \033[1;4munderlined\033[0m and \033[1mbold\033[0m text")
