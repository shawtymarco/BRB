# Use official Node.js image
FROM node:20

# Set working directory
WORKDIR /app

# Copy package files
COPY . .

# Install dependencies
RUN npm install

# Install tsx globally if needed
RUN npm install -g tsx

# Copy the rest of your bot code
COPY . .

# Expose API port (if used)
EXPOSE 8080

# Start the bot
CMD ["tsx", "src/index.ts"]
