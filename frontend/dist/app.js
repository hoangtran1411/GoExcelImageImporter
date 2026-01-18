// Global variable to store update info
let updateInfo = null;

// Wait for Wails runtime to be ready
document.addEventListener('DOMContentLoaded', function () {
    // Listen for progress updates from Go backend
    if (typeof runtime !== 'undefined') {
        runtime.EventsOn('progress', function (progress) {
            updateProgress(progress);
        });

        // Listen for update progress
        runtime.EventsOn('updateProgress', function (message) {
            showStatus(message, 'info');
        });
    }

    // Check for updates on startup
    checkForUpdates();

    // Load current version
    loadVersion();
});

// Load and display current version
async function loadVersion() {
    try {
        const version = await window.go.main.App.GetCurrentVersion();
        document.getElementById('versionText').textContent = version;
    } catch (err) {
        console.error('Failed to get version:', err);
    }
}

// Check for updates
async function checkForUpdates() {
    try {
        updateInfo = await window.go.main.App.CheckForUpdate();

        if (updateInfo && updateInfo.available) {
            // Show update button
            const updateBtn = document.getElementById('updateBtn');
            updateBtn.classList.add('visible');
            updateBtn.title = `Update to ${updateInfo.latestVersion} available! Click to install`;
            console.log('Update available:', updateInfo.latestVersion);
        }
    } catch (err) {
        console.error('Failed to check for updates:', err);
    }
}

// Perform the update
async function performUpdate() {
    if (!updateInfo || !updateInfo.downloadUrl) {
        showStatus('No update information available', 'error');
        return;
    }

    const updateBtn = document.getElementById('updateBtn');
    updateBtn.disabled = true;

    showStatus(`Downloading ${updateInfo.latestVersion}...`, 'info');

    try {
        const result = await window.go.main.App.PerformUpdate(updateInfo.downloadUrl);
        if (result) {
            showStatus('Update installed! Restarting...', 'success');
        }
    } catch (err) {
        showStatus('Update failed: ' + err, 'error');
        updateBtn.disabled = false;
    }
}

// Select Excel file
async function selectExcel() {
    try {
        const path = await window.go.main.App.SelectExcelFile();
        if (path) {
            document.getElementById('excelPath').value = path;
            loadSheets(path);
        }
    } catch (err) {
        showStatus('Error selecting file: ' + err, 'error');
    }
}

// Select image folder
async function selectImageFolder() {
    try {
        const path = await window.go.main.App.SelectImageFolder();
        if (path) {
            document.getElementById('imageDir').value = path;
        }
    } catch (err) {
        showStatus('Error selecting folder: ' + err, 'error');
    }
}

// Load sheets from Excel file
async function loadSheets(excelPath) {
    try {
        const sheets = await window.go.main.App.GetSheets(excelPath);
        const select = document.getElementById('sheetName');

        // Clear existing options
        select.innerHTML = '<option value="">Select sheet...</option>';

        // Add sheets
        sheets.forEach(function (sheet) {
            const option = document.createElement('option');
            option.value = sheet;
            option.textContent = sheet;
            select.appendChild(option);
        });

        // Select first sheet by default
        if (sheets.length > 0) {
            select.value = sheets[0];
        }
    } catch (err) {
        showStatus('Error loading sheets: ' + err, 'error');
    }
}

// Start processing
async function startProcess() {
    const excelPath = document.getElementById('excelPath').value;
    const imageDir = document.getElementById('imageDir').value;

    if (!excelPath) {
        showStatus('Please select an Excel file', 'error');
        return;
    }

    if (!imageDir) {
        showStatus('Please select an image folder', 'error');
        return;
    }

    // Build config object
    const config = {
        excelPath: excelPath,
        imageDir: imageDir,
        codeCol: document.getElementById('codeCol').value || 'A',
        imageCol: document.getElementById('imageCol').value || 'F',
        sheetName: document.getElementById('sheetName').value,
        rowHeight: parseFloat(document.getElementById('rowHeight').value) || 105,
        colWidth: parseFloat(document.getElementById('colWidth').value) || 20,
        workerCount: parseInt(document.getElementById('workerCount').value) || 10
    };

    // Show progress bar
    const progressContainer = document.getElementById('progressContainer');
    progressContainer.classList.add('active');
    updateProgress(0);

    // Disable button
    const btn = document.getElementById('processBtn');
    btn.disabled = true;
    btn.innerHTML = `
        <svg class="processing" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.3"/>
            <path d="M12 2C6.48 2 2 6.48 2 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        Processing...
    `;

    showStatus('Processing images...', 'info');

    try {
        const result = await window.go.main.App.Process(config);

        if (result.success) {
            showStatus('✓ ' + result.message, 'success');

            if (result.missingCodes && result.missingCodes.length > 0) {
                console.log('Missing codes:', result.missingCodes);
            }
        } else {
            showStatus('✗ ' + result.message, 'error');
        }
    } catch (err) {
        showStatus('Error: ' + err, 'error');
    } finally {
        // Reset button
        btn.disabled = false;
        btn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none"><path d="M8 5V19L19 12L8 5Z" fill="currentColor"/></svg>
            Start Processing
        `;

        // Keep progress at 100% on success
        updateProgress(100);
    }
}

// Update progress bar
function updateProgress(percent) {
    const fill = document.getElementById('progressFill');
    const text = document.getElementById('progressText');

    fill.style.width = percent + '%';
    text.textContent = Math.round(percent) + '%';
}

// Show status message
function showStatus(message, type) {
    const section = document.getElementById('statusSection');
    const messageEl = document.getElementById('statusMessage');

    section.classList.add('active');
    messageEl.className = 'status-message ' + type;
    messageEl.textContent = message;
}
