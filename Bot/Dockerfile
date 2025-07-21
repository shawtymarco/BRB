# Use official Node.js image
FROM node:20

# Set working directory
WORKDIR /app

# Copy dependency files
COPY package*.json ./

# Install dependencies
RUN npm install

# Install tsx globally
RUN npm install -g tsx

# Copy source code
COPY . .

# Set default command
CMD ["tsx", "src/index.ts"]