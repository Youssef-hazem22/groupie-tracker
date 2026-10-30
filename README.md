# 🎵 Groupie Trackers

A modern, fast, and responsive web application built with **Go** that interacts with the **Groupie Trackers API**. It provides an engaging interface to explore music artists, band members, release histories, and upcoming or past concert locations across the globe.

🔗 **Live Website:** [groupie-tracker-woad.vercel.app](https://groupie-tracker-woad.vercel.app/)

---

## ✨ Features & Highlights

* 🎨 **Modern & Interactive UI:** A clean, grid-based card layout that seamlessly adapts across all screen sizes (mobile, tablet, and desktop).
* 🔍 **Multi-Field Real-Time Search:** Search dynamically through various criteria including:
  * Artist / Band Names
  * Individual Band Members
  * First Album Release Dates
  * Band Creation Years
  * Concert Cities & Locations
* 🎤 **Comprehensive Artist Profiles:** Dedicated page for each artist featuring high-resolution visuals, active career details, complete member listings, and structured tour schedules.
* 📍 **Concert & Location Tracker:** Organized display of performance venues alongside show counts per location for easy tracking.
* ⚡ **High Performance & Zero Dependencies:** Built using purely Go's standard library for ultra-fast response times and minimal server footprint.
* 🛡️ **Built-in Security:** Automatic HTML output escaping powered by Go's `html/template` package to prevent XSS and template injection attacks.
* 🌐 **Server-Side Data Caching:** Fetches API data at startup and serves in-memory responses for near-instantaneous page loads.

---

## 🛠️ Tech Stack

* **Backend:** Go (Standard Library `net/http`, `html/template`, `encoding/json`)
* **Frontend:** HTML5, CSS3 (Flexbox, CSS Grid, Custom Variables)
* **Hosting & Deployment:** Vercel

---

## 📁 Project Structure

```
.
├── main.go         # Server setup & enhanced HTTP route handling
├── handlers.go     # API integration, search algorithms & request processing
├── objects.go      # Data models and structures
├── go.mod          # Go module file
└── templates/      # User interface templates and styling
    ├── index.html   # Main dashboard (Artist cards & Search bar)
    ├── artist.html  # Artist detail & schedule view
    └── styling.css  # Application stylesheet
```

---

## 🛣️ Navigation & Routes

| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Main page featuring search & artist gallery |
| `GET` | `/artist?id={N}` | Detailed profile page for specific artist |
| `GET` | `/styling.css` | Serves main application styles |

---

## 👤 Author

**Youssef Hazem**

⭐ If you like this project, feel free to give it a star!
