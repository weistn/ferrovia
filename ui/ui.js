const svgNS = "http://www.w3.org/2000/svg";

// Call when DOM has loaded, connect with the server and enabe the UI
async function init() {
    // Install event listeners before connecting to the server.
    go.addEventListener("canvas", (data) => {renderCanvas(data, document.getElementById("view2d"), document.getElementById("view2d-measure"), document.getElementById("view2d-ground"), document.getElementById("view2d-tracks"))});

    document.getElementById("select-switchtower").addEventListener("click", () => {
        document.getElementById("trackdiagram").style.display = "block";
        document.getElementById("view2d").style.display = "none";
    });

    document.getElementById("select-view2d").addEventListener("click", () => {
        document.getElementById("trackdiagram").style.display = "none";
        document.getElementById("view2d").style.display = "block";
    });

    // Connect to the server
    await go.connect();
}
