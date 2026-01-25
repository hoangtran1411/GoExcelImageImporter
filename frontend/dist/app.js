// Global variable to store update info
let updateInfo = null;
let currentOutputPath = null;

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

    // Init worker count
    initWorkerCount();
});

// Initialize worker count based on CPU cores
async function initWorkerCount() {
    try {
        const cores = await window.go.main.App.GetCPUCount();
        if (cores > 0) {
            // Ensure element exists before setting value
            const input = document.getElementById('workerCount');
            if (input) {
                input.value = cores;
            }
        }
    } catch (err) {
        console.error('Failed to get CPU count:', err);
    }
}

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
    // 1. Check if we are in "Open File" mode
    if (currentOutputPath) {
        try {
            await window.go.main.App.OpenFileLocation(currentOutputPath);
        } catch (err) {
            showStatus('Could not open file: ' + err, 'error');
        }

        // Reset UI to initial state
        currentOutputPath = null;
        const btn = document.getElementById('processBtn');
        btn.classList.remove('btn-success'); // Assuming you might add styling
        btn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none"><path d="M8 5V19L19 12L8 5Z" fill="currentColor"/></svg>
            Start Processing
        `;
        updateProgress(0);
        return;
    }

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



    let processingSuccess = false;

    try {
        const result = await window.go.main.App.Process(config);

        if (result.success) {
            showStatus('✓ ' + result.message, 'success');
            processingSuccess = true;

            if (result.missingCodes && result.missingCodes.length > 0) {
                console.log('Missing codes:', result.missingCodes);
            }

            // Store output path for the button action
            if (result.outputPath) {
                currentOutputPath = result.outputPath;
            }
        } else {
            showStatus('✗ ' + result.message, 'error');
        }
    } catch (err) {
        showStatus('Error: ' + err, 'error');
    } finally {
        // Reset button
        btn.disabled = false;

        if (processingSuccess && currentOutputPath) {
            // Change button to "Open Output File"
            btn.innerHTML = `
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                </svg>
                Open Output File
            `;
            // Keep progress at 100%
            updateProgress(100);
        } else {
            // Revert to "Start Processing" on failure
            btn.innerHTML = `
                <svg viewBox="0 0 24 24" fill="none"><path d="M8 5V19L19 12L8 5Z" fill="currentColor"/></svg>
                Start Processing
            `;
            // Only reset progress if failed? Or keep it? keeping logic simple
            if (!processingSuccess) {
                // Maybe keep progress to show it failed? But user might want to retry
            }
        }
    }
}

// Update progress bar
function updateProgress(percent) {
    const fill = document.getElementById('progressFill');
    const text = document.getElementById('progressText');

    fill.style.width = percent + '%';
    text.textContent = Math.round(percent) + '%';
}

// Show sliding toast notification
function showStatus(message, type) {
    const container = document.getElementById('toast-container');

    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    // Icon based on type
    let iconSvg = '';
    if (type === 'success') {
        iconSvg = '<svg viewBox="0 0 24 24" fill="none"><path d="M20 6L9 17L4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>';
    } else if (type === 'error') {
        iconSvg = '<svg viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>';
    } else {
        iconSvg = '<svg viewBox="0 0 24 24" fill="none"><path d="M13 16H11V12H13M12 8H12.01M21 12C21 16.97 16.97 21 12 21C7.03 21 3 16.97 3 12C3 7.03 7.03 3 12 3C16.97 3 21 7.03 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>';
    }

    toast.innerHTML = `
        <div class="toast-icon">${iconSvg}</div>
        <div class="toast-message">${message}</div>
    `;

    // Append to container
    container.appendChild(toast);

    // Auto remove after 4 seconds
    setTimeout(() => {
        toast.classList.add('hiding');
        toast.addEventListener('transitionend', () => {
            if (toast.parentElement) {
                toast.remove();
            }
        });
    }, 4000);
}
