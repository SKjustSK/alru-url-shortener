# Stage 1: Build the React application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package descriptors and lockfile
COPY package*.json ./
RUN npm ci

# Copy frontend source code
COPY . .

# Accept build arguments to configure Vite build environment
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build application bundle
RUN npm run build

# Stage 2: Serve using Nginx
FROM nginx:stable-alpine

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build artifacts from stage 1
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
