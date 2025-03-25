const { app, BrowserWindow } = require('electron');
const path = require('path');
const { spawn } = require('child_process');

let win;
let goServer; // Store the Go server process

function createWindow() {
  try {
    console.log('Creating browser window...');
    win = new BrowserWindow({
      width: 800,
      height: 600,
      webPreferences: {
        nodeIntegration: true,
        //preload: path.join(__dirname, 'preload.js')
      }
    });

    console.log('Loading URL...');
    win.loadURL('http://127.0.0.1:8090/');

    win.on('closed', () => {
      console.log('Window closed.');
      win = null;
    });

    console.log('Window created successfully.');
  } catch (err) {
    console.error('Error creating window:', err);
  }
}

function startGoServer() {
  let goServerPath;
  const appPath = app.getAppPath(); // Get the path to the app's root directory

  if (app.isPackaged) { // If packaged
    if (process.platform === 'win32') { 
      goServerPath = path.join(appPath, '../../server.exe'); // Windows executable
    } else if (process.platform === 'linux' || process.platform === 'darwin') {
      goServerPath = path.join(appPath, '../../server'); // Linux/macOS executable
    } else {
      console.error('Unsupported platform:', process.platform);
      return;
    }
  } else { 
    if (process.platform === 'win32') {
      goServerPath = path.join(appPath, 'server.exe'); // Windows executable
    } else if (process.platform === 'linux' || process.platform === 'darwin') {
      goServerPath = path.join(appPath, 'server'); // Linux/macOS executable
    } else {
      console.error('Unsupported platform:', process.platform);
      return;
    }
  }

  console.log('Starting Go server from:', goServerPath);

  goServer = spawn(goServerPath); // Store the process

  goServer.stdout.on('data', (data) => {
    console.log(`Go server: ${data}`);
  });

  goServer.stderr.on('data', (data) => {
    console.error(`Go server error: ${data}`);
  });

  goServer.on('close', (code) => {
    console.log(`Go server process exited with code ${code}`);
  });

  goServer.on('error', (err) => {
    console.error('Failed to start Go server:', err);
  });
}

const fs = require('fs');

if (process.platform === 'linux') {
  const electronPath = path.join(__dirname, 'node_modules', 'electron', 'dist', 'chrome-sandbox');
  try {
    fs.chmodSync(electronPath, 0o4755); // Set correct permissions (rwsr-xr-x)
    console.log('✅ Fixed chrome-sandbox permissions.');
  } catch (err) {
    console.error('❌ Failed to set chrome-sandbox permissions:', err);
  }
  if (process.env.APPIMAGE) {
    const electronPath = path.join(path.dirname(process.env.APPIMAGE), 'chrome-sandbox');
    try {
      fs.chmodSync(electronPath, 0o4755);
      console.log('✅ Fixed chrome-sandbox permissions (AppImage).');
    } catch (err) {
      console.error('❌ Failed to fix chrome-sandbox permissions:', err);
    }
  }
  
}

app.whenReady().then(() => {
  console.log('App is ready.');
  startGoServer();
  createWindow();
});

app.on('window-all-closed', () => {
  console.log('All windows closed.');
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  console.log('App activated.');
  if (win === null) {
    createWindow();
  }
});

app.on("before-quit", () => {
  if (goServer) {
    console.log("Stopping Go server...");
    goServer.kill("SIGTERM"); // Send termination signal
  }
});
