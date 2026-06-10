// This object holds our camera data
let cameras = [];
let markers = {};

// Initialize Map
const map = L.map("map").setView([51.505, -0.09], 13);
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  attribution: "© OpenStreetMap",
}).addTo(map);

// Function to update the video player
function playVideo(url, name) {
  document.getElementById("cam-name").innerText = name;
  const container = document.getElementById("video-container");
  if (url) {
    // Convert embed URL back to watch URL
    const watchUrl = url.replace(
      "https://www.youtube.com/embed/",
      "https://www.youtube.com/watch?v=",
    );
    container.innerHTML = `
      <div style="text-align: center; padding: 30px;">
        <button onclick="window.go.main.App.OpenURL('${watchUrl}')" style="padding: 14px 28px; background: #ff0000; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; font-size: 1.05em; transition: background 0.2s;" onmouseover="this.style.background='#cc0000'" onmouseout="this.style.background='#ff0000'">
          ▶ Watch Live on YouTube
        </button>
        <p style="margin: 15px 0 0 0; font-size: 0.85em; color: #999;">Opens in your browser</p>
      </div>
    `;
  } else {
    container.innerHTML = `<p style="color: #999; font-size: 0.9em;">No YouTube URL provided.</p>`;
  }
}

// Function to add a marker to the map
function addMarkerToMap(cam) {
  const marker = L.marker([cam.lat, cam.lng]).addTo(map);
  marker.on("click", () => {
    playVideo(cam.youtubeUrl, cam.name);
  });

  let popupText = `<strong>${cam.name}</strong>`;
  if (cam.description) {
    popupText += `<br/>${cam.description}`;
  }
  marker.bindPopup(popupText);

  markers[cam.id] = marker;
}

// Function to refresh map markers and update camera list
function refreshMarkers() {
  Object.values(markers).forEach((marker) => map.removeLayer(marker));
  markers = {};
  cameras.forEach(addMarkerToMap);
  updateCameraList();
}

// Function to update the camera list in the sidebar
function updateCameraList() {
  const listEl = document.getElementById("cameras-list");
  listEl.innerHTML = "";

  cameras.forEach((cam) => {
    const li = document.createElement("li");
    li.textContent = cam.name;
    li.onclick = () => {
      selectCamera(cam);
    };
    li.dataset.camId = cam.id;
    listEl.appendChild(li);
  });
}

// Function to select a camera from the list
function selectCamera(cam) {
  console.log(
    "Selecting camera from list:",
    cam.name,
    "youtubeUrl:",
    cam.youtubeUrl,
  );
  console.log("Cam object keys:", Object.keys(cam));
  console.log("Cam object values:", Object.values(cam));

  // Update active state in list
  document.querySelectorAll("#cameras-list li").forEach((li) => {
    li.classList.remove("active");
  });
  document.querySelector(`[data-cam-id="${cam.id}"]`)?.classList.add("active");

  // Show video and center map
  playVideo(cam.youtubeUrl, cam.name);
  map.setView([cam.lat, cam.lng], 13);
}

// Import file handling using Wails file dialog
document.getElementById("import-btn").addEventListener("click", async () => {
  const statusDiv = document.getElementById("import-status");

  try {
    // Call the Go method to open the native file dialog
    const filePath = await window.go.main.App.SelectKMLFile();

    if (!filePath) {
      // User cancelled the dialog
      return;
    }

    statusDiv.innerHTML = "Importing...";
    statusDiv.style.color = "#666";

    // Call the ImportKML method with the selected file path
    const imported = await window.go.main.App.ImportKML(filePath);

    // Refresh the camera list
    cameras = await window.go.main.App.GetCameras();
    refreshMarkers();

    if (imported && imported.length > 0) {
      statusDiv.innerHTML = `✓ Imported ${imported.length} camera(s)`;
      statusDiv.style.color = "#0a0";
    } else {
      statusDiv.innerHTML = `✗ No cameras found in file`;
      statusDiv.style.color = "#f00";
    }

    setTimeout(() => {
      statusDiv.innerHTML = "";
    }, 3000);
  } catch (err) {
    const statusDiv = document.getElementById("import-status");
    const errMsg = err?.message || String(err) || "Unknown error";
    statusDiv.innerHTML = `✗ Import failed: ${errMsg}`;
    statusDiv.style.color = "#f00";
    console.error("Import error:", err);
  }
});

// This function runs when the app starts
window.addEventListener("DOMContentLoaded", async () => {
  try {
    // CALLING GO: Wails maps Go methods to window.go.[package].[Struct].[Method]
    // Since our package is 'main' and struct is 'App'
    cameras = await window.go.main.App.GetCameras();

    // Add markers to the map for each camera
    cameras.forEach(addMarkerToMap);

    // Update the camera list
    updateCameraList();
  } catch (err) {
    console.error("Failed to fetch cameras from Go:", err);
  }
});
