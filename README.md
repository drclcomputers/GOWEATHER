# 🌤️ **GoWeather** – Your Simple & Accurate Weather Companion  

GoWeather is a sleek and intuitive weather website that provides real-time weather updates, forecasts, and key climate insights for any location worldwide.  

---

## 🚀 **Features**  
✅ **Instant Weather Search** – Get real-time weather data by entering any city name.  
✅ **Detailed Weather Info** – Temperature, humidity, wind speed, visibility, and more.  
✅ **5-Day Forecast** – Stay ahead with future weather predictions.  
✅ **Mobile-Friendly Design** – Fully responsive for a seamless experience on any device.  
✅ **Fast & Lightweight** – Minimalistic UI with quick data fetching.  

---

## 🛠️ **Technologies Used**  
- **Frontend:** HTML, CSS, JavaScript  
- **API:** OpenWeatherMap

---

## 📋 **Requirements (only for building it yourself)**
☑️ **NodeJS >= 18.17.1**<br>
☑️ **GO (golang) >= 1.22**<br>
☑️ **Electron >= 35.0.3**

## 📦 **Installation & Usage**  
**⭐ Run the WebServer and access it via the Web Browser**

1️⃣ Download the release for your operating system.

2️⃣ Extract the content.

3️⃣ In the dearchived folder, run the server executable ("server.exe" for Windows, "server" for other OSes).

‼️ Note: For Linux, you'll have to run 
   ```bash
   sudo chmod +x server
   ```
to make the server an executable.
    
4️⃣ Navigate to **"localhost:8090"** in your Browser.
<br><br><br>
**⭐ Run it as a Desktop WebApp**

1️⃣ Download the release for your operating system.

2️⃣ Extract the content.

‼️ Due to some problems with permissions on Linux, you should run
   ```bash
   sudo chmod +x install.sh
   ```
before, to give permissions to the install script and then run
   ```bash
   ./install.sh
   ```
to install the app.

3️⃣ In the dearchived folder, run the "goweather" executable ("goweather.exe" for Windows, "goweather" for other OSes).
<br><br><br>
**⭐ Build it yourself**

**‼️ Only for professionals**

1️⃣ Use git clone to clone the repo.

2️⃣ In the repo's folder, run
   ```bash
   npm install
   ```
3️⃣ Then compile the server with
   ```bash
   go build -o server
   ```
3️⃣ After this, use 
   ```bash
   npm start
   ```
or
   ```bash
   npm run dist
   ```
to build the app.

‼️ On Linux, you should give the permissions 
   ```bash
   sudo chown root:root ./chrome-sandbox
   sudo chmod 4755 chrome-sandbox
   ```
to chrome-sandbox before, to avoid further problems. Even after doing this, you still won't be able to install the app via the newly created AppImage due to the format's flawed design. I recommend using the first method instead. Build it at your own sanity's risk.

---

## 📜 **License**  
This project is licensed under the **MIT License**.  

---

Let me know if you want any custom changes! 🚀
