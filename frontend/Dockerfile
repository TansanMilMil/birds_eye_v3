FROM node:24
WORKDIR /project
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
CMD ["npm", "run", "build:watch"]
