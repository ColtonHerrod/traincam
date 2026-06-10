// This object holds our camera data
let cameras = [];
let markers = {};

// Initialize Map
const map = L.map("map").setView([51.505, -0.09], 13);
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  attribution: "© OpenStreetMap",
}).addTo(map);

// Function to update the video player
function playVideo(url, name, subscriptionRequired = false) {
  document.getElementById("cam-name").innerText = name;
  const container = document.getElementById("video-container");
  
  if (url) {
    // Convert embed URL back to watch URL
    const watchUrl = url.replace(
      "https://www.youtube.com/embed/",
      "https://www.youtube.com/watch?v=",
    );
    
    if (subscriptionRequired) {
      container.innerHTML = `
        <div style="text-align: center; padding: 30px;">
          <p style="margin-bottom: 15px; color: #e74c3c; font-size: 1.1em;">
            🔒 This camera requires a subscription
          </p>
          <button onclick="window.go.main.App.OpenURL('${watchUrl}')" 
                  style="padding: 14px 28px; background: #ff0000; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; font-size: 1.05em; transition: background 0.2s;"
                  onmouseover="this.style.background='#cc0000'" 
                  onmouseout="this.style.background='#ff0000'">
            ▶ Watch Live on YouTube
          </button>
          <p style="margin: 15px 0 0 0; font-size: 0.85em; color: #999;">Opens in your browser</p>
        </div>
      `;
    } else {
      container.innerHTML = `
        <div style="text-align: center; padding: 30px;">
          <button onclick="window.go.main.App.OpenURL('${watchUrl}')" 
                  style="padding: 14px 28px; background: #ff0000; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; font-size: 1.05em; transition: background 0.2s;"
                  onmouseover="this.style.background='#cc0000'" 
                  onmouseout="this.style.background='#ff0000'">
            ▶ Watch Live on YouTube
          </button>
          <p style="margin: 15px 0 0 0; font-size: 0.85em; color: #999;">Opens in your browser</p>
        </div>
      `;
    }
  } else {
    container.innerHTML = `<p style="color: #999; font-size: 0.9em;">No YouTube URL provided.</p>`;
  }
}

// Function to add a marker to the map
function addMarkerToMap(cam) {
  const marker = L.marker([cam.lat, cam.lng]).addTo(map);
  marker.on("click", () => {
    playVideo(cam.youtubeUrl, cam.name, cam.subscriptionRequired);
  });

  // Construct the popup content using template literals for readability and correctness
  const popupContentHtml = `
    <div style="padding: 10px; user-select: none;">
      <strong style="display: block; font-size: 1.1em; margin-bottom:
5px;">${cam.name}</strong>
      <p style="margin: 0 0 8px 0; font-size: 0.9em;">${cam.description || ""}</p>
      <span style="font-size: 0.9em; padding: 5px 10px; border-radius: 12px;
${cam.subscriptionRequired ? 'background-color: #fbe2e2; color: #d9534f;' :
'background-color: #e9f7ee; color: #5cb85c;'}">
        ${cam.subscriptionRequired ? "🔒 Subscription required" : "🔓 Free access"}
      </span>
    </div>`;

  // Create the popup element and set its HTML content
  const popupContent = document.createElement("div");
  popupContent.innerHTML = popupContentHtml;
  marker.bindPopup(popupContent);
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
    li.style.display = "flex";
    li.style.alignItems = "center";
    li.style.justifyContent = "space-between";
    li.style.padding = "8px";
    li.style.cursor = "pointer";

    const camName = document.createElement("span");
    camName.textContent = cam.name;
    camName.style.flex = "1";
    camName.onclick = () => {
      selectCamera(cam);
    };

    // Set or update the badge visually based on status
    const getStatusHtml = (isSub) => {
      if (isSub) {
        return `<span style="background-color: #fbe2e2; color: #d9534f; padding: 4px 8px; border-radius: 15px; font-size: 0.8em; display: inline-flex; align-items: center;"><span style="margin-right: 6px;">🔒</span> Subscription Required</span>`;
      } else {
        return `<span style="background-color: #e9f7ee; color: #5cb85c; padding: 4px 8px; border-radius: 15px; font-size: 0.8em; display: inline-flex; align-items: center;"><span style="margin-right: 6px;">🔓</span> Free Access</span>`;
      }
    };

    const statusBadgeEl = document.createElement("span");
    statusBadgeEl.innerHTML = getStatusHtml(cam.subscriptionRequired);
    statusBadgeEl.title = cam.subscriptionRequired ? "Click to mark as free" : "Click to mark as subscription required";
    // The click handler will now attach to the parent li for better scope control,
    // but since we are modifying this specific section's logic, let's keep it simple
    // and apply the listener directly to a container if possible.

    // Instead of changing the badge element creation completely, I'll use an outer wrapper element
    // that holds the status and attach the click handler there for better CSS scoping.
    const statusContainer = document.createElement("span");
    statusContainer.innerHTML = getStatusHtml(cam.subscriptionRequired);

    // Re-assigning variable to point to the new container element, not just text content
    // This requires updating how 'subscriptionBadge' was used later.
    // For minimal change footprint, I will update the logic around lines 104 and attach the click listener to the statusContainer itself if possible.

    /* Reverting the full structural change for now. Let's keep it as a span and style it better first. */
    const subscriptionBadge = document.createElement("span");
    subscriptionBadge.textContent = cam.subscriptionRequired ? "🔒" : "🔓";
    // Apply more robust styling to visually enhance the difference
    subscriptionBadge.style.backgroundColor = cam.subscriptionRequired ? '#fee' : '#efe';
    subscriptionBadge.style.color = cam.subscriptionRequired ? '#d9534f' : '#5cb85c';
    subscriptionBadge.style.padding = '4px 8px';
    subscriptionBadge.style.borderRadius = '15px';
    subscriptionBadge.style.fontSize = '0.8em';
    subscriptionBadge.title = cam.subscriptionRequired ? "Click to mark as free" : "Click to mark as subscription required";

    // The click handler logic is correct, it just needs the badge variable updated.
    subscriptionBadge.onclick = (e) => {
      e.stopPropagation();
      cam.toggleSubscriptionStatus(); // Re-using this for now if I cannot scope out the full block edit.
                                    // BUT WAIT: The previous fix was to REMOVE cam.toggleSubscriptionStatus().
                                    // I must use the corrected logic (direct property flip) here, too!
      cam.subscriptionRequired = !cam.subscriptionRequired;
      updateCameraList();
      console.log(`Updated ${cam.name}: ${cam.subscriptionRequired ? "Subscription required" : "Free"}`);
    };
    subscriptionBadge.onclick = (e) => {
      e.stopPropagation();
      // Toggle the subscription status property directly on the camera object
      cam.subscriptionRequired = !cam.subscriptionRequired;
      updateCameraList();
      console.log(`Updated ${cam.name}: ${cam.subscriptionRequired ? "Subscription required" : "Free"}`);
    };

    li.appendChild(camName);
    li.appendChild(subscriptionBadge);
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
  playVideo(cam.youtubeUrl, cam.name, cam.subscriptionRequired);
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
