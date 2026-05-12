# Hey Let's Meet
#### Video Demo: <URL HERE>

#### Description:

Hey Let's Meet is a fullstack web application that suggests a fair walking meeting point between two people. Given two addresses, the app finds the location where both people would walk for the most similar amount of time, minimising the difference in travel time rather than simply splitting the straight-line distance.

The app shows the best meeting point and up to three alternatives on an interactive map. Each result includes walking time and distance for both people, a fairness score (max Eastimated Time of Arrival -ETA- and the difference between the two), and the option to copy the coordinates, open the location in Google Maps, or share it via WhatsApp.


#### App Sample

![Hey Let's Meet](screenshots/Screenshot%202026-05-12%20at%2011.53.35.png)

#### How it works

First, user needs to register, if they are already register they must login. After register/login is succesful the user submits two addresses, the backend geocodes both into coordinates using the OpenRouteService API. It then calculates the geographic midpoint and builds a grid of meeting points around it. For each candidate, it calls the OpenRouteService walking directions API to get real street walking times and distances, not just straight-line estimates. The candidates are then ranked by fairness: first by minimising the maximum ETA (so neither person walks much longer than the other), then by minimising the difference between the two ETAs, and finally by total combined distance. The best result and the top three alternatives are returned to the frontend.

To stay within the OpenRouteService free tier rate limit of 40 requests per minute, the grid is kept small (3×3 = 9 candidates) and all geocoding and routing results are cached in SQLite so repeated queries don't consume extra API calls.


#### Project structure

##### Backend (`backend/`)

The backend is written in Go and exposes a REST API over HTTP.

- **`cmd/api/main.go`:** Entry point. Wires together the database, auth service, meetpoints service, and HTTP server. If an `API_KEY` environment variable is set, it uses real OpenRouteService providers; otherwise it falls back to mock providers for local development without an API key.

- **`internal/httpapi/server.go`:** Sets up the HTTP router and middleware. Protected routes require a valid session cookie. Includes CORS middleware for local development.

- **`internal/auth/`:** Handles user registration, login, logout, and session management. Passwords are hashed with bcrypt. Sessions are stored in SQLite with a configurable Time To Live -TTL-.

- **`internal/meetpoints/logic.go`:** Core algorithm. Geocodes both addresses, builds the candidate grid using Haversine distance and offset math, calls the router for each candidate, sorts results by fairness, and returns the best point plus alternatives. AI was very helpful deciding the logic, and understanding the Haversine formula.

- **`internal/meetpoints/geo.go`:** Geographic math: Haversine distance between two coordinates, midpoint calculation, and coordinate offset by metres in any direction. 

- **`internal/meetpoints/ors.go`:** Real implementation of the Geocoder and Router interfaces using the OpenRouteService API.

- **`internal/meetpoints/mock_geocoder.go` and `mock_router.go`:** Mock implementations used in tests and local development. The mock geocoder uses a hash function to return stable fake coordinates. The mock router uses straight-line walking speed.

- **`internal/meetpoints/cache_repo.go`:** SQLite backed cache for geocoding results and walking routes. Prevents redundant API calls.

- **`migrations/`:** SQL migration files applied in order at startup. Covers users, sessions, and the geocoding and routing cache tables.

##### Frontend (`frontend/`)

The frontend is a React + TypeScript single-page application built with Vite.

- **`src/main.tsx`:** App entry point. Imports global styles and injects the colour palette as CSS variables from `src/lib/colors.ts`.

- **`src/lib/colors.ts`:** Single source of truth for the colour palette (wheat, onyx, oceanic, abyss, nectarine, amber). Values are injected as CSS custom properties at startup so both TypeScript and CSS modules reference the same definitions.

- **`src/lib/api.ts`:** Thin fetch wrapper with JSON handling and structured error extraction.

- **`src/lib/auth.tsx` and `useAuth.ts`:** Auth context and hook. Checks session state on load and exposes login, logout, and register flows to the rest of the app.

- **`src/pages/Home/Home.tsx`:** Main page after login. Contains the address form, submits to the meetpoints API, and renders results. Uses tabs to switch between the best point and alternatives, updating the map marker in real time.

- **`src/components/MapView/MapView.tsx`:** Leaflet map with three markers: origin A, origin B, and the selected meeting point. Autofits the viewport to show all three points. Scroll zoom is disabled unless the mouse is hovering over the map, which prevents accidental zooming while scrolling the page.

- **`src/pages/Login/` and `src/pages/Register/`:** Authentication pages with form validation and error display.


#### Design decisions

**Go over Python:** A personal choice to learn a new language. Despite my more familiarity with Python I preferred to use this opportunity to try Go for the first time.

**SQLite over a cloud database:** SQLite was covered in the CS50 course and I wanted to apply it directly. It is also practical: for a project running on a single server, SQLite is simple to set up, requires no separate database process, and the file lives alongside the backend. A cloud database like Turso would add infrastructure complexity without adding value for a local deployment.

**TypeScript over JavaScript:** TypeScript catches mistakes at compile time rather than at runtime. This prevented several bugs during development that would have been ignored by plain JavaScript.

**React over Vue:** Both are capable frameworks. React was chosen because of recent familiarity and because React Leaflet provided a well-maintained, idiomatic way to integrate the map with React's component model. I spent extra time learning Go. Ergo, I wanted to go faster on the frontend.

**OpenStreetMap over Google Maps:** Google Maps has a beautiful UI but requires billing setup and has usage costs. OpenStreetMap is free, open source, and sufficient for showing three map markers and fitting a viewport. 

**CSS Modules over Bootstrap or Tailwind:** This was a deliberate choice to practice writing CSS rather than relying on a framework. CSS Modules were chosen over plain CSS because they scope class names to the component, preventing naming collisions and making it clear which styles belong to which component.


#### AI assistance

This project was built with the assistance of Claude (Anthropic) as an AI pair programmer. Claude was used to discuss architectural approaches, explore design tradeoffs, review naming and code clarity, and speed up implementation. All design decisions, feature choices, and direction were made by the author. The AI served as an amplifier, not as a replacement for understanding.
