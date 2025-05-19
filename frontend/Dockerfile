# Base image
FROM node:18-alpine

# App directory
WORKDIR /app

# Copy dependencies
COPY package*.json ./

# Install deps
RUN npm install

# Copy source code
COPY . .

# Build for production
RUN npm run build

# Expose Next.js port
EXPOSE 3000

# Run in production mode
ENV NODE_ENV production
CMD ["npm", "start"]
