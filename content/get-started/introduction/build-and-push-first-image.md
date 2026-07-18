docker build -t inventory-backend ./backend
docker run -d -p 4000:4000 \
  -v inventory-data:/data \
  --env-file backend/.env \
  inventory-backend


