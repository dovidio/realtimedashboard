### STAGE 1: Build ###
FROM node:12.18.1-stretch AS build
WORKDIR /usr/src/app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build
RUN ls

### STAGE 2: Run ###
FROM nginx:stable
COPY --from=build /usr/src/app/dist/mapdashboard /usr/share/nginx/html

