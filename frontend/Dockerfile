FROM node:20-alpine

WORKDIR /app

# Copy only dependency files first (cache-friendly)
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest
COPY . .

RUN npm run build

EXPOSE 3000

CMD ["npm", "run", "preview"]