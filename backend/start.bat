@echo off
echo Starting Bioskop Backend Server...
echo.
echo Make sure PostgreSQL is running and database 'db_bioskop' exists!
echo.
go run main.go
pause

