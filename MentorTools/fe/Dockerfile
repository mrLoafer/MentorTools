# fe/Dockerfile

# Stage 1: Build the React app
FROM node:16-alpine AS builder

WORKDIR /app

# Копируем package.json и package-lock.json и устанавливаем зависимости
COPY package.json package-lock.json ./
RUN npm install

# Копируем исходный код и собираем приложение
COPY . .
RUN npm run build

# Stage 2: Serve the React app using Nginx
FROM nginx:stable-alpine

# Копируем собранный проект в Nginx
COPY --from=builder /app/build /usr/share/nginx/html

# Указываем порт для Nginx
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]