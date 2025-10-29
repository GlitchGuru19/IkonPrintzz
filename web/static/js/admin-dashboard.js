const token = localStorage.getItem('admin_token');
const logoutBtn = document.getElementById('logoutBtn');
const foldersContainer = document.getElementById('foldersContainer');
const totalFoldersEl = document.getElementById('totalFolders');
const totalFilesEl = document.getElementById('totalFiles');
const connectionStatusEl = document.getElementById('connectionStatus');

let ws;
let folders = {};
let allFiles = [];

// Check if token exists
if (!token) {
    window.location.href = '/admin';
}

// Logout
logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('admin_token');
    window.location.href = '/admin';
});

// Initialize WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onopen = () => {
        console.log('WebSocket connected');
        connectionStatusEl.className = 'status-connected';
        connectionStatusEl.textContent = '●';
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected');
        connectionStatusEl.className = 'status-disconnected';
        connectionStatusEl.textContent = '●';
        
        // Reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

function handleWebSocketMessage(message) {
    switch (message.type) {
        case 'new_file':
            addFileToUI(message.payload);
            break;
        case 'file_deleted':
            removeFileFromUI(message.payload.id);
            break;
        case 'folder_created':
            addFolderToUI(message.payload);
            break;
    }
}

function addFileToUI(file) {
    allFiles.push(file);
    
    if (!folders[file.folder_id]) {
        folders[file.folder_id] = {
            id: file.folder_id,
            name: file.folder_name,
            files: []
        };
    }
    
    folders[file.folder_id].files.push(file);
    updateStats();
    renderFolders();
}

function removeFileFromUI(fileId) {
    const fileIndex = allFiles.findIndex(f => f.id === fileId);
    if (fileIndex !== -1) {
        const file = allFiles[fileIndex];
        allFiles.splice(fileIndex, 1);
        
        if (folders[file.folder_id]) {
            folders[file.folder_id].files = folders[file.folder_id].files.filter(f => f.id !== fileId);
        }
        
        updateStats();
        renderFolders();
    }
}

function addFolderToUI(folder) {
    if (!folders[folder.id]) {
        folders[folder.id] = {
            id: folder.id,
            name: folder.name,
            files: []
        };
        updateStats();
        renderFolders();
    }
}

// Fetch initial data
async function fetchData() {
    try {
        const response = await fetch('/api/files', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch files');
        }

        allFiles = await response.json() || [];
        
        // Organize files by folder
        folders = {};
        allFiles.forEach(file => {
            if (!folders[file.folder_id]) {
                folders[file.folder_id] = {
                    id: file.folder_id,
                    name: file.folder_name,
                    files: []
                };
            }
            folders[file.folder_id].files.push(file);
        });

        updateStats();
        renderFolders();
    } catch (error) {
        console.error('Error fetching data:', error);
        if (error.message.includes('401')) {
            localStorage.removeItem('admin_token');
            window.location.href = '/admin';
        }
    }
}

function updateStats() {
    totalFoldersEl.textContent = Object.keys(folders).length;
    totalFilesEl.textContent = allFiles.length;
}

function renderFolders() {
    if (Object.keys(folders).length === 0) {
        foldersContainer.innerHTML = `
            <div class="empty-state">
                <h3>📭 No files yet</h3>
                <p>Waiting for users to upload files...</p>
            </div>
        `;
        return;
    }

    foldersContainer.innerHTML = '';
    
    Object.values(folders).forEach(folder => {
        if (folder.files.length === 0) return;

        const folderCard = document.createElement('div');
        folderCard.className = 'folder-card';
        
        const folderHeader = document.createElement('div');
        folderHeader.className = 'folder-header';
        folderHeader.innerHTML = `
            <h3>📁 ${folder.name}</h3>
            <span class="folder-info">${folder.files.length} file(s)</span>
        `;
        
        const filesGrid = document.createElement('div');
        filesGrid.className = 'files-grid';
        
        folder.files.forEach(file => {
            const fileCard = createFileCard(file);
            filesGrid.appendChild(fileCard);
        });
        
        folderCard.appendChild(folderHeader);
        folderCard.appendChild(filesGrid);
        foldersContainer.appendChild(folderCard);
    });
}

function createFileCard(file) {
    const card = document.createElement('div');
    card.className = 'file-card';
    card.innerHTML = `
        <div class="file-name">${file.file_name}</div>
        <div class="file-meta">
            ${formatFileSize(file.file_size)} • ${file.file_type.toUpperCase()}
        </div>
        <div class="file-actions">
            <button class="btn btn-print" onclick="printFile('${file.id}')">🖨️ Print</button>
            <button class="btn btn-danger" onclick="deleteFile('${file.id}')">🗑️ Delete</button>
        </div>
    `;
    return card;
}

async function printFile(fileId) {
    // Open file in new window for printing
    const printWindow = window.open(`/api/files/${fileId}/view`, '_blank');
    
    // Wait for window to load then trigger print dialog
    printWindow.onload = () => {
        setTimeout(() => {
            printWindow.print();
        }, 500);
    };
}

async function deleteFile(fileId) {
    if (!confirm('Are you sure you want to delete this file?')) {
        return;
    }

    try {
        const response = await fetch(`/api/files/${fileId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to delete file');
        }

        // File will be removed via WebSocket message
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Initialize
connectWebSocket();
fetchData();
