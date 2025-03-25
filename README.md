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

## 📋 **Requirements**
☑️ **NodeJS >= 18.17.1**<br>
☑️ **GO (golang) >= 1.22**<br>
☑️ **Electron >= 35.0.3**

## 📦 **Installation & Usage**  
**⭐ Run the WebServer and access it via the Web Browser**

1️⃣ Clone the repository:  
   ```bash
   git clone https://github.com/drclcomputers/GOWEATHER
   ```  
2️⃣ Navigate to the project folder:  
   ```bash
   cd GOWEATHER
   ```
3️⃣ Build the Server executable
   ```bash
   go build
   ```
4️⃣ Run the server executable and navigate to **"localhost:8090"** in your Browser.

**⭐ Run it as a Desktop WebApp**

1️⃣ Clone the repository:  
   ```bash
   git clone https://github.com/drclcomputers/GOWEATHER
   ```  
2️⃣ Navigate to the project folder:  
   ```bash
   cd GOWEATHER
   ```
3️⃣ Install the dependencies with npm:
   ```bash
   npm install
   ```
4️⃣ Build the Server executable:
   ```bash
   go build -o server
   ```
   ‼️ Note: For Windows, write **"server.exe"**. Otherwise leave it as is.
   
5️⃣ Test the App:
   ```bash
   npm start
   ```
   
6️⃣ Build an executable:
   ```bash
   npm run dist
   ```

!!! On Fedora, from what I've tested, this method works. However, in Debian-based or Ununtu-like distros you have to run the commands below to give some permisions.
   ```bash
   sudo chown root:root node_modules/electron/dist/chrome-sandbox
   sudo chmod 4755 node_modules/electron/dist/chrome-sandbox
   ```
!!! I am still unable to figure out why the Appimage won't work. Until then, please use **"npm start"** in the main folder or additionally go to **"dist/linux-unpacked*"**and run **"./goweather"**. 

---

## 📜 **License**  
This project is licensed under the **MIT License**.  

---

Let me know if you want any custom changes! 🚀
