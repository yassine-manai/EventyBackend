@echo off

::set /p version="Enter Docker version (e.g., 0.0.x): "


echo Building Docker image
docker build -t yassinemanai/eventy_backend:1.0.2 .

if %errorlevel% neq 0 (
    echo Building failed with error code %errorlevel%.
    exit /b 1
)

echo Building completed successfully.

echo Pushing to DOCKER HUB
docker push yassinemanai/eventy_backend:1.0.2

if %errorlevel% neq 0 (
    echo Pushing failed with error code %errorlevel%.
    exit /b 1
)

echo Push completed successfully.
