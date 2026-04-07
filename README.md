# 🅿️ ParkQ: Real-Time Location-Sharing System

Welcome to the **ParkQ Location-Sharing System**! This interactive frontend is built with React, Vite, and Leaflet. It visualizes live movement across a group using real-time Geolocation tracking and WebSocket synchronization.

## ✨ Features

- **🗺️ Interactive Maps**: Choose between multiple map styles (Monochrome, Terra, Standard, Satellite) powered by React-Leaflet.
- **📍 Real-Time Tracking**: Polling device GPS via `navigator.geolocation.watchPosition` to ensure live location updates.
- **👥 Group Synchronization**: Uses WebSockets to manage concurrent connections, allowing peers to share and broadcast their coordinates instantly (Currently mocked for frontend demonstration).
- **🧮 Haversine Proximity**: Accurately calculates the real-world distance between users mathematically using the Haversine formula.
- **📱 Mobile Optimized**: Clean, native-app-like UI built with Framer Motion, fully optimized for mobile viewports (`100dvh`).

## 🚀 Getting Started

### Prerequisites

Make sure you have [Node.js](https://nodejs.org/) installed.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Shafwansafi06/location-sharing-system.git
   cd location-sharing-system
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

## 🛠️ Tech Stack

- **Frontend Framework**: [React](https://reactjs.org/) + [Vite](https://vitejs.dev/)
- **Map Rendering**: [Leaflet](https://leafletjs.com/) & [React-Leaflet](https://react-leaflet.js.org/)
- **Animations**: [Framer Motion](https://www.framer.com/motion/)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Backend (Target Pipeline)**: Built for [Golang WebSockets](https://go.dev/) (Part 1 backend specification mapped in `src/App.jsx`).

## 📝 Current Status (Mocked Environment)

To facilitate immediate UI/UX testing, the current build includes a **mocked API layer**:
- **GPS Simulation**: Device coordinates are simulated to drift naturally from a starting point.
- **WebSocket Switchboard**: Group connections are simulated locally to dispatch mocked peers (Alex, Mila) joining the map and jittering around to demonstrate smooth marker animations and real-time distance calculations.

Once the Golang backend is spun up, simply uncomment the production WebSocket configuration in `src/App.jsx` and delete the mockup code block!

---
*Built for the ParkQ Smart Parking Initiative.*
