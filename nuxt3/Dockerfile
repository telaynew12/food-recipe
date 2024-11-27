FROM node:18.17.1

# Set working directory
WORKDIR /app

# Copy package manager files first to leverage Docker cache
COPY package.json package-lock.json ./

# Force clean cache and install dependencies
RUN npm cache clean --force && rm -rf node_modules package-lock.json
RUN yarn install --omit=optional --legacy-peer-deps

# Copy the rest of the application code
COPY . .

# Build Nuxt 3 app
RUN npm run build

# Expose Nuxt's default port
EXPOSE 3000

# Start Nuxt in production
CMD ["npm", "run", "start"]
